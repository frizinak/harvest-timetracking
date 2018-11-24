package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/frizinak/harvest-timetracking/harvest"
)

const dateFormat = "2006-01-02"

type Config struct {
	AccountID     string   `json:"account_id"`
	Token         string   `json:"token"`
	ExcludedDates []string `json:"exclude_dates"`
	excludedMap   map[string]struct{}
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

	return nil
}

func (c *Config) Excluded(t time.Time) bool {
	_, ok := c.excludedMap[t.Format(dateFormat)]
	return ok
}

type Timetracking struct {
	l      *log.Logger
	conf   *Config
	client *harvest.Harvest
	user   *harvest.User
}

func New(l *log.Logger, c *Config) (*Timetracking, error) {
	aid, err := strconv.Atoi(c.AccountID)
	if err != nil {
		return nil, errors.New("account_id should be a numeric value")
	}

	return &Timetracking{l: l, conf: c, client: harvest.New(aid, c.Token)}, nil
}

func (t *Timetracking) SetUID(uid int) (err error) {
	t.user = nil
	if uid == 0 {
		t.user, err = t.client.GetMe()
		return
	}

	t.user, err = t.client.GetUser(uid)
	return
}

func (t *Timetracking) User() *harvest.User {
	return t.user
}

func (t *Timetracking) GetRecentDays(
	amount int,
	from time.Time,
	actualDays bool,
) (harvest.TimeEntries, error) {
	params := &harvest.TimeEntriesParams{UserID: &t.User().ID, To: &from}

	entries := make(harvest.TimeEntries, 0, amount)
	counter := make(map[string]struct{})

	if actualDays {
		d := from
		for {
			d = shiftWeekend(d)
			for t.conf.Excluded(d) {
				d = shiftWeekend(d.AddDate(0, 0, -1))
			}
			counter[d.Format(dateFormat)] = struct{}{}
			entries = append(
				entries,
				&harvest.TimeEntry{
					Hours:     harvest.DurationHours{0},
					SpentDate: &harvest.Date{d},
				},
			)
			d = d.AddDate(0, 0, -1)
			if len(counter) == amount {
				break
			}
		}
	}

outer:
	for {
		res, err := t.client.GetTimeEntries(params)
		if err != nil {
			return nil, err
		}

		for _, e := range res.TimeEntries {
			if e.SpentDate == nil {
				continue
			}

			d := shiftWeekend(e.SpentDate.Time)
			for t.conf.Excluded(d) {
				d = shiftWeekend(d.AddDate(0, 0, -1))
			}

			df := d.Format(dateFormat)
			if _, ok := counter[df]; !ok && len(counter) == amount {
				break outer
			}
			counter[df] = struct{}{}
			entries = append(entries, e)
		}

		if res.NextPage == nil {
			break
		}

		params.Page = res.NextPage
	}

	return entries, nil
}

func shiftWeekend(d time.Time) time.Time {
	wd := d.Weekday()
	if wd == time.Sunday {
		return d.AddDate(0, 0, -2)
	} else if wd == time.Saturday {
		return d.AddDate(0, 0, -1)
	}
	return d
}
