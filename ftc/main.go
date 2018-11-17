package ftc

import (
	"io"
	"net/http"
	"os"
)

func request(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, "https://theorangealliance.org/"+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-TOA-Key", os.Getenv("TOA_API_KEY"))
	req.Header.Add("X-Application-Origin", os.Getenv("TOA_APPLICATION_ORIGIN"))
	return req, nil
}
