package harvest

import (
	"net/url"
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

type TimeEntries []*TimeEntry

func (t TimeEntries) Days(includeRunning bool) []*Day {
	d := make([]*Day, 0, len(t))
	lookup := make(map[string]int)
	for _, e := range t {
		if e.SpentDate == nil {
			continue
		}
		if !includeRunning && e.Running {
			continue
		}

		df := e.SpentDate.Format(timeFormatDate)
		if _, ok := lookup[df]; !ok {
			lookup[df] = len(d)
			d = append(d, &Day{SpentDate: e.SpentDate.Time})
		}

		day := d[lookup[df]]
		day.Hours += e.Hours.Duration
		if e.Running {
			day.Running = true
		}
	}

	return d
}

type Day struct {
	Running   bool
	SpentDate time.Time
	Hours     time.Duration
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
