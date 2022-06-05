package runner

import (
	"context"
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
	Run(context.Context) (Result, error)
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

func (r *runner) Init(ctx context.Context) error {
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

func (r *runner) Run(ctx context.Context) (Result, error) {
	tasks := func() []*proto.Task {
		var t []*proto.Task
		for _, item := range r.cfg.Data.Spec.Tasks {
			t = append(t, &proto.Task{
				Name:     item.Name,
				Commands: item.Commands,
				Depends:  item.Depends,
			})
		}
		return t
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Config.Spec.Scheduler.Timeout)*time.Second)
	defer cancel()

	reply, err := r.client.SendServer(ctx, &proto.ServerRequest{
		ApiVersion: r.cfg.Data.ApiVersion,
		Kind:       r.cfg.Data.Kind,
		Metadata: &proto.Metadata{
			Name: r.cfg.Data.Metadata.Name,
		},
		Spec: &proto.Spec{
			Tasks: tasks,
		},
	})

	if err != nil {
		return Result{}, errors.Wrap(err, "failed to send")
	}

	return Result{Output: reply.GetOutput(), Error: reply.GetError()}, nil
}
