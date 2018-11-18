package main

import (
	"fmt"
	"os"

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

func main() {
	toa := ftc.OrangeAllianceConfig{
		APIKey:            os.Getenv("TOA_API_KEY"),
		ApplicationOrigin: os.Getenv("TOA_APPLICATION_ORIGIN"),
	}

	atomic := ftc.Team{Key: "4174"}
	quantum := ftc.Team{Key: "6051"}
	isquared := ftc.Team{Key: "6081"}

	atomic.FetchTeamDetails(&toa)
	quantum.FetchTeamDetails(&toa)
	isquared.FetchTeamDetails(&toa)

	fmt.Println(atomic)
	fmt.Println(quantum)
	fmt.Println(isquared)
}
