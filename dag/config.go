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
	Language runner.Language
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
	Language runner.Language
}

type Edge struct {
	From string
	To   string
}
