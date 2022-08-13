package dag

import (
	"github.com/pipego/dag/runner"
)

type Task struct {
	Name     string
	File     runner.File
	Commands []string
	Depends  []string
}

type Dag struct {
	Vertex []Vertex
	Edge   []Edge
}

type Vertex struct {
	Name     string
	File     runner.File
	Commands []string
}

type Edge struct {
	From string
	To   string
}
