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
	Tasks  []Task `json:"tasks"`
	Glance Glance `json:"glance"`
}

type Task struct {
	Name     string      `json:"name"`
	File     TaskFile    `json:"file"`
	Params   []TaskParam `json:"params"`
	Commands []string    `json:"commands"`
	Log      TaskLog     `json:"log"`
	Timeout  TaskTimeout `json:"timeout"`
	Depends  []string    `json:"depends"`
}

type TaskFile struct {
	Content string `json:"content"`
	Gzip    bool   `json:"gzip"`
}

type TaskParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TaskLog struct {
	Width int64 `json:"width"`
}

type TaskTimeout struct {
	Time int64  `json:"time"`
	Unit string `json:"unit"`
}

type TaskResult struct {
	Output TaskOutput `json:"output"`
	Error  string     `json:"error"`
}

type TaskOutput struct {
	Pos     int64  `json:"pos"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

type Glance struct {
	Dir     GlanceDirReq  `json:"dir"`
	File    GlanceFileReq `json:"file"`
	Sys     GlanceSysReq  `json:"sys"`
	Timeout GlanceTimeout `json:"timeout"`
}

type GlanceDirReq struct {
	Path string `json:"path"`
}

type GlanceFileReq struct {
	Path    string `json:"path"`
	MaxSize int64  `json:"maxSize"`
}

type GlanceSysReq struct {
	Enable bool `json:"enable"`
}

type GlanceTimeout struct {
	Time int64  `json:"time"`
	Unit string `json:"unit"`
}

type GlanceReply struct {
	Dir   GlanceDirRep  `json:"dir"`
	File  GlanceFileRep `json:"file"`
	Sys   GlanceSysRep  `json:"sys"`
	Error string        `json:"error"`
}

type GlanceDirRep struct {
	Entries []GlanceEntry `json:"entries"`
}

type GlanceEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
	Time  string `json:"time"`
	User  string `json:"user"`
	Group string `json:"group"`
	Mode  string `json:"mode"`
}

type GlanceFileRep struct {
	Content  string `json:"content"`
	Readable bool   `json:"readable"`
}

type GlanceSysRep struct {
	Resource GlanceResource `json:"resource"`
	Stats    GlanceStats    `json:"stats"`
}

type GlanceResource struct {
	Allocatable GlanceAllocatable `json:"allocatable"`
	Requested   GlanceRequested   `json:"requested"`
}

type GlanceAllocatable struct {
	MilliCPU int64 `json:"milliCPU"`
	Memory   int64 `json:"memory"`
	Storage  int64 `json:"storage"`
}

type GlanceRequested struct {
	MilliCPU int64 `json:"milliCPU"`
	Memory   int64 `json:"memory"`
	Storage  int64 `json:"storage"`
}

type GlanceStats struct {
	CPU       GlanceCPU       `json:"cpu"`
	Host      string          `json:"host"`
	Memory    GlanceMemory    `json:"memory"`
	OS        string          `json:"os"`
	Storage   GlanceStorage   `json:"storage"`
	Processes []GlanceProcess `json:"processes"`
}

type GlanceCPU struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

type GlanceMemory struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

type GlanceStorage struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

type GlanceProcess struct {
	Name    string  `json:"name"`
	Cmdline string  `json:"cmdline"`
	Memory  int64   `json:"memory"`
	Percent float64 `json:"percent"`
	Pid     int64   `json:"pid"`
	Ppid    int64   `json:"ppid"`
}
