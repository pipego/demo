package runner

import (
	"bytes"
	"compress/gzip"
	"context"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipego/cli/config"
	"github.com/pipego/cli/dag"
	proto "github.com/pipego/cli/runner/proto"
	dagRunner "github.com/pipego/dag/runner"
)

const (
	Livelog = 5000
	Time    = 12
	Unit    = "hour"
)

type Runner interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
	Tail(ctx context.Context) dagRunner.Livelog
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
	log    dagRunner.Livelog
}

var (
	UnitMap = map[string]time.Duration{
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
	}
)

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

func (r *runner) Tail(ctx context.Context) dagRunner.Livelog {
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
			File:     dagRunner.File(item.File),
			Commands: item.Commands,
			Livelog:  item.Livelog,
			Depends:  item.Depends,
		})
	}

	r.log = dagRunner.Livelog{
		Error: make(chan error, Livelog),
		Line:  make(chan *dagRunner.Line, Livelog),
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

func (r *runner) routine(name string, file dagRunner.File, args []string, _len int64, log dagRunner.Livelog) error {
	task := func() *proto.Task {
		return &proto.Task{
			Name: name,
			File: &proto.File{
				Content: r.contentHelper([]byte(file.Content), file.Gzip),
				Gzip:    file.Gzip,
			},
			Commands: args,
			Livelog:  _len,
		}
	}()

	output := func(s proto.ServerProto_SendServerClient) {
		done := make(chan bool)
		go func(s proto.ServerProto_SendServerClient, log dagRunner.Livelog, done chan bool) {
			for {
				if recv, err := s.Recv(); err == nil {
					if recv.GetOutput().GetMessage() == "EOF" {
						_ = s.CloseSend()
						done <- true
						return
					}
					log.Line <- &dagRunner.Line{
						Pos:     recv.GetOutput().GetPos(),
						Time:    recv.GetOutput().GetTime(),
						Message: recv.GetOutput().GetMessage(),
					}
				} else {
					_ = s.CloseSend()
					done <- false
					return
				}
			}
		}(s, log, done)
		<-done
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.setTimeout(name))
	defer cancel()

	reply, err := r.client.SendServer(ctx)
	defer func() {
		_ = reply.CloseSend()
	}()

	if err != nil {
		return errors.Wrap(err, "failed to set")
	}

	if err := reply.Send(&proto.ServerRequest{
		ApiVersion: r.cfg.Data.ApiVersion,
		Kind:       r.cfg.Data.Kind,
		Metadata: &proto.Metadata{
			Name: r.cfg.Data.Metadata.Name,
		},
		Spec: &proto.Spec{
			Task: task,
		},
	}); err != nil {
		return errors.Wrap(err, "failed to send")
	}

	output(reply)

	return nil
}

func (r *runner) contentHelper(data []byte, compressed bool) []byte {
	var b bytes.Buffer

	if !compressed {
		return data
	}

	gz := gzip.NewWriter(&b)
	defer func(gz *gzip.Writer) {
		_ = gz.Close()
	}(gz)

	if _, err := gz.Write(data); err != nil {
		return data
	}

	if err := gz.Flush(); err != nil {
		return data
	}

	return b.Bytes()
}

func (r *runner) setTimeout(name string) time.Duration {
	helper := func(t Task) time.Duration {
		tm := int64(Time)
		unit := int64(UnitMap[Unit])
		if t.Timeout.Time != 0 {
			tm = t.Timeout.Time
		}
		if t.Timeout.Unit != "" {
			if val, ok := UnitMap[t.Timeout.Unit]; ok {
				unit = int64(val)
			}
		}
		return time.Duration(tm * unit)
	}

	var t Task

	for _, item := range r.cfg.Data.Spec.Tasks {
		if name == item.Name {
			t = item
			break
		}
	}

	return helper(t)
}
