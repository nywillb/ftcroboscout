package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nywillb/ftcroboscout/ftc"
)

type configuration struct {
	TOA      ftc.OrangeAllianceConfig `toml:"OrangeAlliance"`
	Database databaseConfiguration    `toml:"Database"`
	Server   serverConfiguration      `toml:"Server"`
}

type databaseConfiguration struct {
	Username, Password, Database string
}

type serverConfiguration struct {
	Port int
}

var config *configuration

func configure() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}
}

var intro = `
  __ _       ____       _                               _   
 / _| |_ ___|  _ \ ___ | |__   ___  ___  ___ ___  _   _| |_ 
| |_| __/ __| |_) / _ \| '_ \ / _ \/ __|/ __/ _ \| | | | __|
|  _| || (__|  _ < (_) | |_) | (_) \__ \ (_| (_) | |_| | |_ 
|_|  \__\___|_| \_\___/|_.__/ \___/|___/\___\___/ \__,_|\__|
														   
                  (c) 2018 William Barkoff
				 
`
