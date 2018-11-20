package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nywillb/ftcroboscout/ftc"
)

type configuration struct {
	TOA      ftc.OrangeAllianceConfig `toml:"OrangeAlliance"`
	Database databaseConfiguration    `toml:"Database"`
}

type databaseConfiguration struct {
	Username string
	Password string
	Database string
}

var config *configuration

func configure() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}
}
