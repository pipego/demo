package cmd

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v3"

	"github.com/pipego/demo/config"
	"github.com/pipego/demo/pipeline"
	"github.com/pipego/demo/runner"
	"github.com/pipego/demo/scheduler"
)

var (
	app        = kingpin.New("demo", "pipego demo").Version(config.Version + "-build-" + config.Build)
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
)

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg, err := initConfig(ctx, *configFile)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	s, err := initScheduler(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init scheduler")
	}

	r, err := initRunner(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init runner")
	}

	p, err := initPipeline(ctx, cfg, s, r)
	if err != nil {
		return errors.Wrap(err, "failed to init pipeline")
	}

	if err := runPipeline(ctx, cfg, p); err != nil {
		return errors.Wrap(err, "failed to run pipeline")
	}

	return nil
}

func initConfig(_ context.Context, name string) (*config.Config, error) {
	c := config.New()

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, _ := io.ReadAll(fi)

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

func initScheduler(ctx context.Context, cfg *config.Config) (scheduler.Scheduler, error) {
	c := scheduler.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg

	return scheduler.New(ctx, c), nil
}

func initRunner(ctx context.Context, cfg *config.Config) (runner.Runner, error) {
	c := runner.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg

	return runner.New(ctx, c), nil
}

func initPipeline(ctx context.Context, cfg *config.Config, sched scheduler.Scheduler, run runner.Runner) (pipeline.Pipeline, error) {
	c := pipeline.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Runner = run
	c.Scheduler = sched

	return pipeline.New(ctx, c), nil
}

func runPipeline(ctx context.Context, cfg *config.Config, pipe pipeline.Pipeline) error {
	// TODO: runPipeline
	return nil
}
