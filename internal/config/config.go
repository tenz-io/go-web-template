package config

type Config struct {
	Verbose bool      `yaml:"verbose" json:"verbose"`
	Env     string    `yaml:"env" json:"env"`
	App     AppConfig `yaml:"app" json:"app"`
}

type AppConfig struct {
	Name string `yaml:"name" json:"name"`
	Port string `yaml:"port" json:"port"`
	Web  string `yaml:"web" json:"web"`
}
