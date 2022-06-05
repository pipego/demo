package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/pipego/cli/config"
	"github.com/pipego/cli/runner"
	"github.com/pipego/cli/scheduler"
)

type Pipeline interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (scheduler.Result, runner.Result, error)
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
	_ = p.cfg.Runner.Deinit(ctx)
	_ = p.cfg.Scheduler.Deinit(ctx)

	return nil
}

func (p *pipeline) Run(ctx context.Context) (s scheduler.Result, r runner.Result, e error) {
	resScheduler, err := p.cfg.Scheduler.Run(ctx)
	if err != nil {
		return scheduler.Result{}, runner.Result{}, errors.Wrap(err, "failed to issuerail scheduler")
	}

	resRunner, err := p.cfg.Runner.Run(ctx)
	if err != nil {
		return scheduler.Result{}, runner.Result{}, errors.Wrap(err, "failed to run runner")
	}

	return resScheduler, resRunner, nil
}
