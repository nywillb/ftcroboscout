package ftc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Season holds the information regarding a season.
type Season struct {
	Key         string `json:"season_key"`
	Description string `json:"description"`
	Active      bool   `json:"is_active"`
	Events      []Event
}

// FetchSeasons gets all the seasons and returns them in a []seasons.
func FetchSeasons() ([]Season, error) {
	client := &http.Client{}

	req, err := request("GET", "api/seasons", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	seasons := make([]Season, 0)
	err = json.Unmarshal(body, &seasons)
	if err != nil {
		return nil, err
	}
	return seasons, nil
}

// FetchEvents gets all the events for a praticular season returns them in a []Event.
func (s *Season) FetchEvents() ([]Event, error) {
	client := &http.Client{}

	reqBody := []byte(`{"season_key":"` + s.Key + `"}`)

	req, err := request("GET", "api/event", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
