package ftc

import (
	"errors"
	"io"
	"net/http"
)

// OrangeAllianceConfig holds an configuration for the API
type OrangeAllianceConfig struct {
	APIKey            string
	ApplicationOrigin string
	Client            http.Client
}

var config OrangeAllianceConfig

// Configure sets up the api
func Configure(configuration OrangeAllianceConfig) {
	config = configuration
}

// MakeRequest makes an API request
func MakeRequest(method string, endpoint string, body io.Reader) (*http.Response, error) {
	if config.APIKey == "" || config.ApplicationOrigin == "" {
		return nil, errors.New("invalid TOA configuration")
	}
	req, err := http.NewRequest(method, "https://theorangealliance.org/api/"+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-TOA-Key", config.APIKey)
	req.Header.Add("X-Application-Origin", config.ApplicationOrigin)

	res, err := config.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
