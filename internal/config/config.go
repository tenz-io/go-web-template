package config

type Config struct {
	Verbose bool      `yaml:"verbose" json:"verbose"`
	Env     string    `yaml:"env" json:"env"`
	App     AppConfig `yaml:"app" json:"app"`
	DB      DBConfig  `yaml:"db" json:"db"`
}

type AppConfig struct {
	Name      string `yaml:"name" json:"name"`
	Port      int    `yaml:"port" json:"port"`
	Web       string `yaml:"web" json:"web"`
	Secret    string `yaml:"secret" json:"secret"`
	AdminUser string `yaml:"admin_user" json:"admin_user"`
	AdminPass string `yaml:"admin_pass" json:"admin_pass"`
}

type DBConfig struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	User string `yaml:"user" json:"user"`
	Pass string `yaml:"pass" json:"pass"`
	DB   string `yaml:"db" json:"db"`
}
