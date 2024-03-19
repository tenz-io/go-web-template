package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	envPort = "PORT"
)

type Config struct {
	path string    `yaml:"-"` //nolint
	Env  string    `yaml:"env" json:"env"`
	App  AppConfig `yaml:"app" json:"app"`
}

func NewConfig(path string) *Config {
	return &Config{
		path: path,
	}
}

func (c *Config) Init() error {
	log.Println("[Config] init config")
	if err := c.load(); err != nil {
		return err
	}

	// override port from env
	if port := os.Getenv(envPort); port != "" {
		c.App.Port = port
		log.Println("port from env: ", c.App.Port)
	}

	return nil
}

func (c *Config) load() error {
	if c.path == "" {
		return fmt.Errorf("config path is empty")
	}

	yamlFile, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("read config file fail, err: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("unmarshal config file fail, err: %w", err)
	}

	return nil
}

type AppConfig struct {
	Name string `yaml:"name" json:"name"`
	Port string `yaml:"port" json:"port" default:"8080"`
	Web  string `yaml:"web" json:"web" default:"web"`
}
