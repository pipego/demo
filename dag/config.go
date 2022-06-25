package dag

type Task struct {
	Name     string
	Commands []string
	Depends  []string
}

type Dag struct {
	Vertex []Vertex
	Edge   []Edge
}

type Vertex struct {
	Name string
	Run  []string
}

type Edge struct {
	From string
	To   string
}
