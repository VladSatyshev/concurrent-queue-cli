package config

import "time"

const TestConfigPath = "../config/test_config.yaml"

type Config struct {
	Server ServerConfig `yaml:"Server"`
	Queues QueuesConfig `yaml:"Queues"`
	Logger LoggerConfig `yaml:"Logger"`
}

type ServerConfig struct {
	Port              string        `yaml:"Port"`
	Mode              string        `yaml:"Mode"`
	Timeout           time.Duration `yaml:"Timeout"`
	CtxDefaultTimeout time.Duration `yaml:"CtxDefaultTimeout"`
}

type QueuesConfig []QueueConfig

type QueueConfig struct {
	Name              string `yaml:"Name"`
	Length            uint   `yaml:"Length"`
	SubscribersAmount uint   `yaml:"SubscribersAmount"`
}

type LoggerConfig struct {
	Development       bool   `yaml:"Development"`
	DisableCaller     bool   `yaml:"DisableCaller"`
	DisableStacktrace bool   `yaml:"DisableStacktrace"`
	Encoding          string `yaml:"Encoding"`
	Level             string `yaml:"Level"`
}
