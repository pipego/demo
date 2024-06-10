package dag

import (
	"context"

	"github.com/pipego/cli/config"
	"github.com/pipego/dag/runner"
)

type DAG interface {
	Init(context.Context, []Task) error
	Deinit(context.Context) error
	Run(context.Context, func(string, runner.File, []runner.Param, []string, int64, runner.Language, runner.Log) error, runner.Log) error
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
			Width:    tasks[index].Width,
			Language: tasks[index].Language,
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

func (d *dag) Run(_ context.Context, routine func(string, runner.File, []runner.Param, []string, int64, runner.Language, runner.Log) error,
	log runner.Log) error {
	for i := range d.vertex {
		d.runner.AddVertex(d.vertex[i].Name, routine, d.vertex[i].File, d.vertex[i].Params, d.vertex[i].Commands,
			d.vertex[i].Width, d.vertex[i].Language)
	}

	for _, edge := range d.edge {
		d.runner.AddEdge(edge.From, edge.To)
	}

	return d.runner.Run(log)
}
