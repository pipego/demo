package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v3"

	"github.com/pipego/cli/config"
	"github.com/pipego/cli/pipeline"
	"github.com/pipego/cli/runner"
	"github.com/pipego/cli/scheduler"
)

var (
	app           = kingpin.New("cli", "pipego cli").Version(config.Version + "-build-" + config.Build)
	configFile    = app.Flag("config-file", "Config file (.yml)").Required().String()
	runnerFile    = app.Flag("runner-file", "Runner file (.json)").Required().String()
	schedulerFile = app.Flag("scheduler-file", "Scheduler file (.json)").Required().String()
)

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg, err := initConfig(ctx, *configFile)
	if err != nil {
		return errors.Wrap(err, "failed to init scheduler")
	}

	r, err := initRunner(ctx, cfg, *runnerFile)
	if err != nil {
		return errors.Wrap(err, "failed to init runner")
	}

	s, err := initScheduler(ctx, cfg, *schedulerFile)
	if err != nil {
		return errors.Wrap(err, "failed to init scheduler")
	}

	p, err := initPipeline(ctx, cfg, r, s)
	if err != nil {
		return errors.Wrap(err, "failed to init pipeline")
	}

	if err := runPipeline(ctx, p); err != nil {
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

func loadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return buf, nil
}

func initRunner(ctx context.Context, cfg *config.Config, name string) (runner.Runner, error) {
	c := runner.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg

	buf, err := loadFile(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load")
	}

	if err := json.Unmarshal(buf, &c.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return runner.New(ctx, c), nil
}

func initScheduler(ctx context.Context, cfg *config.Config, name string) (scheduler.Scheduler, error) {
	c := scheduler.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg

	buf, err := loadFile(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load")
	}

	if err := json.Unmarshal(buf, &c.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return scheduler.New(ctx, c), nil
}

func initPipeline(ctx context.Context, cfg *config.Config, run runner.Runner, sched scheduler.Scheduler) (pipeline.Pipeline, error) {
	c := pipeline.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Runner = run
	c.Scheduler = sched

	return pipeline.New(ctx, c), nil
}

// nolint: gosec
func runPipeline(ctx context.Context, pipe pipeline.Pipeline) error {
	if err := pipe.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init")
	}

	resScheduler, resRunner, err := pipe.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to run")
	}

	_ = pipe.Deinit(ctx)

	fmt.Println("   Run: scheduler")
	fmt.Println("  Name:", resScheduler.Name)
	fmt.Println(" Error:", resScheduler.Error)
	fmt.Println()
	fmt.Println("   Run: runner")
	for _, item := range resRunner {
		fmt.Println("Output:", item.Output)
		fmt.Println(" Error:", item.Error)
	}

	return nil
}
