package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initalizeDatabase() {
	var err error
	db, err = sql.Open("mysql", config.Database.Username+":"+config.Database.Password+"@/"+config.Database.Database+"?charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func deinitializeDatabase() {
	db.Close()
}

func fixTime(format string, date string) string {
	timeObj, err := time.Parse(format, date)
	if err != nil {
		panic(err)
	}
	return timeObj.Format("2006-01-02 15:04:05")
}

func fixBool(b bool) int {
	if b {
		return 1
	}
	return 0
}
