package runner

import (
	"context"

	"github.com/pipego/demo/config"
)

type Runner interface {
	Init(context.Context) error
	Run(context.Context) (Result, error)
}

type Config struct {
	Config config.Config
	Data   Proto
}

type runner struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Runner {
	return &runner{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *runner) Init(ctx context.Context) error {
	// TODO: Init
	return nil
}

func (r *runner) Run(ctx context.Context) (Result, error) {
	// TODO: Run
	return Result{}, nil
}
