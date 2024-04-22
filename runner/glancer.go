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

type Glancer interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (GlanceReply, error)
}

type GlancerConfig struct {
	Config config.Config
	Data   Proto
}

type glancer struct {
	cfg    *GlancerConfig
	client proto.ServerProtoClient
	conn   *grpc.ClientConn
}

var (
	GlancerUnitMap = map[string]time.Duration{
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
	}
)

func GlancerNew(_ context.Context, cfg *GlancerConfig) Glancer {
	return &glancer{
		cfg: cfg,
	}
}

func GlancerDefaultConfig() *GlancerConfig {
	return &GlancerConfig{}
}

func (g *glancer) Init(ctx context.Context) error {
	if err := g.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	return nil
}

func (g *glancer) Deinit(ctx context.Context) error {
	_ = g.deinitConn(ctx)

	return nil
}

// nolint: funlen
func (g *glancer) Run(ctx context.Context) (rep GlanceReply, err error) {
	entries := func(ent []*proto.GlanceEntry) []GlanceEntry {
		var buf []GlanceEntry
		for _, item := range ent {
			buf = append(buf, GlanceEntry{
				Name:  item.GetName(),
				IsDir: item.GetIsDir(),
				Size:  item.GetSize(),
				Time:  item.GetTime(),
				User:  item.GetUser(),
				Group: item.GetGroup(),
				Mode:  item.GetMode(),
			})
		}
		return buf
	}

	threads := func(threads []*proto.GlanceThread) []GlanceThread {
		var buf []GlanceThread
		for _, item := range threads {
			buf = append(buf, GlanceThread{
				Name:    item.GetName(),
				Cmdline: item.GetCmdline(),
				Memory:  item.GetMemory(),
				Time:    item.GetTime(),
				Pid:     item.GetPid(),
			})
		}
		return buf
	}

	processes := func(proc []*proto.GlanceProcess) []GlanceProcess {
		var buf []GlanceProcess
		for _, item := range proc {
			buf = append(buf, GlanceProcess{
				Process: GlanceThread{
					Name:    item.GetProcess().GetName(),
					Cmdline: item.GetProcess().GetCmdline(),
					Memory:  item.GetProcess().GetMemory(),
					Time:    item.GetProcess().GetTime(),
					Pid:     item.GetProcess().GetPid(),
				},
				Threads: threads(item.GetThreads()),
			})
		}
		return buf
	}

	output := func(r *proto.GlanceReply) GlanceReply {
		return GlanceReply{
			Dir: GlanceDirRep{
				Entries: entries(r.GetDir().GetEntries()),
			},
			File: GlanceFileRep{
				Content:  r.GetFile().GetContent(),
				Readable: r.GetFile().GetReadable(),
			},
			Sys: GlanceSysRep{
				Resource: GlanceResource{
					Allocatable: GlanceAllocatable{
						MilliCPU: r.GetSys().GetResource().GetAllocatable().GetMilliCPU(),
						Memory:   r.GetSys().GetResource().GetAllocatable().GetMemory(),
						Storage:  r.GetSys().GetResource().GetAllocatable().GetStorage(),
					},
					Requested: GlanceRequested{
						MilliCPU: r.GetSys().GetResource().GetRequested().GetMilliCPU(),
						Memory:   r.GetSys().GetResource().GetRequested().GetMemory(),
						Storage:  r.GetSys().GetResource().GetRequested().GetStorage(),
					},
				},
				Stats: GlanceStats{
					CPU: GlanceCPU{
						Total: r.GetSys().GetStats().GetCpu().GetTotal(),
						Used:  r.GetSys().GetStats().GetCpu().GetUsed(),
					},
					Host: r.GetSys().GetStats().GetHost(),
					Memory: GlanceMemory{
						Total: r.GetSys().GetStats().GetMemory().GetTotal(),
						Used:  r.GetSys().GetStats().GetMemory().GetUsed(),
					},
					OS: r.GetSys().GetStats().GetOs(),
					Storage: GlanceStorage{
						Total: r.GetSys().GetStats().GetStorage().GetTotal(),
						Used:  r.GetSys().GetStats().GetStorage().GetUsed(),
					},
					Processes: processes(r.GetSys().GetStats().GetProcesses()),
				},
			},
			Error: r.GetError(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, g.setTimeout(g.cfg.Data.Spec.Glance.Timeout))
	defer cancel()

	reply, e := g.client.SendGlance(ctx)
	defer func() {
		_ = reply.CloseSend()
	}()

	if e != nil {
		return rep, errors.Wrap(e, "failed to set")
	}

	if e = reply.Send(&proto.GlanceRequest{
		ApiVersion: g.cfg.Data.ApiVersion,
		Kind:       g.cfg.Data.Kind,
		Metadata: &proto.GlanceMetadata{
			Name: g.cfg.Data.Metadata.Name,
		},
		Spec: &proto.GlanceSpec{
			Glance: &proto.Glance{
				Dir: &proto.GlanceDirReq{
					Path: g.cfg.Data.Spec.Glance.Dir.Path,
				},
				File: &proto.GlanceFileReq{
					Path:    g.cfg.Data.Spec.Glance.File.Path,
					MaxSize: g.cfg.Data.Spec.Glance.File.MaxSize,
				},
				Sys: &proto.GlanceSysReq{
					Enable: g.cfg.Data.Spec.Glance.Sys.Enable,
				},
			},
		},
	}); e != nil {
		return rep, errors.Wrap(e, "failed to send")
	}

	recv, e := reply.Recv()
	if e != nil {
		return rep, errors.Wrap(e, "failed to recv")
	}

	return output(recv), nil
}

func (g *glancer) initConn(_ context.Context) error {
	var err error

	host := g.cfg.Config.Spec.Runner.Host
	port := g.cfg.Config.Spec.Runner.Port

	g.conn, err = grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	g.client = proto.NewServerProtoClient(g.conn)

	return nil
}

func (g *glancer) deinitConn(_ context.Context) error {
	return g.conn.Close()
}

func (g *glancer) setTimeout(timeout GlanceTimeout) time.Duration {
	tm := int64(Time)
	unit := int64(GlancerUnitMap[Unit])

	if timeout.Time != 0 {
		tm = timeout.Time
	}

	if timeout.Unit != "" {
		if val, ok := GlancerUnitMap[timeout.Unit]; ok {
			unit = int64(val)
		}
	}

	return time.Duration(tm * unit)
}
