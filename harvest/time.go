package harvest

import (
	"net/url"
	"sort"
	"strconv"
	"time"
)

type TimeEntriesParams struct {
	UserID       *int
	ClientID     *int
	ProjectID    *int
	Billed       *bool
	Running      *bool
	UpdatedSince *time.Time
	From         *time.Time
	To           *time.Time
	Page         *int
	PerPage      *int
}

func (t *TimeEntriesParams) Values() url.Values {
	v := make(url.Values)
	v.Set("page", "1")

	if t.UserID != nil {
		v.Set("user_id", strconv.Itoa(*t.UserID))
	}
	if t.ClientID != nil {
		v.Set("client_id", strconv.Itoa(*t.ClientID))
	}
	if t.ProjectID != nil {
		v.Set("project_id", strconv.Itoa(*t.ProjectID))
	}
	if t.Billed != nil {
		v.Set("is_billed", boolToString(*t.Billed))
	}
	if t.Running != nil {
		v.Set("is_running", boolToString(*t.Running))
	}
	if t.UpdatedSince != nil {
		v.Set("updated_since", t.UpdatedSince.Format(timeFormatDateTime))
	}
	if t.From != nil {
		v.Set("from", t.From.Format(timeFormatDate))
	}
	if t.To != nil {
		v.Set("to", t.To.Format(timeFormatDate))
	}
	if t.Page != nil {
		v.Set("page", strconv.Itoa(*t.Page))
	}
	if t.PerPage != nil {
		v.Set("per_page", strconv.Itoa(*t.PerPage))
	}

	return v
}

type TimeEntriesResponse struct {
	NextPage     *int        `json:"next_page"`
	TotalEntries int         `json:"total_entries"`
	Page         int         `json:"page"`
	TimeEntries  TimeEntries `json:"time_entries"`
}

type Grouper func(t *TimeEntry) (key string, include bool)

type TimeEntries []*TimeEntry

func (t TimeEntries) SortSpent() TimeEntries {
	sort.SliceStable(
		t,
		func(i, j int) bool {
			ie, je := t[i], t[j]
			if ie.SpentDate == nil {
				return je.SpentDate == nil
			} else if je.SpentDate == nil {
				return ie.SpentDate != nil
			}

			return !ie.SpentDate.Before(je.SpentDate.Time)
		},
	)

	return t
}

func (t TimeEntries) Group(groupBy Grouper) Grouped {
	d := make(Grouped, 0, len(t))
	lookup := make(map[string]int)
	for _, e := range t {
		k, ok := groupBy(e)
		if !ok {
			continue
		}

		if _, ok := lookup[k]; !ok {
			lookup[k] = len(d)
			var spent time.Time
			if e.SpentDate != nil {
				spent = e.SpentDate.Time

			}
			d = append(
				d,
				&Group{
					FirstSpentDate: spent,
					SpentDates:     make([]time.Time, 0, 1),
				},
			)
		}

		group := d[lookup[k]]
		group.Hours += e.Hours.Duration
		if e.SpentDate != nil {
			group.SpentDates = append(group.SpentDates, e.SpentDate.Time)
		}
		if e.Running {
			group.Running = true
		}
	}

	return d
}

type Grouped []*Group

func (g Grouped) SortSpent() Grouped {
	sort.SliceStable(
		g,
		func(i, j int) bool {
			ie, je := g[i], g[j]
			return !ie.FirstSpentDate.Before(je.FirstSpentDate)
		},
	)

	return g
}

type Group struct {
	Running        bool
	FirstSpentDate time.Time
	SpentDates     []time.Time
	Hours          time.Duration
}

type TimeEntry struct {
	ID int `json:"id"`

	User              *UserRef           `json:"user"`
	UserAssignment    *UserAssignment    `json:"user_assignment"`
	Client            *ClientRef         `json:"client"`
	Project           *ProjectRef        `json:"project"`
	Task              *TaskRef           `json:"task"`
	TaskAssignment    *TaskAssignment    `json:"task"`
	ExternalReference *ExternalReference `json:"external_reference"`
	Invoice           *InvoiceRef        `json:"invoice"`

	Hours          DurationHours `json:"hours"`
	Notes          string        `json:"notes"`
	Locked         bool          `json:"is_locked"`
	LockedReason   string        `json:"locked_reason"`
	Closed         bool          `json:"is_closed"`
	Billed         bool          `json:"is_billed"`
	SpentDate      *Date         `json:"spent_date"`
	TimerStartedAt *DateTime     `json:"timer_started_at"`
	StartedTime    string        `json:"started_time"`
	EndedTime      string        `json:"ended_time"`
	Running        bool          `json:"is_running"`
	Billable       bool          `json:"billable"`
	Budgeted       bool          `json:"budgeted"`
	BillableRate   float64       `json:"billable_rate"`
	CostRate       float64       `json:"cost_rate"`
	CreatedAt      *DateTime     `json:"created_at"`
	UpdatedAt      *DateTime     `json:"updated_at"`
}

func (h *Harvest) GetTimeEntries(p *TimeEntriesParams) (*TimeEntriesResponse, error) {
	v := &TimeEntriesResponse{}
	return v, h.get("/time_entries", p.Values(), v)
}
