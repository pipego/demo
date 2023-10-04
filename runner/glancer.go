package runner

import (
	"context"

	"github.com/pipego/cli/config"
)

type Glancer interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Dir(context.Context, string) ([]GlanceEntry, error)
	File(context.Context, string, int64) (string, bool, error)
	Sys(context.Context) (GlanceAllocatable, GlanceRequested, GlanceCPU, GlanceMemory, GlanceStorage, string, string, error)
}

type GlancerConfig struct {
	Config config.Config
}

type glancer struct {
	cfg *GlancerConfig
}

func GlancerNew(_ context.Context, cfg *GlancerConfig) Glancer {
	return &glancer{
		cfg: cfg,
	}
}

func GlancerDefaultConfig() *GlancerConfig {
	return &GlancerConfig{}
}

func (g *glancer) Init(ctx context.Context) error {
	return nil
}

func (g *glancer) Deinit(ctx context.Context) error {
	return nil
}

func (g *glancer) Dir(ctx context.Context, path string) (entries []GlanceEntry, err error) {
	return entries, err
}

func (g *glancer) File(ctx context.Context, path string, maxSize int64) (content string, readable bool, err error) {
	return content, readable, err
}

// nolint: gocritic,lll
func (g *glancer) Sys(ctx context.Context) (allocatable GlanceAllocatable, requested GlanceRequested, _cpu GlanceCPU, _memory GlanceMemory, _storage GlanceStorage, _host, _os string, err error) {
	return allocatable, requested, _cpu, _memory, _storage, _host, _os, err
}
