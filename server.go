package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var serverPaused = false

type teamResponseData struct {
	Number      int             `json:"team"`
	Name        string          `json:"name"`
	Affiliation string          `json:"affiliation"`
	City        string          `json:"city"`
	Region      string          `json:"region"`
	RookieYear  int             `json:"rookieYear"`
	FullMatch   statisticalData `json:"fullMatchData"`
	Auto        statisticalData `json:"autoData"`
	TeleOp      statisticalData `json:"teleOpData"`
	End         statisticalData `json:"endData"`
}

type statisticalData struct {
	ExpO     int     `json:"expO"`
	Variance float64 `json:"variance"`
	Opar     float64 `json:"opar"`
}

type serverResponse struct {
	Status string             `json:"status"`
	Data   []teamResponseData `json:"data"`
}

type serverErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func initalizeServer() {
	router := httprouter.New()
	router.GET("/", handleDataRequest)
	router.POST("/update", handleDataUpdate)

	fmt.Println("Now serving on port", config.Server.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Server.Port), router))
}

func handleDataRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := fetchData()
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*.roboscout.app")
	if serverPaused {
		w.WriteHeader(http.StatusServiceUnavailable)
		encoder.Encode(serverErrorResponse{Status: "error", Error: "server paused"})
	}
	encoder.Encode(serverResponse{Status: "ok", Data: data})
}

func handleDataUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*.roboscout.app")
	encoder := json.NewEncoder(w)
	if r.RemoteAddr == "127.0.0.1" {
		w.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(serverErrorResponse{Status: "error", Error: "not authorized"})
	} else {
		encoder.Encode(serverResponse{Status: "ok"})
		refresh()
	}
}
