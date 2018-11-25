package harvest

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type UsersParams struct {
	Active       *bool
	UpdatedSince *time.Time
	Page         *int
	PerPage      *int
}

func (u *UsersParams) Values() url.Values {
	v := make(url.Values)
	v.Set("page", "1")

	if u.Active != nil {
		v.Set("is_active", boolToString(*u.Active))
	}
	if u.UpdatedSince != nil {
		v.Set("updated_since", u.UpdatedSince.Format(TimeFormatDateTime))
	}
	if u.Page != nil {
		v.Set("page", strconv.Itoa(*u.Page))
	}
	if u.PerPage != nil {
		v.Set("per_page", strconv.Itoa(*u.PerPage))
	}

	return v
}

type UsersResponse struct {
	NextPage     int     `json:"next_page"`
	TotalEntries int     `json:"total_entries"`
	Page         int     `json:"page"`
	Users        []*User `json:"users"`
}

type User struct {
	ID                  int       `json:"id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	Tel                 string    `json:"telephone"`
	TZ                  string    `json:"timezone"`
	FutureProjectAccess bool      `json:"has_access_to_all_future_projects"`
	Contractor          bool      `json:"is_contractor"`
	Admin               bool      `json:"is_admin"`
	PM                  bool      `json:"is_project_manager"`
	CanSeeRates         bool      `json:"can_see_rates"`
	CanCreateProjects   bool      `json:"can_create_projects"`
	CanCreateInvoices   bool      `json:"can_create_invoices"`
	Active              bool      `json:"is_active"`
	WeeklyCapacity      int       `json:"weekly_capacity"`
	DefaultHourRate     float64   `json:"default_hourly_rate"`
	CostRate            float64   `json:"cost_rate"`
	Roles               []string  `json:"roles"`
	AvatarURL           string    `json:"avatar_url"`
	CreatedAt           *DateTime `json:"created_at"`
	UpdatedAt           *DateTime `json:"updated_at"`
}

func (u *User) Capacity() time.Duration {
	return time.Duration(u.WeeklyCapacity) * time.Second
}

func (h *Harvest) GetUsers(u *UsersParams) (*UsersResponse, error) {
	v := &UsersResponse{}
	return v, h.get("/users", u.Values(), v)
}

func (h *Harvest) GetMe() (*User, error) {
	v := &User{}
	return v, h.get("/users/me", nil, v)
}

func (h *Harvest) GetUser(id int) (*User, error) {
	v := &User{}
	return v, h.get(fmt.Sprintf("/users/%d", id), nil, v)
}
