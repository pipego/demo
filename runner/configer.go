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

type Configer interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (ConfigReply, error)
}

type ConfigerConfig struct {
	Config config.Config
	Data   Proto
}

type configer struct {
	cfg    *ConfigerConfig
	client proto.ServerProtoClient
	conn   *grpc.ClientConn
}

func ConfigerNew(_ context.Context, cfg *ConfigerConfig) Configer {
	return &configer{
		cfg: cfg,
	}
}

func ConfigerDefaultConfig() *ConfigerConfig {
	return &ConfigerConfig{}
}

func (c *configer) Init(ctx context.Context) error {
	if err := c.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	return nil
}

func (c *configer) Deinit(ctx context.Context) error {
	_ = c.deinitConn(ctx)

	return nil
}

func (c *configer) Run(ctx context.Context) (rep ConfigReply, err error) {
	output := func(r *proto.ConfigReply) ConfigReply {
		return ConfigReply{
			Version: r.GetVersion(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, c.setTimeout(c.cfg.Data.Spec.Config.Timeout))
	defer cancel()

	reply, e := c.client.SendConfig(ctx)
	defer func() {
		_ = reply.CloseSend()
	}()

	if e != nil {
		return rep, errors.Wrap(e, "failed to set")
	}

	if e = reply.Send(&proto.ConfigRequest{
		ApiVersion: c.cfg.Data.ApiVersion,
		Kind:       c.cfg.Data.Kind,
		Metadata: &proto.ConfigMetadata{
			Name: c.cfg.Data.Metadata.Name,
		},
		Spec: &proto.ConfigSpec{
			Config: &proto.Config{
				Version: c.cfg.Data.Spec.Config.Version,
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

func (c *configer) initConn(_ context.Context) error {
	var err error

	host := c.cfg.Config.Spec.Runner.Host
	port := c.cfg.Config.Spec.Runner.Port

	c.conn, err = grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	c.client = proto.NewServerProtoClient(c.conn)

	return nil
}

func (c *configer) deinitConn(_ context.Context) error {
	return c.conn.Close()
}

func (c *configer) setTimeout(timeout string) time.Duration {
	duration, _ := time.ParseDuration(timeout)

	return duration
}
