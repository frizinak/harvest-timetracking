package harvest

type Company struct {
	BaseURL              *URL   `json:"base_uri"`
	FullDomain           string `json:"full_domain"`
	Name                 string `json:"name"`
	Active               bool   `json:"is_active"`
	WeekStart            string `json:"week_start_day"`
	WantsTimestampTimers bool   `json:"wants_timestamp_timers"`
	TimeFormat           string `json:"time_format"`
	PlanType             string `json:"plan_type"`
	Clock                string `json:"clock"`
	Decimal              string `json:"decimal_symbol"`
	Thousands            string `json:"thousands_separator"`
	ColorScheme          string `json:"color_scheme"`
	ExpenseFeature       bool   `json:"expense_feature"`
	InvoiceFeature       bool   `json:"invoice_feature"`
	EstimateFeature      bool   `json:"estimate_feature"`
	ApprovalFeature      bool   `json:"approval_feature"`
}

func (h *Harvest) GetCompany() (*Company, error) {
	v := &Company{}
	return v, h.get("/company", nil, v)
}
