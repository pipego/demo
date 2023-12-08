package dag

import (
	"context"

	"github.com/pipego/cli/config"
	"github.com/pipego/dag/runner"
)

type DAG interface {
	Init(context.Context, []Task) error
	Deinit(context.Context) error
	Run(context.Context, func(string, runner.File, []runner.Param, []string, int64, int64, runner.Livelog) error, runner.Livelog) error
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
	for index := range tasks {
		v := Vertex{
			Name:     tasks[index].Name,
			File:     tasks[index].File,
			Params:   tasks[index].Params,
			Commands: tasks[index].Commands,
			Count:    tasks[index].Count,
			Width:    tasks[index].Width,
		}
		d.vertex = append(d.vertex, v)

		for _, dep := range tasks[index].Depends {
			e := Edge{
				From: dep,
				To:   tasks[index].Name,
			}
			d.edge = append(d.edge, e)
		}
	}

	return nil
}

func (d *dag) Deinit(_ context.Context) error {
	return nil
}

func (d *dag) Run(_ context.Context, routine func(string, runner.File, []runner.Param, []string, int64, int64, runner.Livelog) error,
	log runner.Livelog) error {
	for _, vertex := range d.vertex {
		d.runner.AddVertex(vertex.Name, routine, vertex.File, vertex.Params, vertex.Commands, vertex.Count, vertex.Width)
	}

	for _, edge := range d.edge {
		d.runner.AddEdge(edge.From, edge.To)
	}

	return d.runner.Run(log)
}
