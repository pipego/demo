package scheduler

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipego/demo/config"
	proto "github.com/pipego/demo/scheduler/proto"
)

type Scheduler interface {
	Init(context.Context) error
	Run(context.Context) (Result, error)
}

type Config struct {
	Config config.Config
	Data   Proto
}

type scheduler struct {
	cfg    *Config
	client proto.ServerProtoClient
}

func New(_ context.Context, cfg *Config) Scheduler {
	return &scheduler{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (s *scheduler) Init(_ context.Context) error {
	host := s.cfg.Config.Spec.Scheduler.Host
	port := s.cfg.Config.Spec.Scheduler.Port

	conn, err := grpc.Dial(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	defer func() { _ = conn.Close() }()

	s.client = proto.NewServerProtoClient(conn)

	return nil
}

func (s *scheduler) Run(_ context.Context) (Result, error) {
	task := func() *proto.Task {
		return &proto.Task{
			Name:         s.cfg.Data.Spec.Task.Name,
			NodeName:     s.cfg.Data.Spec.Task.NodeName,
			NodeSelector: s.cfg.Data.Spec.Task.NodeSelector,
			RequestedResource: &proto.RequestedResource{
				MilliCPU: s.cfg.Data.Spec.Task.RequestedResource.MilliCPU,
				Memory:   s.cfg.Data.Spec.Task.RequestedResource.Memory,
				Storage:  s.cfg.Data.Spec.Task.RequestedResource.Storage,
			},
			ToleratesUnschedulable: s.cfg.Data.Spec.Task.ToleratesUnschedulable,
		}
	}()

	nodes := func() []*proto.Node {
		var n []*proto.Node
		for _, item := range s.cfg.Data.Spec.Nodes {
			n = append(n, &proto.Node{
				Name:  item.Name,
				Host:  item.Host,
				Label: item.Label,
				AllocatableResource: &proto.AllocatableResource{
					MilliCPU: item.AllocatableResource.MilliCPU,
					Memory:   item.AllocatableResource.Memory,
					Storage:  item.AllocatableResource.Storage,
				},
				RequestedResource: &proto.RequestedResource{
					MilliCPU: item.RequestedResource.MilliCPU,
					Memory:   item.RequestedResource.Memory,
					Storage:  item.RequestedResource.Storage,
				},
				Unschedulable: item.Unschedulable,
			})
		}
		return n
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.Config.Spec.Scheduler.Timeout)*time.Second)
	defer cancel()

	reply, err := s.client.SendServer(ctx, &proto.ServerRequest{
		ApiVersion: s.cfg.Data.ApiVersion,
		Kind:       s.cfg.Data.Kind,
		Metadata: &proto.Metadata{
			Name: s.cfg.Data.Metadata.Name,
		},
		Spec: &proto.Spec{
			Task:  task,
			Nodes: nodes,
		},
	})

	if err != nil {
		return Result{}, errors.Wrap(err, "failed to send")
	}

	return Result{Name: reply.GetName(), Error: reply.GetError()}, nil
}
