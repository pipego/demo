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

type Mainter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (MaintReply, error)
}

type MainterConfig struct {
	Config config.Config
	Data   Proto
}

type mainter struct {
	cfg    *MainterConfig
	client proto.ServerProtoClient
	conn   *grpc.ClientConn
}

func MainterNew(_ context.Context, cfg *MainterConfig) Mainter {
	return &mainter{
		cfg: cfg,
	}
}

func MainterDefaultConfig() *MainterConfig {
	return &MainterConfig{}
}

func (m *mainter) Init(ctx context.Context) error {
	if err := m.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	return nil
}

func (m *mainter) Deinit(ctx context.Context) error {
	_ = m.deinitConn(ctx)

	return nil
}

func (m *mainter) Run(ctx context.Context) (rep MaintReply, err error) {
	output := func(r *proto.MaintReply) MaintReply {
		return MaintReply{
			Clock: MaintClockRep{
				Sync: MaintClockSync{
					Status: r.GetClock().GetSync().GetStatus(),
				},
				Diff: MaintClockDiff{
					Time:      r.GetClock().GetDiff().GetTime(),
					Dangerous: r.GetClock().GetDiff().GetDangerous(),
				},
			},
			Error: r.GetError(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, m.setTimeout(m.cfg.Data.Spec.Maint.Timeout))
	defer cancel()

	reply, e := m.client.SendMaint(ctx)
	defer func() {
		_ = reply.CloseSend()
	}()

	if e != nil {
		return rep, errors.Wrap(e, "failed to set")
	}

	if e = reply.Send(&proto.MaintRequest{
		ApiVersion: m.cfg.Data.ApiVersion,
		Kind:       m.cfg.Data.Kind,
		Metadata: &proto.MaintMetadata{
			Name: m.cfg.Data.Metadata.Name,
		},
		Spec: &proto.MaintSpec{
			Maint: &proto.Maint{
				Clock: &proto.MaintClockReq{
					Sync: m.cfg.Data.Spec.Maint.Clock.Sync,
					Time: m.cfg.Data.Spec.Maint.Clock.Time,
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

func (m *mainter) initConn(_ context.Context) error {
	var err error

	host := m.cfg.Config.Spec.Runner.Host
	port := m.cfg.Config.Spec.Runner.Port

	m.conn, err = grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	m.client = proto.NewServerProtoClient(m.conn)

	return nil
}

func (m *mainter) deinitConn(_ context.Context) error {
	return m.conn.Close()
}

func (m *mainter) setTimeout(timeout string) time.Duration {
	duration, _ := time.ParseDuration(timeout)

	return duration
}
