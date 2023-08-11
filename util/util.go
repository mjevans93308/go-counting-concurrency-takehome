package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	DOMAIN_ENV_KEY   = "EXTERNAL_DOMAIN_ENV"
	DOMAIN_LOCALHOST = "localhost"
	ENV              = "ENV"
	ENV_TEST         = "TEST"
	ENV_STAGING      = "STAGING"
	ENV_PROD         = "PRODUCTION"
)

type API struct {
	Client  *http.Client
	baseURL string
}

var Api = NewApi("")

func NewApi(url string) *API {
	api := API{}
	if api.Client == nil {
		api.Client = &http.Client{
			Timeout: 1 * time.Second,
		}
		if url == "" {
			api.baseURL = BuildAddr()
		} else {
			api.baseURL = url
		}
	}
	return &api
}

// using `omitempty` should ensure that we avoid processing an empty response value
type ExternalResponse struct {
	Value int `json:"value,omitempty"`
}

// GetInteger makes the api request, processes the response and parses into a struct for easier handling
// takes the current ctx and current iteration var
func (api *API) GetInteger(ctx context.Context, iter int) (*ExternalResponse, error) {
	resp, err := api.Client.Get(api.baseURL + "/integers/" + fmt.Sprint(iter))
	if err != nil {
		ctx.Done()
		return nil, err
	}
	defer resp.Body.Close()
	response := &ExternalResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("could not decode response from external server for iteration %d: %s", iter, err)
	}
	return response, nil
}

func getDomain() string {
	switch os.Getenv(ENV) {
	case ENV_TEST:
		return DOMAIN_LOCALHOST + ":80"
	case ENV_STAGING, ENV_PROD:
		return os.Getenv(DOMAIN_ENV_KEY)
	default:
		return DOMAIN_LOCALHOST + ":8080"
	}
}

func getHTTPProtocol() string {
	switch os.Getenv(ENV) {
	case ENV_TEST:
		return "http://"
	case ENV_STAGING, ENV_PROD:
		return "https://"
	default:
		return "http://"
	}
}

func BuildAddr() string {
	return getHTTPProtocol() + getDomain()
}

// isEven checks whether a given value `input` is even or odd
// returns a boolean true if even, false if odd
// if `input` is blank, it will default to 0, which is an even value
func IsEven(input int) bool {
	return input%2 == 0
}
