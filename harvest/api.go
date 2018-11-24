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
	client    *http.Client
	accountID int
	token     string
}

func New(accountID int, token string) *Harvest {
	return &Harvest{http.DefaultClient, accountID, token}
}

func (h *Harvest) get(path string, query url.Values, v interface{}) error {
	u, err := url.Parse(endpoint + path)
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

	req.Header.Set("Harvest-Account-ID", strconv.Itoa(h.accountID))
	req.Header.Set("Authorization", "Bearer "+h.token)
	res, err := h.client.Do(req)
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
