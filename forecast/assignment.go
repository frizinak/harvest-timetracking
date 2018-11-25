package forecast

import (
	"net/url"
	"strconv"
	"time"

	"github.com/frizinak/harvest-timetracking/harvest"
)

type AssignmentsParams struct {
	ProjectID *int
	PersonID  *int
	StartDate *time.Time
	EndDate   *time.Time
}

func (a *AssignmentsParams) Values() url.Values {
	v := make(url.Values)
	if a.PersonID != nil {
		v.Set("person_id", strconv.Itoa(*a.PersonID))
	}
	if a.ProjectID != nil {
		v.Set("project_id", strconv.Itoa(*a.ProjectID))
	}
	if a.StartDate != nil {
		v.Set("start_date", a.StartDate.Format(harvest.TimeFormatDateTime))
	}
	if a.EndDate != nil {
		v.Set("end_date", a.EndDate.Format(harvest.TimeFormatDateTime))
	}

	return v
}

type Assignment struct {
	ID                      int                     `json:"id"`
	StartDate               *harvest.Date           `json:"start_date"`
	EndDate                 *harvest.Date           `json:"end_date"`
	Allocation              harvest.DurationSeconds `json:"allocation"`
	Notes                   string                  `json:"notes"`
	UpdatedAt               *harvest.DateTime       `json:"updated_at"`
	UpdatedByID             int                     `json:"updated_by_id"`
	ProjectID               int                     `json:"project_id"`
	PersonID                int                     `json:"person_id"`
	PlaceholderID           int                     `json:"placeholder_id"`
	RepeatedAssignmentSetID int                     `json:"repeated_assignment_set_id"`
	ActiveOnDaysOff         bool                    `json:"active_on_days_off"`
}

type AssignmentsResponse struct {
	Assignments []*Assignment `json:"assignments"`
}

func (f *Forecast) GetAssignments(p *AssignmentsParams) (*AssignmentsResponse, error) {
	v := &AssignmentsResponse{}
	return v, f.get("/assignments", p.Values(), v)
}
