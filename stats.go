package main

import (
	"fmt"

	"github.com/montanaflynn/stats"
	"github.com/nywillb/ftcroboscout/ftc"
)

type teamStatistics struct {
	Team         int
	Scores       []int
	AutoScores   []int
	TeleOpScores []int
	EndScores    []int
	A            int
	ExpO         int
	Var          float64
	Opar         float64
	AAuto        int
	AutoExpO     int
	AutoVar      float64
	AutoOpar     float64
	ATeleOp      int
	TeleOpExpO   int
	TeleOpVar    float64
	TeleOpOpar   float64
	AEnd         int
	EndExpO      int
	EndVar       float64
	EndOpar      float64
	Teammates    []int
}

type matchParticipant struct {
	Key      string
	Team     int
	Alliance string
}

func calculateRankings() {
	statistics := makeStatisticsMap()
	statistics = generateStats(statistics)
	saveStats(statistics)
}

func makeStatisticsMap() map[int]teamStatistics {
	statistics := map[int]teamStatistics{}

	matchRows, err := db.Query("SELECT tournamentLevel, matchKey, redScore, blueScore, redAutoScore, blueAutoScore, redTeleOpScore, blueTeleOpScore, redEndScore, blueEndScore FROM matches")
	if err != nil {
		panic(err)
	}
	defer matchRows.Close()

	for matchRows.Next() {
		match := ftc.Match{}
		var level int
		err := matchRows.Scan(
			&level,
			&match.Key,
			&match.RedScore,
			&match.BlueScore,
			&match.RedAutoScore,
			&match.BlueAutoScore,
			&match.RedTeleOpScore,
			&match.BlueTeleOpScore,
			&match.RedEndScore,
			&match.BlueEndScore,
		)
		if err != nil {
			panic(err)
		}

		if level > 1 {
			continue //don't process matches with more than 2 players, roboscout can't handle these
		}

		participantRows, err := db.Query("SELECT team, isBlueAlliance FROM match_participants WHERE matchKey=?", match.Key)
		if err != nil {
			panic(err)
		}
		defer participantRows.Close()

		redAlliance := []teamStatistics{}
		blueAlliance := []teamStatistics{}

		for participantRows.Next() {
			participant := matchParticipant{}
			var alliance int
			participantRows.Scan(&participant.Team, &alliance)
			if alliance == 0 {
				participant.Alliance = "red"
			} else {
				participant.Alliance = "blue"
			}

			thisTeamStatistics := statistics[participant.Team]
			thisTeamStatistics.Team = participant.Team

			if participant.Alliance == "blue" {
				thisTeamStatistics.AutoScores = append(thisTeamStatistics.AutoScores, match.BlueAutoScore)
				thisTeamStatistics.TeleOpScores = append(thisTeamStatistics.TeleOpScores, match.BlueTeleOpScore)
				thisTeamStatistics.EndScores = append(thisTeamStatistics.EndScores, match.BlueEndScore)
				thisTeamStatistics.Scores = append(thisTeamStatistics.Scores, match.BlueAutoScore+match.BlueTeleOpScore+match.BlueEndScore)
				blueAlliance = append(blueAlliance, thisTeamStatistics)
			} else {
				thisTeamStatistics.AutoScores = append(thisTeamStatistics.AutoScores, match.RedAutoScore)
				thisTeamStatistics.TeleOpScores = append(thisTeamStatistics.TeleOpScores, match.RedTeleOpScore)
				thisTeamStatistics.EndScores = append(thisTeamStatistics.EndScores, match.RedEndScore)
				thisTeamStatistics.Scores = append(thisTeamStatistics.Scores, match.RedAutoScore+match.RedTeleOpScore+match.RedEndScore)
				redAlliance = append(redAlliance, thisTeamStatistics)
			}
		}

		redAlliance = addTeamates(redAlliance)
		blueAlliance = addTeamates(blueAlliance)

		for _, v := range redAlliance {
			statistics[v.Team] = v
		}

		for _, v := range blueAlliance {
			statistics[v.Team] = v
		}
	}
	return statistics
}

func addTeamates(alliance []teamStatistics) []teamStatistics {
	for i, v := range alliance {
		for _, v2 := range alliance {
			if v.Team == v2.Team {
				continue
			}
			alliance[i].Teammates = append(v.Teammates, v2.Team)
		}
	}
	return alliance
}

func generateStats(statistics map[int]teamStatistics) map[int]teamStatistics {
	// Compute all the averages
	for i, v := range statistics {
		v.A = average(v.Scores)
		v.AAuto = average(v.AutoScores)
		v.ATeleOp = average(v.TeleOpScores)
		v.AEnd = average(v.EndScores)
		statistics[i] = v
	}

	// Compute ExpO and Variance
	var expOs, autoExpOs, teleOpExpOs, endExpOs []int

	for i, v := range statistics {
		var tmA []int //Teammate averages
		var tmAAuto []int
		var tmATeleOp []int
		var tmAEnd []int
		for _, w := range v.Teammates {
			t := statistics[w]

			tmA = append(tmA, t.A)
			tmAAuto = append(tmAAuto, t.AAuto)
			tmATeleOp = append(tmATeleOp, t.ATeleOp)
			tmAEnd = append(tmAEnd, t.AEnd)
		}
		v.ExpO = v.A - (average(tmA) / 2)
		v.AutoExpO = v.AAuto - (average(tmAAuto) / 2)
		v.TeleOpExpO = v.ATeleOp - (average(tmATeleOp) / 2)
		v.EndExpO = v.AEnd - (average(tmAEnd) / 2)

		expOs = append(expOs, v.ExpO)
		autoExpOs = append(autoExpOs, v.AutoExpO)
		teleOpExpOs = append(teleOpExpOs, v.TeleOpExpO)
		endExpOs = append(endExpOs, v.EndExpO)

		matchStats := stats.LoadRawData(v.Scores)
		autoStats := stats.LoadRawData(v.AutoScores)
		teleOpStats := stats.LoadRawData(v.TeleOpScores)
		endStats := stats.LoadRawData(v.EndScores)

		variance, err := matchStats.SampleVariance()
		autoVar, err := autoStats.SampleVariance()
		teleOpVar, err := teleOpStats.SampleVariance()
		endVar, err := endStats.SampleVariance()
		if err != nil {
			panic(err)
		}

		v.Var = variance
		v.AutoVar = autoVar
		v.TeleOpVar = teleOpVar
		v.EndVar = endVar

		statistics[i] = v
	}

	avgExpO := average(expOs)
	avgAutoExpO := average(autoExpOs)
	avgTeleOpExpO := average(teleOpExpOs)
	avgEndExpO := average(endExpOs)

	for i, v := range statistics {
		v.Opar = float64(v.ExpO) / float64(avgExpO)
		v.AutoOpar = float64(v.AutoExpO) / float64(avgAutoExpO)
		v.TeleOpOpar = float64(v.TeleOpExpO) / float64(avgTeleOpExpO)
		v.EndOpar = float64(v.EndExpO) / float64(avgEndExpO)

		fmt.Println(v.Opar)

		statistics[i] = v
	}

	return statistics
}

func saveStats(statistics map[int]teamStatistics) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("UPDATE `ftcroboscout`.`teams` SET `autoOpar` = ?, `endExpo` = ?, `teleOpOpar` = ?, `opar` = ?, `autoExpo` = ?, `variance` = ?, `endVariance` = ?, `teleOpExpo` = ?, `autoVariance` = ?, `teleOpVariance` = ?, `endOpar` = ?, `expo` = ? WHERE `number` = ?")

	for i, v := range statistics {
		_, err := stmt.Exec(
			v.AutoOpar,
			v.EndExpO,
			v.TeleOpOpar,
			v.Opar,
			v.AutoExpO,
			v.Var,
			v.EndVar,
			v.TeleOpExpO,
			v.AutoVar,
			v.TeleOpVar,
			v.EndOpar,
			v.ExpO,
			i,
		)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	tx.Commit()
}

func average(n []int) int {
	var tot int
	for _, v := range n {
		tot += v
	}
	return tot / len(n)
}
