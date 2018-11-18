package ftc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Season holds the information regarding a season.
type Season struct {
	Key         string `json:"season_key"`
	Description string `json:"description"`
	Active      bool   `json:"is_active"`
	Events      []Event
}

// FetchSeasons gets all the seasons and returns them in a []seasons.
func FetchSeasons(toa *OrangeAllianceConfig) ([]Season, error) {
	res, err := toa.MakeRequest("GET", "api/seasons", nil)
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
func (s *Season) FetchEvents(toa *OrangeAllianceConfig) ([]Event, error) {
	reqBody := []byte(`{"season_key":"` + s.Key + `"}`)

	res, err := toa.MakeRequest("GET", "event", bytes.NewBuffer(reqBody))
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
