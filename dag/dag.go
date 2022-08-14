package dag

import (
	"context"

	"github.com/pipego/cli/config"
	"github.com/pipego/dag/runner"
)

type DAG interface {
	Init(context.Context, []Task) error
	Deinit(context.Context) error
	Run(context.Context, func(string, runner.File, []string, int64, runner.Livelog) error, runner.Livelog) error
}

type Config struct {
	Config config.Config
}

type dag struct {
	cfg    *Config
	edge   []Edge
	runner runner.Runner
	vertex []Vertex
}

func New(_ context.Context, cfg *Config) DAG {
	return &dag{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (d *dag) Init(_ context.Context, tasks []Task) error {
	for _, task := range tasks {
		v := Vertex{
			Name:     task.Name,
			File:     task.File,
			Commands: task.Commands,
			Livelog:  task.Livelog,
		}
		d.vertex = append(d.vertex, v)

		for _, dep := range task.Depends {
			e := Edge{
				From: dep,
				To:   task.Name,
			}
			d.edge = append(d.edge, e)
		}
	}

	return nil
}

func (d *dag) Deinit(_ context.Context) error {
	return nil
}

func (d *dag) Run(_ context.Context, routine func(string, runner.File, []string, int64, runner.Livelog) error,
	log runner.Livelog) error {
	for _, vertex := range d.vertex {
		d.runner.AddVertex(vertex.Name, routine, vertex.File, vertex.Commands, vertex.Livelog)
	}

	for _, edge := range d.edge {
		d.runner.AddEdge(edge.From, edge.To)
	}

	return d.runner.Run(log)
}
