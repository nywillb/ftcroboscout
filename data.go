package main

import (
	"time"

	"github.com/nywillb/ftcroboscout/ftc"
)

func importData() {

	/* Make sure event list is up to date */
	events, err := currentSeason.FetchEvents(&config.TOA)
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		rows, err := db.Query("SELECT * FROM events WHERE eventKey=?", event.Key)
		if err != nil {
			panic(err)
		}

		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		addEvent, err := tx.Prepare(
			"INSERT INTO events (eventKey, code, type, name, start, end, city, state, country, venue, website, timeZone, active, public, tournamentLevel, allianceCount, fieldCount) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		)
		if !rows.Next() {
			_, err := addEvent.Exec(
				event.Key,
				event.Code,
				event.Type,
				event.Name,
				fixTime(time.RFC3339Nano, event.StartDate),
				fixTime(time.RFC3339Nano, event.EndDate),
				event.City,
				event.State,
				event.Country,
				event.Venue,
				event.Website,
				event.TimeZone,
				fixBool(event.Active),
				fixBool(event.Public),
				event.TournamentLevel,
				event.AllianceCount,
				event.FieldCount,
			)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
		}
		rows.Close()
		tx.Commit()
	}

	/* Get all events that have ended && don't have data */
	rows, err := db.Query("SELECT eventKey FROM events WHERE `end` < NOW() AND eventKey NOT IN (SELECT eventKey FROM matches)")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var eventKey string
		rows.Scan(&eventKey)
		event := ftc.Event{Key: eventKey}
		matches, err := event.FetchMatches(&config.TOA)
		if err != nil {
			panic(err)
		}

		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		insertMatch, err := tx.Prepare(
			"INSERT INTO matches (matchKey, eventKey, tournamentLevel, name, playNumber, fieldNumber, redScore, blueScore, redPenalty, bluePenalty, redAutoScore, blueAutoScore, redTeleOpScore, blueTeleOpScore, redEndScore, blueEndScore, videoUrl)" +
				"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		)
		if err != nil {
			panic(err)
		}

		insertParticipant, err := tx.Prepare(
			"INSERT INTO match_participants (matchKey, team, isBlueAlliance) VALUES(?, ?, ?)",
		)
		if err != nil {
			panic(err)
		}

		for _, match := range matches {
			_, err := insertMatch.Exec(
				match.Key,
				match.EventKey,
				match.TournamentLevel,
				match.MatchName,
				match.PlayNumber,
				match.FieldNumber,
				match.RedScore,
				match.BlueScore,
				match.RedPenalty,
				match.BluePenalty,
				match.RedAutoScore,
				match.BlueAutoScore,
				match.RedTeleOpScore,
				match.BlueTeleOpScore,
				match.RedEndScore,
				match.BlueEndScore,
				match.VideoURL,
			)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			for _, participant := range match.Participants {
				_, err := insertParticipant.Exec(
					match.Key,
					participant.Team,
					fixBool(participant.IsBlue()),
				)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
			}
		}

		tx.Commit()
	}
	rows.Close()

	/* Refresh teams db */
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec("TRUNCATE TABLE `teams`")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	teams, err := ftc.FetchTeams(&config.TOA)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	addTeam, err := db.Prepare("INSERT INTO teams (teamKey, region, number, name, affiliation, city, state, zipCode, country, website, lastActive, rookieYear) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, team := range teams {
		_, err = addTeam.Exec(
			team.Key,
			team.Region,
			team.Number,
			team.Name,
			team.Affiliation,
			team.City,
			team.State,
			team.ZipCode,
			team.Country,
			team.Website,
			team.LastActive,
			team.RookieYear,
		)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	tx.Commit()
}

func fetchData() []teamResponseData {
	var data []teamResponseData

	rows, err := db.Query("SELECT region, number, name, affiliation, city, rookieYear, expo, variance, opar, autoExpo, autoVariance, autoOpar, teleOpExpo, teleOpVariance, teleOpOpar, endExpo, endVariance, endOpar  FROM teams WHERE expo > -999 ORDER BY expo DESC")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var team teamResponseData
		err := rows.Scan(
			&team.Region,
			&team.Number,
			&team.Name,
			&team.Affiliation,
			&team.City,
			&team.RookieYear,
			&team.FullMatch.ExpO,
			&team.FullMatch.Variance,
			&team.FullMatch.Opar,
			&team.Auto.ExpO,
			&team.Auto.Variance,
			&team.Auto.Opar,
			&team.TeleOp.ExpO,
			&team.TeleOp.Variance,
			&team.TeleOp.Opar,
			&team.End.ExpO,
			&team.End.Variance,
			&team.End.Opar,
		)
		if err != nil {
			panic(err)
		}
		data = append(data, team)
	}

	return data
}
