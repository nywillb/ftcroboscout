package main

import (
	"github.com/nywillb/ftcroboscout/ftc"
)

var currentSeason = ftc.Season{Key: "1819"}

func main() {
	configure()                  //load the config and store config values
	ftc.Configure(config.TOA)    //Configure api
	initalizeDatabase()          //initalize the database
	defer deinitializeDatabase() //deinitalize the database when done using it
	refresh()                    //refresh the data and rankings
}

func refresh() {
	importData()        //Get the new data
	calculateRankings() //Recalculate the rankings with the new data
}
