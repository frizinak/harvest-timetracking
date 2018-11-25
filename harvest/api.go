package harvest

import (
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

type Api struct {
	Client          *http.Client
	AccountID       int
	Token           string
	Endpoint        string
	AccountIDHeader string
}

func (a *Api) Get(path string, query url.Values, v interface{}) error {
	u, err := url.Parse(a.Endpoint + path)
	if err != nil {
		return err
	}
	q := u.Query()
	for i := range query {
		q.Set(i, query.Get(i))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set(a.AccountIDHeader, strconv.Itoa(a.AccountID))
	req.Header.Set("Authorization", "Bearer "+a.Token)
	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		all, _ := ioutil.ReadAll(res.Body)
		return errors.New("Unexpected api error: " + string(all))
	}
	// --
	// debug, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(debug))
	// os.Exit(0)
	// --
	return json.NewDecoder(res.Body).Decode(v)
}
