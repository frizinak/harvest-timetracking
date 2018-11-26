package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/frizinak/harvest-timetracking/config"
)

func getConfig(l *log.Logger) (*config.ConfigLoader, *Config, error) {
	confLoader, err := config.DotFile(
		".timetracking",
		&Config{
			"-- your account id --",
			"-- your forecast account id (optional)--",
			defaultToken,
			[]string{"saturday", "sunday"},
			[]string{},
			Tasks{},
			nil,
			nil,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	conf := &Config{}
	if err := confLoader.Read(conf); err != nil {
		if os.IsNotExist(err) {
			l.Printf(
				"Config file %s does not exist, creating example. [https://id.getharvest.com/developers to create an access token]",
				confLoader.Path(),
			)
			return nil, nil, confLoader.CreateDefault()
		}

		l.Printf("Failed to parse config")
		return nil, nil, err
	}

	if conf.Token == defaultToken {
		l.Printf(
			"You should fill in your access token and account id in '%s'",
			confLoader.Path(),
		)

		return nil, nil, nil
	}

	return confLoader, conf, nil
}

type Config struct {
	AccountID         string   `json:"account_id"`
	ForecastAccountID string   `json:"forecast_account_id"`
	Token             string   `json:"token"`
	WeekdaysOff       []string `json:"weekdays_off"`
	ExcludedDates     []string `json:"exclude_dates"`
	Tasks             Tasks    `json:"tasks"`
	excludedMap       map[string]struct{}
	weekdaysOffMap    map[time.Weekday]struct{}
}

func (c *Config) Validate() error {
	c.excludedMap = make(map[string]struct{})
	for _, v := range c.ExcludedDates {
		c.excludedMap[v] = struct{}{}
		_, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
	}

	c.weekdaysOffMap = make(map[time.Weekday]struct{})
	wds := map[string]time.Weekday{
		strings.ToLower(time.Monday.String()):    time.Monday,
		strings.ToLower(time.Tuesday.String()):   time.Tuesday,
		strings.ToLower(time.Wednesday.String()): time.Wednesday,
		strings.ToLower(time.Thursday.String()):  time.Thursday,
		strings.ToLower(time.Friday.String()):    time.Friday,
		strings.ToLower(time.Saturday.String()):  time.Saturday,
		strings.ToLower(time.Sunday.String()):    time.Sunday,
	}
	for _, v := range c.WeekdaysOff {
		wd, ok := wds[strings.ToLower(v)]
		if !ok {
			return fmt.Errorf("Invalid weekday '%s'", v)
		}

		c.weekdaysOffMap[wd] = struct{}{}
	}

	if len(c.weekdaysOffMap) > 6 {
		return errors.New("What are you using this program for, if you take every day off?")
	}

	return nil
}

func (c *Config) Excluded(t time.Time) bool {
	_, ok := c.excludedMap[t.Format(dateFormat)]
	return ok
}

func (c *Config) Off(t time.Time) bool {
	_, ok := c.weekdaysOffMap[t.Weekday()]
	return ok
}

func (c *Config) AmountOff() int {
	return len(c.weekdaysOffMap)
}

func (c *Config) WorkWeek() int {
	return 7 - len(c.weekdaysOffMap)
}
