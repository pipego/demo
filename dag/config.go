package dag

import (
	"github.com/pipego/dag/runner"
)

type Task struct {
	Name     string
	File     runner.File
	Params   []runner.Param
	Commands []string
	Width    int64
	Depends  []string
}

type Dag struct {
	Vertex []Vertex
	Edge   []Edge
}

type Vertex struct {
	Name     string
	File     runner.File
	Params   []runner.Param
	Commands []string
	Width    int64
}

type Edge struct {
	From string
	To   string
}
