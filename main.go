package main

import (
	"fmt"

	"github.com/nywillb/ftcroboscout/ftc"
)

var currentSeason = ftc.Season{Key: "1819"}

func main() {
	fmt.Println(intro)
	fmt.Print("Starting up (may take some time)...")
	configure()                  //load the config and store config values
	ftc.Configure(config.TOA)    //Configure api
	initalizeDatabase()          //initalize the database
	defer deinitializeDatabase() //deinitalize the database when done using it
	refresh()                    //refresh the data and rankings
	fmt.Println("Done!")         //we're done with startup!
	initalizeServer()
}

func refresh() {
	serverPaused = true  //Pause the web server
	importData()         //Get the new data
	calculateRankings()  //Recalculate the rankings with the new data
	serverPaused = false //Unpause the web server
}
