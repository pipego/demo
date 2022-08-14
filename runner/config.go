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
	File     File     `json:"file"`
	Commands []string `json:"commands"`
	Livelog  int64    `json:"livelog"`
	Timeout  Timeout  `json:"timeout"`
	Depends  []string `json:"depends"`
}

type File struct {
	Content string `json:"content"`
	Gzip    bool   `json:"gzip"`
}

type Timeout struct {
	Time int64  `json:"time"`
	Unit string `json:"unit"`
}

type Result struct {
	Output Output `json:"output"`
	Error  string `json:"error"`
}

type Output struct {
	Pos     int64  `json:"pos"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}
