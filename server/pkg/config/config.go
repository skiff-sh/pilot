package config

import (
	"errors"
	"log/slog"

	"github.com/skiff-sh/config"
)

type Config struct {
	Log    config.Log    `koanf:"log" yaml:"log" json:"log"`
	Server config.Server `koanf:"server" yaml:"server" json:"server"`
	Test   Test          `koanf:"test" yaml:"test" json:"test"`
}

// Test defines config that's only relevant in E2E tests.
type Test struct {
	Image      string `koanf:"image" json:"image" yaml:"image"`
	DeployName string `koanf:"deployname" json:"deployname" yaml:"deployname"`
	Namespace  string `koanf:"namespace" json:"namespace" yaml:"namespace"`
}

func New() (*Config, error) {
	k := config.InitKoanf("pilot", Default())
	conf := new(Config)
	err := k.Unmarshal("", conf)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal"), err)
	}

	return conf, nil
}

func Default() *Config {
	return &Config{
		Log: config.Log{
			Level:   slog.LevelInfo.String(),
			Outputs: "stdout",
		},
		Server: config.Server{
			Addr: ":8080",
		},
		Test: Test{
			Image:      "ghcr.io/skiff-sh/pilot:latest",
			DeployName: "pilot",
			Namespace:  "pilot",
		},
	}
}
