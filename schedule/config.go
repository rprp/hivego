package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type HiveConfig struct {
	Maxprocs int    `toml:"maxprocs"`
	Port     string `toml:"port"`
	Conn     string `toml:"conn"`
	Loglevel uint8  `toml:"loglevel"`
}

func LoadHiveConfig(configPath string) (config *HiveConfig) {

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println(err)
		panic(err)
		return nil
	}

	return config
}
