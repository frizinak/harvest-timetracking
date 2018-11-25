package forecast

import (
	"net/http"
	"net/url"

	"github.com/frizinak/harvest-timetracking/harvest"
)

const (
	endpoint = "https://api.forecastapp.com"
)

type Forecast struct {
	api harvest.Api
}

func (f *Forecast) get(path string, query url.Values, v interface{}) error {
	return f.api.Get(path, query, v)
}

func New(accountID int, token string) *Forecast {
	return &Forecast{
		harvest.Api{
			http.DefaultClient,
			accountID,
			token,
			"https://api.forecastapp.com",
			"Forecast-Account-ID",
		},
	}
}
