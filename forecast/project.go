package forecast

import "github.com/frizinak/harvest-timetracking/harvest"

type Project struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Color       string            `json:"color"`
	Code        string            `json:"code"`
	Notes       string            `json:"notes"`
	StartDate   *harvest.Date     `json:"start_date"`
	EndDate     *harvest.Date     `json:"end_date"`
	HarvestID   int               `json:"harvest_id"`
	Archived    bool              `json:"archived"`
	UpdatedAt   *harvest.DateTime `json:"updated_at"`
	UpdatedByID int               `json:"updated_by_id"`
}

type ProjectsResponse struct {
	Projects []*Project `json:"projects"`
}

func (f *Forecast) GetProjects() (*ProjectsResponse, error) {
	v := &ProjectsResponse{}
	return v, f.get("/projects", nil, v)
}
