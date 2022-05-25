package config

type Config struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	MetaData   MetaData `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type MetaData struct {
	Name string `yaml:"name"`
}

type Spec struct {
	Runner    Server `json:"runner"`
	Scheduler Server `json:"scheduler"`
}

type Server struct {
	Host string
	Port int64
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
