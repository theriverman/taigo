package taigo

import (
	"time"
)

func genericToEpic(anyEpicObject interface{}) *Epic {
	payloadEpic := Epic{}
	convertStructViaJSON(&anyEpicObject, &payloadEpic)
	return &payloadEpic
}

func genericToEpics(anyEpicObjectSlice interface{}) []Epic {
	payloadEpicsSlice := []Epic{}
	convertStructViaJSON(&anyEpicObjectSlice, &payloadEpicsSlice)
	return payloadEpicsSlice
}

// Epic represents the mandatory fields of an Epic only
type Epic struct {
	TaigaBaseObject
	ID                int      `json:"id,omitempty"`
	Ref               int      `json:"ref,omitempty"`
	Version           int      `json:"version,omitempty"`
	AssignedTo        int      `json:"assigned_to,omitempty"`
	BlockedNote       string   `json:"blocked_note,omitempty"`
	ClientRequirement bool     `json:"client_requirement,omitempty"`
	Color             string   `json:"color,omitempty"`
	Description       string   `json:"description,omitempty"`
	EpicsOrder        int64    `json:"epics_order,omitempty"`
	IsBlocked         bool     `json:"is_blocked,omitempty"`
	Project           int      `json:"project,omitempty"`
	Status            int      `json:"status,omitempty"`
	Subject           string   `json:"subject,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	TeamRequirement   bool     `json:"team_requirement,omitempty"`
	Watchers          []int    `json:"watchers,omitempty"`
	EpicDetail        *EpicDetail
	EpicDetailGET     *EpicDetailGET
	EpicDetailLIST    *EpicDetailLIST
}

// GetID returns the ID
func (e *Epic) GetID() int {
	return e.ID
}

// GetRef returns the Ref
func (e *Epic) GetRef() int {
	return e.Ref
}

// GetVersion return the version
func (e *Epic) GetVersion() int {
	return e.Version
}

// GetSubject returns the subject
func (e *Epic) GetSubject() string {
	return e.Subject
}

// GetProject returns the project ID
func (e *Epic) GetProject() int {
	return e.Project
}

// ListRelatedUserStories => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-list
func (e *Epic) ListRelatedUserStories(client *Client) ([]EpicRelatedUserStoryDetail, error) {
	return client.Epic.ListRelatedUserStories(e.ID)
}

// EpicDetailLIST -> Epic detail (LIST)
// https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-detail-list
type EpicDetailLIST []struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments         []interface{}       `json:"attachments,omitempty"`
	BlockedNote         string              `json:"blocked_note,omitempty"`
	ClientRequirement   bool                `json:"client_requirement,omitempty"`
	Color               string              `json:"color,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	EpicsOrder          int64               `json:"epics_order,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsBlocked           bool                `json:"is_blocked,omitempty"`
	IsClosed            bool                `json:"is_closed,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	ModifiedDate        time.Time           `json:"modified_date,omitempty"`
	Owner               int                 `json:"owner,omitempty"`
	OwnerExtraInfo      OwnerExtraInfo      `json:"owner_extra_info,omitempty"`
	Project             int                 `json:"project,omitempty"`
	ProjectExtraInfo    Project             `json:"project_extra_info,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Status              int                 `json:"status,omitempty"`
	StatusExtraInfo     StatusExtraInfo     `json:"status_extra_info,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	Tags                Tags                `json:"tags,omitempty"`
	TeamRequirement     bool                `json:"team_requirement,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	UserStoriesCounts   UserStoriesCounts   `json:"user_stories_counts,omitempty"`
	Version             int                 `json:"version,omitempty"`
	Watchers            []int               `json:"watchers,omitempty"`
}

// AsEpics packs the returned EpicDetailLIST into a generic Epic struct
func (e *EpicDetailLIST) AsEpics() ([]Epic, error) {
	epics := genericToEpics(&e)
	for i := 0; i < len(epics); i++ {
		epics[i].EpicDetailLIST = e
	}
	return epics, nil
}

// EpicDetailGET => Epic detail (GET) https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-detail-get
type EpicDetailGET struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments         []interface{}       `json:"attachments,omitempty"`
	BlockedNote         string              `json:"blocked_note,omitempty"`
	BlockedNoteHTML     string              `json:"blocked_note_html,omitempty"`
	ClientRequirement   bool                `json:"client_requirement,omitempty"`
	Color               string              `json:"color,omitempty"`
	Comment             string              `json:"comment,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	Description         string              `json:"description,omitempty"`
	DescriptionHTML     string              `json:"description_html,omitempty"`
	EpicsOrder          int64               `json:"epics_order,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsBlocked           bool                `json:"is_blocked,omitempty"`
	IsClosed            bool                `json:"is_closed,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	ModifiedDate        time.Time           `json:"modified_date,omitempty"`
	Neighbors           Neighbors           `json:"neighbors,omitempty"`
	Owner               int                 `json:"owner,omitempty"`
	OwnerExtraInfo      OwnerExtraInfo      `json:"owner_extra_info,omitempty"`
	Project             int                 `json:"project,omitempty"`
	ProjectExtraInfo    ProjectExtraInfo    `json:"project_extra_info,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Status              int                 `json:"status,omitempty"`
	StatusExtraInfo     StatusExtraInfo     `json:"status_extra_info,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	Tags                Tags                `json:"tags,omitempty"`
	TeamRequirement     bool                `json:"team_requirement,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	UserStoriesCounts   UserStoriesCounts   `json:"user_stories_counts,omitempty"`
	Version             int                 `json:"version,omitempty"`
	Watchers            []int               `json:"watchers,omitempty"`
}

// AsEpic packs the returned EpicDetailGET into a generic Epic struct
func (e *EpicDetailGET) AsEpic() (*Epic, error) {
	epic := genericToEpic(&e)
	epic.EpicDetailGET = e // Backmap original EpicDetailLIST
	return epic, nil
}

// EpicDetail => Epic detail https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-detail
type EpicDetail struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments         []interface{}       `json:"attachments,omitempty"`
	BlockedNote         string              `json:"blocked_note,omitempty"`
	BlockedNoteHTML     string              `json:"blocked_note_html,omitempty"`
	ClientRequirement   bool                `json:"client_requirement,omitempty"`
	Color               string              `json:"color,omitempty"`
	Comment             string              `json:"comment,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	Description         string              `json:"description,omitempty"`
	DescriptionHTML     string              `json:"description_html,omitempty"`
	EpicsOrder          int64               `json:"epics_order,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsBlocked           bool                `json:"is_blocked,omitempty"`
	IsClosed            bool                `json:"is_closed,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	ModifiedDate        time.Time           `json:"modified_date,omitempty"`
	Neighbors           Neighbors           `json:"neighbors,omitempty"`
	Owner               int                 `json:"owner,omitempty"`
	OwnerExtraInfo      OwnerExtraInfo      `json:"owner_extra_info,omitempty"`
	Project             int                 `json:"project"` // Mandatory
	ProjectExtraInfo    ProjectExtraInfo    `json:"project_extra_info,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Status              int                 `json:"status,omitempty"`
	StatusExtraInfo     StatusExtraInfo     `json:"status_extra_info,omitempty"`
	Subject             string              `json:"subject"` // Mandatory
	Tags                [][]string          `json:"tags,omitempty"`
	TeamRequirement     bool                `json:"team_requirement,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	UserStoriesCounts   UserStoriesCounts   `json:"user_stories_counts,omitempty"`
	Version             int                 `json:"version,omitempty"`
	Watchers            []int               `json:"watchers,omitempty"`
}

// AsEpic packs the returned EpicDetail into a generic Epic struct
func (e *EpicDetail) AsEpic() (*Epic, error) {
	epic := genericToEpic(&e)
	epic.EpicDetail = e
	return epic, nil
}

// EpicFiltersDataDetail => Epic filters data detail https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-filters-data
type EpicFiltersDataDetail struct {
	AssignedTo []struct {
		Count    int    `json:"count,omitempty"`
		FullName string `json:"full_name,omitempty"`
		ID       int    `json:"id,omitempty"`
	} `json:"assigned_to,omitempty"`
	Owners []struct {
		Count    int    `json:"count,omitempty"`
		FullName string `json:"full_name,omitempty"`
		ID       int    `json:"id,omitempty"`
	} `json:"owners,omitempty"`
	Statuses []struct {
		Color string `json:"color,omitempty"`
		Count int    `json:"count,omitempty"`
		ID    int    `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Order int    `json:"order,omitempty"`
	} `json:"statuses,omitempty"`
	Tags []struct {
		Color TagsColors `json:"color,omitempty"`
		Count int        `json:"count,omitempty"`
		Name  string     `json:"name,omitempty"`
	} `json:"tags,omitempty"`
}

// EpicRelatedUserStoryDetail => Epic related user story detail https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-related-user-story-detail
type EpicRelatedUserStoryDetail struct {
	EpicID      int   `json:"epic,omitempty"`
	Order       int64 `json:"order,omitempty"`
	UserStoryID int   `json:"user_story,omitempty"`
}

// GetUserStory returns the UserStory referred in the EpicRelatedUserStoryDetail
func (e *EpicRelatedUserStoryDetail) GetUserStory(c *Client) (*UserStory, error) {
	return c.UserStory.Get(e.UserStoryID)
}

// GetEpic returns the Epic referred in the EpicRelatedUserStoryDetail
func (e *EpicRelatedUserStoryDetail) GetEpic(c *Client) (*Epic, error) {
	return c.Epic.Get(e.EpicID)
}

// EpicWatcherDetail => Epic watcher detail https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-watcher-detail
type EpicWatcherDetail struct {
	FullName string `json:"full_name,omitempty"`
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

// EpicsQueryParams holds fields to be used as URL query parameters to filter the queried objects
type EpicsQueryParams struct {
	Project        int    `url:"project,omitempty"`
	ProjectSlug    string `url:"project__slug,omitempty"`
	AssignedTo     int    `url:"assigned_to,omitempty"`
	StatusIsClosed bool   `url:"status__is_closed,omitempty"`
}

// EpicMinimal represent a small subset of a full Epic object
type EpicMinimal struct {
	Color   string  `json:"color"`
	ID      int     `json:"id"`
	Project Project `json:"project"`
	Ref     int     `json:"ref"`
	Subject string  `json:"subject"`
}
