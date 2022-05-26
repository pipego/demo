package scheduler

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
	Task  Task   `json:"task"`
	Nodes []Node `json:"nodes"`
}

type Task struct {
	Name                   string   `json:"name"`
	NodeName               string   `json:"nodeName"`
	NodeSelector           []string `json:"nodeSelector"`
	RequestedResource      Resource `json:"requestedResource"`
	ToleratesUnschedulable bool     `json:"toleratesUnschedulable"`
}

type Node struct {
	Name                string   `json:"name"`
	Host                string   `json:"host"`
	Label               string   `json:"label"`
	AllocatableResource Resource `json:"allocatableResource"`
	RequestedResource   Resource `json:"requestedResource"`
	Unschedulable       bool     `json:"unschedulable"`
}

type Resource struct {
	MilliCPU int64 `json:"milliCPU"`
	Memory   int64 `json:"memory"`
	Storage  int64 `json:"storage"`
}

type Result struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}
