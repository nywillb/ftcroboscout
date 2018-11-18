package ftc

import (
	"encoding/json"
	"io/ioutil"
)

// EventType holds information regaarding an event type.
type EventType struct {
	Key         string `json:"event_type_key"`
	Description string `json:"description"`
}

// Event holds information regarding an event.
type Event struct {
	Key             string `json:"event_key"`
	Season          string `json:"season_key"`
	Code            string `json:"event_code"`
	Type            string `json:"event_type_key"`
	Name            string `json:"event_name"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	City            string `json:"city"`
	State           string `json:"state_prov"`
	Country         string `json:"country"`
	Venue           string `json:"venue"`
	Website         string `json:"website"`
	TimeZone        string `json:"time_zone"`
	Active          bool   `json:"is_active"`
	Public          bool   `json:"is_public"`
	TournamentLevel int    `json:"active_tournament_level"`
	AllianceCount   int    `json:"alliance_count"`
	FieldCound      int    `json:"field_count"`
}

// FetchMatches gets all the matches and returns them in an []Match.
func (e *Event) FetchMatches(toa *OrangeAllianceConfig) ([]Match, error) {
	res, err := toa.MakeRequest("GET", "event/"+e.Key+"/matches", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	matches := make([]Match, 0)
	err = json.Unmarshal(body, &matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// FetchEvents gets all the events and returns them in an []Event.
func FetchEvents(toa *OrangeAllianceConfig) ([]Event, error) {
	res, err := toa.MakeRequest("GET", "event", nil)
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
