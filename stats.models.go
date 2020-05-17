package taigo

// DiscoverStats => https://taigaio.github.io/taiga-doc/dist/api.html#object-discover-stats
type DiscoverStats struct {
	Projects struct {
		Total int `json:"total"`
	} `json:"projects"`
}

// SystemStats => https://taigaio.github.io/taiga-doc/dist/api.html#object-system-stats
type SystemStats struct {
	Projects struct {
		AverageLastFiveWorkingDays  float64 `json:"average_last_five_working_days"`
		AverageLastSevenDays        float64 `json:"average_last_seven_days"`
		PercentWithBacklog          float64 `json:"percent_with_backlog"`
		PercentWithBacklogAndKanban float64 `json:"percent_with_backlog_and_kanban"`
		PercentWithKanban           float64 `json:"percent_with_kanban"`
		Today                       int     `json:"today"`
		Total                       int     `json:"total"`
		TotalWithBacklog            int     `json:"total_with_backlog"`
		TotalWithBacklogAndKanban   int     `json:"total_with_backlog_and_kanban"`
		TotalWithKanban             int     `json:"total_with_kanban"`
	} `json:"projects"`
	Users struct {
		AverageLastFiveWorkingDays float64        `json:"average_last_five_working_days"`
		AverageLastSevenDays       float64        `json:"average_last_seven_days"`
		CountsLastYearPerWeek      map[string]int `json:"counts_last_year_per_week"`
		Today                      int            `json:"today"`
		Total                      int            `json:"total"`
	} `json:"users"`
	Userstories struct {
		AverageLastFiveWorkingDays float64 `json:"average_last_five_working_days"`
		AverageLastSevenDays       float64 `json:"average_last_seven_days"`
		Today                      int     `json:"today"`
		Total                      int     `json:"total"`
	} `json:"userstories"`
}
