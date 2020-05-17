package taigo

import "time"

// TODO
/*

	Find a way to fetch and store custom Taiga headers:
	  - Taiga-Info-Total-Opened-Milestones
	  - Taiga-Info-Total-Closed-Milestones

*/

// Milestone represents all fields of a Milestone(Sprint)
//
// https://taigaio.github.io/taiga-doc/dist/api.html#object-milestone-detail
type Milestone struct {
	Closed           bool             `json:"closed"`
	ClosedPoints     int              `json:"closed_points"`
	CreatedDate      time.Time        `json:"created_date"`
	Disponibility    float64          `json:"disponibility"`
	EstimatedFinish  string           `json:"estimated_finish"`
	EstimatedStart   string           `json:"estimated_start"`
	ID               int              `json:"id"`
	ModifiedDate     time.Time        `json:"modified_date"`
	Name             string           `json:"name"`
	Order            int              `json:"order"`
	Owner            int              `json:"owner"`
	Project          int              `json:"project"`
	ProjectExtraInfo ProjectExtraInfo `json:"project_extra_info"`
	Slug             string           `json:"slug"`
	TotalPoints      float64          `json:"total_points"`
	UserStories      []UserStory      `json:"user_stories"`
}

// MilestonesQueryParams holds fields to be used as URL query parameters to filter the queried objects
type MilestonesQueryParams struct {
	Project int  `url:"project,omitempty"`
	Closed  bool `url:"closed,omitempty"`
}
