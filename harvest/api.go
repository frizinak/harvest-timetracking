package harvest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	endpoint = "https://api.harvestapp.com/v2"
)

type Harvest struct {
	api Api
}

func New(accountID int, token string) *Harvest {
	return &Harvest{
		Api{
			http.DefaultClient,
			accountID,
			token,
			"https://api.harvestapp.com/v2",
			"Harvest-Account-ID",
		},
	}
}

func (h *Harvest) get(path string, query url.Values, v interface{}) error {
	return h.api.Get(path, query, v)
}

func (h *Harvest) post(path string, query url.Values, body interface{}, v interface{}) error {
	return h.api.Post(path, query, body, v)
}

type Api struct {
	Client          *http.Client
	AccountID       int
	Token           string
	Endpoint        string
	AccountIDHeader string
}

func (a *Api) Get(path string, query url.Values, v interface{}) error {
	req, err := a.prepareRequest(path, query)
	if err != nil {
		return err
	}

	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		all, _ := ioutil.ReadAll(res.Body)
		return errors.New("Unexpected api error: " + string(all))
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func (a *Api) Post(path string, query url.Values, body interface{}, v interface{}) error {
	req, err := a.prepareRequest(path, query)
	if err != nil {
		return err
	}

	var rw bytes.Buffer
	e := json.NewEncoder(&rw)
	if err := e.Encode(body); err != nil {
		return err
	}

	req.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(&rw)

	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		all, _ := ioutil.ReadAll(res.Body)
		return errors.New("Unexpected api error: " + string(all))
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func (a *Api) prepareRequest(path string, query url.Values) (*http.Request, error) {
	u, err := url.Parse(a.Endpoint + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for i := range query {
		q.Set(i, query.Get(i))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(a.AccountIDHeader, strconv.Itoa(a.AccountID))
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return req, nil
}
