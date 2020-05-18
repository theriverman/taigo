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
	ID               int              `json:"id,omitempty"`
	Slug             string           `json:"slug,omitempty"`
	Name             string           `json:"name"`
	EstimatedFinish  string           `json:"estimated_finish"`
	EstimatedStart   string           `json:"estimated_start"`
	Closed           bool             `json:"closed,omitempty"`
	ClosedPoints     int              `json:"closed_points,omitempty"`
	CreatedDate      time.Time        `json:"created_date,omitempty"`
	Disponibility    float64          `json:"disponibility,omitempty"`
	ModifiedDate     time.Time        `json:"modified_date,omitempty"`
	Order            int              `json:"order,omitempty"`
	Owner            int              `json:"owner,omitempty"`
	Project          int              `json:"project"`
	ProjectExtraInfo ProjectExtraInfo `json:"project_extra_info,omitempty"`
	TotalPoints      float64          `json:"total_points,omitempty"`
	UserStories      []UserStory      `json:"user_stories,omitempty"`
}

// MilestonesQueryParams holds fields to be used as URL query parameters to filter the queried objects
type MilestonesQueryParams struct {
	Project int  `url:"project,omitempty"`
	Closed  bool `url:"closed,omitempty"`
}
