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
	Runner    Server `yaml:"runner"`
	Scheduler Server `yaml:"scheduler"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
