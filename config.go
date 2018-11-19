package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nywillb/ftcroboscout/ftc"
)

type teamWithStats struct {
	*ftc.Team
	AverageScore  float64 `json:"average_score"`
	MatchesPlayed int     `json:"matches_played"`
	ExpO          float64 `json:"expO"`
	Variance      float64 `json:"variance"`
	Opar          float64 `json:"opar"`
}

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
