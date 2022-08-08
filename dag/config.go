package dag

import (
	"github.com/pipego/dag/runner"
)

type Task struct {
	Name     string
	File     runner.File
	Commands []string
	Depends  []string
	Timeout  runner.Timeout
}

type Dag struct {
	Vertex []Vertex
	Edge   []Edge
}

type Vertex struct {
	Name     string
	File     runner.File
	Commands []string
	Timeout  runner.Timeout
}

type Edge struct {
	From string
	To   string
}
