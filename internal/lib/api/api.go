package api

import (
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = "invalid status code"
)

func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func (req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if  err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %s: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	defer resp.Body.Close()

	return resp.Header.Get("Location"), nil
}