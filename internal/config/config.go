package config

import "time"

type Config struct {
	Verbose bool      `yaml:"verbose" json:"verbose"`
	Env     string    `yaml:"env" json:"env"`
	App     AppConfig `yaml:"app" json:"app"`
	DB      DBConfig  `yaml:"db" json:"db"`
	JWT     JWTConfig `yaml:"jwt" json:"jwt"`
}

type AppConfig struct {
	Name string `yaml:"name" json:"name"`
	Port int    `yaml:"port" json:"port"`
	Web  string `yaml:"web" json:"web"`
	// 新增配置项
	Debug    bool   `yaml:"debug" json:"debug"`
	LogLevel string `yaml:"log_level" json:"log_level"`
	LogFile  string `yaml:"log_file" json:"log_file"`
}

type DBConfig struct {
	// SQLite 配置
	Path  string `yaml:"path" json:"path"`
	Debug bool   `yaml:"debug" json:"debug"`
	// 数据库连接池配置（SQLite 不需要，但保留兼容性）
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret" json:"secret"`
	ExpireTime time.Duration `yaml:"expire_time" json:"expire_time"`
	Issuer     string        `yaml:"issuer" json:"issuer"`
}
