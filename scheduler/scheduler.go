package scheduler

import (
	"context"

	"github.com/pipego/demo/config"
)

type Scheduler interface {
	Init(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config config.Config
}

type scheduler struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Scheduler {
	return &scheduler{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (s *scheduler) Init(ctx context.Context) error {
	// TODO: Init
	return nil
}

func (s *scheduler) Run(ctx context.Context) error {
	// TODO: Run
	return nil
}
