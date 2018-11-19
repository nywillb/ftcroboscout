package ftc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Team holds the information about an FTC team.
type Team struct {
	Key         string `json:"team_key"`
	Region      string `json:"region_key"`
	Number      int    `json:"team_number"`
	Name        string `json:"team_name_short"`
	Affiliation string `json:"team_name_long"`
	City        string `json:"city"`
	State       string `json:"state_prov"`
	ZipCode     string `json:"zip_code"`
	Country     string `json:"country"`
	Website     string `json:"website"`
	LastActive  string `json:"last_active"`
	RookieYear  int    `json:"rookie_year"`
}

// FetchTeamDetails fills in unknown information about the team
func (t *Team) FetchTeamDetails(toa *OrangeAllianceConfig) error {
	res, err := toa.MakeRequest("GET", "team/"+t.Key, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	responseArray := make([]Team, 1)

	err = json.Unmarshal(body, &responseArray)
	if err != nil {
		return err
	}

	if len(responseArray) < 1 {
		return errors.New("team not found")
	}

	*t = responseArray[0]
	return nil
}

func FetchTeams(toa *OrangeAllianceConfig) ([]Team, error) {
	res, err := toa.MakeRequest("GET", "team", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseArray := make([]Team, 1)

	err = json.Unmarshal(body, &responseArray)
	if err != nil {
		return nil, err
	}

	return responseArray, nil
}
