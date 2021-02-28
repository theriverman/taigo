package taigo

import (
	"net/http"
	"strconv"
	"time"
)

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
	ID                 int              `json:"id,omitempty"`
	Slug               string           `json:"slug,omitempty"`
	Name               string           `json:"name"`
	EstimatedFinish    string           `json:"estimated_finish"`
	EstimatedStart     string           `json:"estimated_start"`
	Closed             bool             `json:"closed,omitempty"`
	ClosedPoints       float64          `json:"closed_points,omitempty"`
	CreatedDate        time.Time        `json:"created_date,omitempty"`
	Disponibility      float64          `json:"disponibility,omitempty"`
	ModifiedDate       time.Time        `json:"modified_date,omitempty"`
	Order              int              `json:"order,omitempty"`
	Owner              int              `json:"owner,omitempty"`
	Project            int              `json:"project"`
	ProjectExtraInfo   ProjectExtraInfo `json:"project_extra_info,omitempty"`
	TotalPoints        float64          `json:"total_points,omitempty"`
	UserStories        []UserStory      `json:"user_stories,omitempty"`
	IncludeAttachments bool             `url:"include_attachments,omitempty"`
}

// MilestonesQueryParams holds fields to be used as URL query parameters to filter the queried objects
type MilestonesQueryParams struct {
	Project int  `url:"project,omitempty"`
	Closed  bool `url:"closed,omitempty"`
}

// MilestoneTotalInfo holds the two extra headers returned by Taiga when filtering for milestones
//
// Taiga-Info-Total-Opened-Milestones: the number of opened milestones for this project
// Taiga-Info-Total-Closed-Milestones: the number of closed milestones for this project
//
// https://taigaio.github.io/taiga-doc/dist/api.html#milestones-list
type MilestoneTotalInfo struct {
	TaigaInfoTotalOpenedMilestones int // Taiga-Info-Total-Opened-Milestones
	TaigaInfoTotalClosedMilestones int // Taiga-Info-Total-Closed-Milestones
}

// LoadFromHeaders accepts an *http.Response struct and reads the relevant
// pagination headers returned by Taiga
func (m *MilestoneTotalInfo) LoadFromHeaders(response *http.Response) {
	m.TaigaInfoTotalOpenedMilestones, _ = strconv.Atoi(response.Header.Get("Taiga-Info-Total-Opened-Milestones"))
	m.TaigaInfoTotalClosedMilestones, _ = strconv.Atoi(response.Header.Get("Taiga-Info-Total-Closed-Milestones"))
}
