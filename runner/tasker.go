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

type Tasker interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
	Tail(ctx context.Context) dagRunner.Livelog
	Tasks(ctx context.Context) []Task
}

type TaskerConfig struct {
	Config config.Config
	Dag    dag.DAG
	Data   Proto
}

type tasker struct {
	cfg    *TaskerConfig
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

func TaskerNew(_ context.Context, cfg *TaskerConfig) Tasker {
	return &tasker{
		cfg: cfg,
	}
}

func TaskerDefaultConfig() *TaskerConfig {
	return &TaskerConfig{}
}

func (t *tasker) Init(ctx context.Context) error {
	if err := t.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	if err := t.initDag(ctx); err != nil {
		return errors.Wrap(err, "failed to init dag")
	}

	return nil
}

func (t *tasker) Deinit(ctx context.Context) error {
	_ = t.deinitDag(ctx)
	_ = t.deinitConn(ctx)

	return nil
}

func (t *tasker) Run(ctx context.Context) error {
	return t.runDag(ctx)
}

func (t *tasker) Tail(ctx context.Context) dagRunner.Livelog {
	return t.log
}

func (t *tasker) Tasks(_ context.Context) []Task {
	return t.cfg.Data.Spec.Tasks
}

func (t *tasker) initConn(_ context.Context) error {
	var err error

	host := t.cfg.Config.Spec.Runner.Host
	port := t.cfg.Config.Spec.Runner.Port

	t.conn, err = grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	t.client = proto.NewServerProtoClient(t.conn)

	return nil
}

func (t *tasker) deinitConn(_ context.Context) error {
	return t.conn.Close()
}

func (t *tasker) initDag(ctx context.Context) error {
	helper := func(p []TaskParam) []dagRunner.Param {
		var buf []dagRunner.Param
		for _, item := range p {
			buf = append(buf, dagRunner.Param{
				Name:  item.Name,
				Value: item.Value,
			})
		}
		return buf
	}

	var tasks []dag.Task

	for i := range t.cfg.Data.Spec.Tasks {
		tasks = append(tasks, dag.Task{
			Name:     t.cfg.Data.Spec.Tasks[i].Name,
			File:     dagRunner.File(t.cfg.Data.Spec.Tasks[i].File),
			Params:   helper(t.cfg.Data.Spec.Tasks[i].Params),
			Commands: t.cfg.Data.Spec.Tasks[i].Commands,
			Livelog:  t.cfg.Data.Spec.Tasks[i].Livelog,
			Depends:  t.cfg.Data.Spec.Tasks[i].Depends,
		})
	}

	t.log = dagRunner.Livelog{
		Error: make(chan error, Livelog),
		Line:  make(chan *dagRunner.Line, Livelog),
	}

	return t.cfg.Dag.Init(ctx, tasks)
}

func (t *tasker) deinitDag(ctx context.Context) error {
	_ = t.cfg.Dag.Deinit(ctx)

	close(t.log.Error)
	close(t.log.Line)

	return nil
}

func (t *tasker) runDag(ctx context.Context) error {
	return t.cfg.Dag.Run(ctx, t.routine, t.log)
}

func (t *tasker) routine(name string, file dagRunner.File, envs []dagRunner.Param, args []string, _len int64, log dagRunner.Livelog) error {
	params := func(p []dagRunner.Param) []*proto.TaskParam {
		var buf []*proto.TaskParam
		for _, item := range p {
			buf = append(buf, &proto.TaskParam{
				Name:  item.Name,
				Value: item.Value,
			})
		}
		return buf
	}

	task := func() *proto.Task {
		return &proto.Task{
			Name: name,
			File: &proto.TaskFile{
				Content: t.contentHelper([]byte(file.Content), file.Gzip),
				Gzip:    file.Gzip,
			},
			Params:   params(envs),
			Commands: args,
			Livelog:  _len,
		}
	}()

	output := func(s proto.ServerProto_SendTaskClient) {
		done := make(chan bool)
		go func(s proto.ServerProto_SendTaskClient, log dagRunner.Livelog, done chan bool) {
			for {
				if recv, err := s.Recv(); err == nil {
					log.Line <- &dagRunner.Line{
						Pos:     recv.GetOutput().GetPos(),
						Time:    recv.GetOutput().GetTime(),
						Message: recv.GetOutput().GetMessage(),
					}
					if recv.GetOutput().GetMessage() == "EOF" {
						_ = s.CloseSend()
						done <- true
						return
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

	ctx, cancel := context.WithTimeout(context.Background(), t.setTimeout(name))
	defer cancel()

	reply, err := t.client.SendTask(ctx)
	defer func() {
		_ = reply.CloseSend()
	}()

	if err != nil {
		return errors.Wrap(err, "failed to set")
	}

	if err := reply.Send(&proto.TaskRequest{
		ApiVersion: t.cfg.Data.ApiVersion,
		Kind:       t.cfg.Data.Kind,
		Metadata: &proto.TaskMetadata{
			Name: t.cfg.Data.Metadata.Name,
		},
		Spec: &proto.TaskSpec{
			Task: task,
		},
	}); err != nil {
		return errors.Wrap(err, "failed to send")
	}

	output(reply)

	return nil
}

func (t *tasker) contentHelper(data []byte, compressed bool) []byte {
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

func (t *tasker) setTimeout(name string) time.Duration {
	helper := func(_t Task) time.Duration {
		tm := int64(Time)
		unit := int64(UnitMap[Unit])
		if _t.Timeout.Time != 0 {
			tm = _t.Timeout.Time
		}
		if _t.Timeout.Unit != "" {
			if val, ok := UnitMap[_t.Timeout.Unit]; ok {
				unit = int64(val)
			}
		}
		return time.Duration(tm * unit)
	}

	var _t Task

	for i := range t.cfg.Data.Spec.Tasks {
		if name == t.cfg.Data.Spec.Tasks[i].Name {
			_t = t.cfg.Data.Spec.Tasks[i]
			break
		}
	}

	return helper(_t)
}
