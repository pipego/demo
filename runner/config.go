package runner

type Proto struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
	Depends  []string `json:"depends"`
}

type Result struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}
