package runner

import (
	"context"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipego/cli/config"
	proto "github.com/pipego/cli/runner/proto"
)

type Runner interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) ([]Result, error)
}

type Config struct {
	Config config.Config
	Data   Proto
}

type runner struct {
	cfg    *Config
	client proto.ServerProtoClient
	conn   *grpc.ClientConn
}

func New(_ context.Context, cfg *Config) Runner {
	return &runner{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *runner) Init(_ context.Context) error {
	var err error

	host := r.cfg.Config.Spec.Runner.Host
	port := r.cfg.Config.Spec.Runner.Port

	r.conn, err = grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	r.client = proto.NewServerProtoClient(r.conn)

	return nil
}

func (r *runner) Deinit(_ context.Context) error {
	return r.conn.Close()
}

func (r *runner) Run(ctx context.Context) ([]Result, error) {
	task := func() *proto.Task {
		return &proto.Task{
			Name:     r.cfg.Data.Spec.Task.Name,
			Commands: r.cfg.Data.Spec.Task.Commands,
		}
	}()

	output := func(s proto.ServerProto_SendServerClient) []Result {
		var res []Result
		done := make(chan bool)
		go func() {
			for {
				recv, err := s.Recv()
				if err == io.EOF {
					done <- true
					return
				}
				if err != nil {
					res = append(res, Result{Error: err.Error()})
					return
				}
				output := Output{
					Pos:     recv.GetOutput().GetPos(),
					Time:    recv.GetOutput().GetTime(),
					Message: recv.GetOutput().GetMessage(),
				}
				res = append(res, Result{Output: output, Error: recv.GetError()})
			}
		}()
		<-done
		return res
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Config.Spec.Scheduler.Timeout)*time.Second)
	defer cancel()

	reply, err := r.client.SendServer(ctx, &proto.ServerRequest{
		ApiVersion: r.cfg.Data.ApiVersion,
		Kind:       r.cfg.Data.Kind,
		Metadata: &proto.Metadata{
			Name: r.cfg.Data.Metadata.Name,
		},
		Spec: &proto.Spec{
			Task: task,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to send")
	}

	return output(reply), nil
}
