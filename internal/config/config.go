package config

import (
	_ "embed"

	"github.com/zedisdog/ty/config"
)

func NewConfig() (conf config.IConfig) {
	return config.NewWithBytesContent("yaml", c)
}

//go:embed config.yaml
var c []byte
