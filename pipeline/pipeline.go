package pipeline

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pipego/demo/config"
	"github.com/pipego/demo/runner"
	"github.com/pipego/demo/scheduler"
)

type Pipeline interface {
	Init(context.Context) error
	Deinit(context.Context) error
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
	if err := p.cfg.Scheduler.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init scheduler")
	}

	if err := p.cfg.Runner.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init runner")
	}

	return nil
}

func (p *pipeline) Deinit(ctx context.Context) error {
	return p.cfg.Scheduler.Deinit(ctx)
}

func (p *pipeline) Run(ctx context.Context) error {
	resScheduler, err := p.cfg.Scheduler.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to run scheduler")
	}

	fmt.Println("Scheduler:")
	fmt.Println("Name:", resScheduler.Name)
	fmt.Println("Error:", resScheduler.Error)

	resRunner, err := p.cfg.Runner.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to run runner")
	}

	fmt.Println("Runner:")
	fmt.Println("Message:", resRunner.Message)

	return nil
}
