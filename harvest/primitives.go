package harvest

import (
	"encoding/json"
	"net/url"
	"time"
)

const (
	TimeFormatDate     = "2006-01-02"
	TimeFormatDateTime = "2006-01-02T15:04:05Z07:00"
)

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func unmarshalDate(b []byte, format string) (*time.Time, error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	if s == "" {
		return nil, nil
	}

	nd, err := time.Parse(format, s)
	return &nd, err
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	nd, err := unmarshalDate(b, TimeFormatDate)
	if nd != nil {
		*d = Date{*nd}
	}
	return err
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(TimeFormatDate))
}

type DateTime struct {
	time.Time
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	nd, err := unmarshalDate(b, TimeFormatDateTime)
	if nd != nil {
		*d = DateTime{*nd}
	}
	return err
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(TimeFormatDateTime))
}

type DurationSeconds struct {
	time.Duration
}

func (d *DurationSeconds) UnmarshalJSON(b []byte) error {
	var s float64
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*d = DurationSeconds{time.Duration(s * float64(time.Second))}

	return nil
}

type DurationHours struct {
	time.Duration
}

func (d *DurationHours) UnmarshalJSON(b []byte) error {
	var s float64
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*d = DurationHours{time.Duration(s * float64(time.Hour))}

	return nil
}

type URL struct {
	url.URL
}

func (u *URL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	p, err := url.Parse(s)
	if err != nil {
		return err
	}
	*u = URL{*p}
	return nil
}

type UserRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ClientRef struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type ProjectRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UserAssignmentRef struct {
	ID             int       `json:"id"`
	ProjectManager bool      `json:"is_project_manager"`
	Active         bool      `json:"is_active"`
	Budget         Budget    `json:"budget"`
	CreatedAt      *DateTime `json:"created_at"`
	UpdatedAt      *DateTime `json:"updated_at"`
	HourlyRate     float64   `json:"hourly_rate"`
}

type TaskAssignment struct {
	ID         int       `json:"id"`
	Billable   bool      `json:"billable"`
	Active     bool      `json:"is_active"`
	CreatedAt  *DateTime `json:"created_at"`
	UpdatedAt  *DateTime `json:"updated_at"`
	HourlyRate float64   `json:"hourly_rate"`
	Budget     Budget    `json:"budget"`
	Task       TaskRef   `json:"task"`
}

type TaskRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ExternalReference struct{}

type InvoiceRef struct{}

type Budget struct{}
