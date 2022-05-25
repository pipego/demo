package pipeline

import (
	"context"

	"github.com/pipego/demo/config"
	"github.com/pipego/demo/runner"
	"github.com/pipego/demo/scheduler"
)

type Pipeline interface {
	Init(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config    config.Config
	Runner    runner.Runner
	Scheduler scheduler.Scheduler
}

type pipeline struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Pipeline {
	return &pipeline{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (p *pipeline) Init(ctx context.Context) error {
	// TODO: Init
	return nil
}

func (p *pipeline) Run(ctx context.Context) error {
	// TODO: Run
	return nil
}
