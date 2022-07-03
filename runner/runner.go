package runner

import (
	"context"
	"io"
	"math"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipego/cli/config"
	"github.com/pipego/cli/dag"
	proto "github.com/pipego/cli/runner/proto"
	livelog "github.com/pipego/dag/runner"
)

const (
	LIVELOG = 5000
)

type Runner interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
	Tail(ctx context.Context) livelog.Livelog
	Tasks(ctx context.Context) []Task
}

type Config struct {
	Config config.Config
	Dag    dag.DAG
	Data   Proto
}

type runner struct {
	cfg    *Config
	client proto.ServerProtoClient
	conn   *grpc.ClientConn
	log    livelog.Livelog
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
	if err := r.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	if err := r.initDag(ctx); err != nil {
		return errors.Wrap(err, "failed to init dag")
	}

	return nil
}

func (r *runner) Deinit(ctx context.Context) error {
	_ = r.deinitDag(ctx)
	_ = r.deinitConn(ctx)

	return nil
}

func (r *runner) Run(ctx context.Context) error {
	return r.runDag(ctx)
}

func (r *runner) Tail(ctx context.Context) livelog.Livelog {
	return r.log
}

func (r *runner) Tasks(_ context.Context) []Task {
	return r.cfg.Data.Spec.Tasks
}

func (r *runner) initConn(_ context.Context) error {
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

func (r *runner) deinitConn(_ context.Context) error {
	return r.conn.Close()
}

func (r *runner) initDag(ctx context.Context) error {
	var tasks []dag.Task

	for _, item := range r.cfg.Data.Spec.Tasks {
		tasks = append(tasks, dag.Task{
			Name:     item.Name,
			Commands: item.Commands,
			Depends:  item.Depends,
		})
	}

	r.log = livelog.Livelog{
		Error: make(chan error, LIVELOG),
		Line:  make(chan *livelog.Line, LIVELOG),
	}

	return r.cfg.Dag.Init(ctx, tasks)
}

func (r *runner) deinitDag(ctx context.Context) error {
	_ = r.cfg.Dag.Deinit(ctx)

	close(r.log.Error)
	close(r.log.Line)

	return nil
}

func (r *runner) runDag(ctx context.Context) error {
	return r.cfg.Dag.Run(ctx, r.routine, r.log)
}

func (r *runner) routine(name string, args []string, log livelog.Livelog) error {
	task := func() *proto.Task {
		return &proto.Task{
			Name:     name,
			Commands: args,
		}
	}()

	output := func(s proto.ServerProto_SendServerClient) {
		done := make(chan bool)
		go func() {
			for {
				recv, err := s.Recv()
				if err == io.EOF {
					done <- true
					return
				}
				if err != nil {
					log.Error <- err
					return
				}
				log.Line <- &livelog.Line{
					Pos:     recv.GetOutput().GetPos(),
					Time:    recv.GetOutput().GetTime(),
					Message: recv.GetOutput().GetMessage(),
				}
			}
		}()
		<-done
	}

	reply, err := r.client.SendServer(context.Background(), &proto.ServerRequest{
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
		return errors.Wrap(err, "failed to send")
	}

	output(reply)

	return nil
}
