package ftc

import (
	"io"
	"net/http"
)

// OrangeAllianceConfig holds an configuration for the API
type OrangeAllianceConfig struct {
	APIKey            string
	ApplicationOrigin string
	Client            http.Client
}

// MakeRequest makes an API request
func (toa *OrangeAllianceConfig) MakeRequest(method string, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://theorangealliance.org/api/"+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-TOA-Key", toa.APIKey)
	req.Header.Add("X-Application-Origin", toa.ApplicationOrigin)

	res, err := toa.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
