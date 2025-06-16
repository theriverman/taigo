package taigo

import "time"

func genericToUserStory(anyUsObject interface{}) *UserStory {
	object := UserStory{}
	convertStructViaJSON(&anyUsObject, &object)
	return &object
}

func genericToUserStories(anyUsObjectSlice interface{}) []UserStory {
	objects := []UserStory{}
	convertStructViaJSON(&anyUsObjectSlice, &objects)
	return objects
}

// UserStory represents a subset of (UserStoryDetail, UserStoryDetailGET, UserStoryDetailLIST) structs for creating new objects
type UserStory struct {
	TaigaBaseObject
	ID                  int         `json:"id,omitempty"`
	Ref                 int         `json:"ref,omitempty"`
	Version             int         `json:"version,omitempty"`
	AssignedTo          int         `json:"assigned_to,omitempty"`
	BacklogOrder        int64       `json:"backlog_order,omitempty"`
	BlockedNote         string      `json:"blocked_note,omitempty"`
	ClientRequirement   bool        `json:"client_requirement,omitempty"`
	Description         string      `json:"description,omitempty"`
	ExternalReference   []string    `json:"external_reference,omitempty"`
	IsBlocked           bool        `json:"is_blocked,omitempty"`
	IsClosed            bool        `json:"is_closed,omitempty"`
	KanbanOrder         int64       `json:"kanban_order,omitempty"`
	Milestone           int         `json:"milestone,omitempty"`
	Points              AgilePoints `json:"points,omitempty"`
	Project             int         `json:"project"`
	SprintOrder         int         `json:"sprint_order,omitempty"`
	Status              int         `json:"status,omitempty"`
	Subject             string      `json:"subject"`
	Tags                [][]string  `json:"tags,omitempty"`
	TeamRequirement     bool        `json:"team_requirement,omitempty"`
	Watchers            []int       `json:"watchers,omitempty"`
	UserStoryDetail     *UserStoryDetail
	UserStoryDetailGET  *UserStoryDetailGET
	UserStoryDetailLIST *UserStoryDetailLIST
}

// GetID returns the ID
func (us *UserStory) GetID() int {
	return us.ID
}

// GetRef returns the Ref
func (us *UserStory) GetRef() int {
	return us.Ref
}

// GetVersion return the version
func (us *UserStory) GetVersion() int {
	return us.Version
}

// GetSubject returns the subject
func (us *UserStory) GetSubject() string {
	return us.Subject
}

// GetProject returns the project ID
func (us *UserStory) GetProject() int {
	return us.Project
}

// UserStoryOrigin stores the bare minimum fields of the original User Story as a reference
// https://github.com/taigaio/taiga-back/blob/886ef47d0621388a11aeea61a269f00d13d32f8d/taiga/projects/userstories/serializers.py#L34
type UserStoryOrigin struct {
	ID      int    `json:"id,omitempty"`
	Ref     int    `json:"ref,omitempty"`
	Subject string `json:"subject,omitempty"`
}

// UserStoryDetailLIST => https://taigaio.github.io/taiga-doc/dist/api.html#object-userstory-detail-list
type UserStoryDetailLIST []struct {
	AssignedTo          int                       `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo       `json:"assigned_to_extra_info,omitempty"`
	AssignedUsers       []int                     `json:"assigned_users,omitempty"`
	Attachments         []GenericObjectAttachment `json:"attachments,omitempty"`
	BacklogOrder        int                       `json:"backlog_order,omitempty"`
	BlockedNote         string                    `json:"blocked_note,omitempty"`
	ClientRequirement   bool                      `json:"client_requirement,omitempty"`
	Comment             string                    `json:"comment,omitempty"`
	CreatedDate         time.Time                 `json:"created_date,omitempty"`
	DueDate             string                    `json:"due_date"`
	DueDateReason       string                    `json:"due_date_reason"`
	DueDateStatus       string                    `json:"due_date_status"`
	EpicOrder           int                       `json:"epic_order,omitempty"`
	Epics               []EpicMinimal             `json:"epics,omitempty"`
	ExternalReference   []string                  `json:"external_reference,omitempty"`
	FinishDate          time.Time                 `json:"finish_date,omitempty"`
	GeneratedFromIssue  int                       `json:"generated_from_issue,omitempty"`
	GeneratedFromTask   int                       `json:"generated_from_task,omitempty"`
	ID                  int                       `json:"id,omitempty"`
	IsBlocked           bool                      `json:"is_blocked,omitempty"`
	IsClosed            bool                      `json:"is_closed,omitempty"`
	IsVoter             bool                      `json:"is_voter,omitempty"`
	IsWatcher           bool                      `json:"is_watcher,omitempty"`
	KanbanOrder         int                       `json:"kanban_order,omitempty"`
	Milestone           int                       `json:"milestone,omitempty"`
	MilestoneName       string                    `json:"milestone_name,omitempty"`
	MilestoneSlug       string                    `json:"milestone_slug,omitempty"`
	ModifiedDate        time.Time                 `json:"modified_date,omitempty"`
	OriginIssue         *UserStoryOrigin          `json:"origin_issue,omitempty"`
	OriginTask          *UserStoryOrigin          `json:"origin_task,omitempty"`
	Owner               int                       `json:"owner,omitempty"`
	OwnerExtraInfo      OwnerExtraInfo            `json:"owner_extra_info,omitempty"`
	Points              AgilePoints               `json:"points,omitempty"`
	Project             int                       `json:"project,omitempty"`
	ProjectExtraInfo    ProjectExtraInfo          `json:"project_extra_info,omitempty"`
	Ref                 int                       `json:"ref,omitempty"`
	SprintOrder         int                       `json:"sprint_order,omitempty"`
	Status              int                       `json:"status,omitempty"`
	StatusExtraInfo     StatusExtraInfo           `json:"status_extra_info,omitempty"`
	Subject             string                    `json:"subject,omitempty"`
	Tags                Tags                      `json:"tags,omitempty"`
	Tasks               []UserStoryNestedTask     `json:"tasks"`
	TeamRequirement     bool                      `json:"team_requirement,omitempty"`
	TotalAttachments    int                       `json:"total_attachments,omitempty"`
	TotalComments       int                       `json:"total_comments,omitempty"`
	TotalPoints         float64                   `json:"total_points,omitempty"`
	TotalVoters         int                       `json:"total_voters,omitempty"`
	TotalWatchers       int                       `json:"total_watchers,omitempty"`
	TribeGig            TribeGig                  `json:"tribe_gig,omitempty"`
	Version             int                       `json:"version,omitempty"`
	Watchers            []int                     `json:"watchers,omitempty"`
}

// AsUserStory packs the returned UserStoryDetailLIST into a generic UserStory struct
func (u *UserStoryDetailLIST) AsUserStory() ([]UserStory, error) {
	userstories := genericToUserStories(&u)
	for i := 0; i < len(userstories); i++ {
		userstories[i].UserStoryDetailLIST = u
	}
	return userstories, nil
}

// UserStoryDetail => https://taigaio.github.io/taiga-doc/dist/api.html#object-userstory-detail
type UserStoryDetail struct {
	AssignedTo          int                       `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo       `json:"assigned_to_extra_info,omitempty"`
	AssignedUsers       []int                     `json:"assigned_users,omitempty"`
	Attachments         []GenericObjectAttachment `json:"attachments,omitempty"`
	BacklogOrder        int64                     `json:"backlog_order,omitempty"`
	BlockedNote         string                    `json:"blocked_note,omitempty"`
	BlockedNoteHTML     string                    `json:"blocked_note_html,omitempty"`
	ClientRequirement   bool                      `json:"client_requirement,omitempty"`
	Comment             string                    `json:"comment,omitempty"`
	CreatedDate         time.Time                 `json:"created_date,omitempty"`
	Description         string                    `json:"description,omitempty"`
	DescriptionHTML     string                    `json:"description_html,omitempty"`
	DueDate             string                    `json:"due_date,omitempty"`
	DueDateReason       string                    `json:"due_date_reason,omitempty"`
	DueDateStatus       string                    `json:"due_date_status,omitempty"`
	EpicOrder           int                       `json:"epic_order,omitempty"`
	Epics               []EpicMinimal             `json:"epics,omitempty"`
	ExternalReference   []string                  `json:"external_reference,omitempty"`
	FinishDate          string                    `json:"finish_date,omitempty"`
	GeneratedFromIssue  int                       `json:"generated_from_issue,omitempty"`
	GeneratedFromTask   int                       `json:"generated_from_task,omitempty"`
	ID                  int                       `json:"id,omitempty"`
	IsBlocked           bool                      `json:"is_blocked,omitempty"`
	IsClosed            bool                      `json:"is_closed,omitempty"`
	IsVoter             bool                      `json:"is_voter,omitempty"`
	IsWatcher           bool                      `json:"is_watcher,omitempty"`
	KanbanOrder         int64                     `json:"kanban_order,omitempty"`
	Milestone           int                       `json:"milestone,omitempty"`
	MilestoneName       string                    `json:"milestone_name,omitempty"`
	MilestoneSlug       string                    `json:"milestone_slug,omitempty"`
	ModifiedDate        time.Time                 `json:"modified_date,omitempty"`
	Neighbors           struct {
		Next struct {
			ID      int    `json:"id,omitempty"`
			Ref     int    `json:"ref,omitempty"`
			Subject string `json:"subject,omitempty"`
		} `json:"next,omitempty"`
		Previous struct {
			ID      int    `json:"id,omitempty"`
			Ref     int    `json:"ref,omitempty"`
			Subject string `json:"subject,omitempty"`
		} `json:"previous,omitempty"`
	} `json:"neighbors,omitempty"`
	OriginIssue      *UserStoryOrigin      `json:"origin_issue,omitempty"`
	OriginTask       *UserStoryOrigin      `json:"origin_task,omitempty"`
	Owner            int                   `json:"owner,omitempty"`
	OwnerExtraInfo   OwnerExtraInfo        `json:"owner_extra_info,omitempty"`
	Points           AgilePoints           `json:"points,omitempty"`
	Project          int                   `json:"project"`
	ProjectExtraInfo ProjectExtraInfo      `json:"project_extra_info,omitempty"`
	Ref              int                   `json:"ref,omitempty"`
	SprintOrder      int                   `json:"sprint_order,omitempty"`
	Status           int                   `json:"status,omitempty"`
	StatusExtraInfo  StatusExtraInfo       `json:"status_extra_info,omitempty"`
	Subject          string                `json:"subject"`
	Tags             Tags                  `json:"tags,omitempty"`
	Tasks            []UserStoryNestedTask `json:"tasks"`
	TeamRequirement  bool                  `json:"team_requirement,omitempty"`
	TotalAttachments int                   `json:"total_attachments,omitempty"`
	TotalComments    int                   `json:"total_comments,omitempty"`
	TotalPoints      float64               `json:"total_points,omitempty"`
	TotalVoters      int                   `json:"total_voters,omitempty"`
	TotalWatchers    int                   `json:"total_watchers,omitempty"`
	TribeGig         TribeGig              `json:"tribe_gig,omitempty"`
	Version          int                   `json:"version,omitempty"`
	Watchers         []int                 `json:"watchers,omitempty"`
}

// AsUserStory packs the returned UserStoryDetail into a generic UserStory struct
func (u *UserStoryDetail) AsUserStory() (*UserStory, error) {
	userstory := genericToUserStory(&u)
	userstory.UserStoryDetail = u
	return userstory, nil
}

// UserStoryDetailGET => https://taigaio.github.io/taiga-doc/dist/api.html#object-userstory-detail-get
type UserStoryDetailGET struct {
	AssignedTo          int                       `json:"assigned_to"`
	AssignedToExtraInfo AssignedToExtraInfo       `json:"assigned_to_extra_info"`
	AssignedUsers       []int                     `json:"assigned_users"`
	Attachments         []GenericObjectAttachment `json:"attachments"`
	BacklogOrder        int64                     `json:"backlog_order"`
	BlockedNote         string                    `json:"blocked_note"`
	BlockedNoteHTML     string                    `json:"blocked_note_html"`
	ClientRequirement   bool                      `json:"client_requirement"`
	Comment             string                    `json:"comment"`
	CreatedDate         time.Time                 `json:"created_date"`
	Description         string                    `json:"description"`
	DescriptionHTML     string                    `json:"description_html"`
	DueDate             string                    `json:"due_date"`
	DueDateReason       string                    `json:"due_date_reason"`
	DueDateStatus       string                    `json:"due_date_status"`
	EpicOrder           int                       `json:"epic_order"`
	Epics               []EpicMinimal             `json:"epics"`
	ExternalReference   []string                  `json:"external_reference"`
	FinishDate          time.Time                 `json:"finish_date"`
	GeneratedFromIssue  int                       `json:"generated_from_issue"`
	GeneratedFromTask   int                       `json:"generated_from_task"`
	ID                  int                       `json:"id"`
	IsBlocked           bool                      `json:"is_blocked"`
	IsClosed            bool                      `json:"is_closed"`
	IsVoter             bool                      `json:"is_voter"`
	IsWatcher           bool                      `json:"is_watcher"`
	KanbanOrder         int64                     `json:"kanban_order"`
	Milestone           int                       `json:"milestone"`
	MilestoneName       string                    `json:"milestone_name"`
	MilestoneSlug       string                    `json:"milestone_slug"`
	ModifiedDate        time.Time                 `json:"modified_date"`
	Neighbors           Neighbors                 `json:"neighbors"`
	OriginIssue         int                       `json:"origin_issue"`
	OriginTask          *UserStoryOrigin          `json:"origin_task"`
	Owner               int                       `json:"owner"`
	OwnerExtraInfo      OwnerExtraInfo            `json:"owner_extra_info"`
	Points              Points                    `json:"points"`
	Project             int                       `json:"project"`
	ProjectExtraInfo    ProjectExtraInfo          `json:"project_extra_info"`
	Ref                 int                       `json:"ref"`
	SprintOrder         int                       `json:"sprint_order"`
	Status              int                       `json:"status"`
	StatusExtraInfo     StatusExtraInfo           `json:"status_extra_info"`
	Subject             string                    `json:"subject"`
	Tags                Tags                      `json:"tags"`
	Tasks               []UserStoryNestedTask     `json:"tasks"`
	TeamRequirement     bool                      `json:"team_requirement"`
	TotalAttachments    int                       `json:"total_attachments"`
	TotalComments       int                       `json:"total_comments"`
	TotalPoints         float64                   `json:"total_points"`
	TotalVoters         int                       `json:"total_voters"`
	TotalWatchers       int                       `json:"total_watchers"`
	TribeGig            TribeGig                  `json:"tribe_gig"`
	Version             int                       `json:"version"`
	Watchers            []int                     `json:"watchers"`
}

// AsUserStory packs the returned UserStoryDetailGET into a generic UserStory struct
func (u *UserStoryDetailGET) AsUserStory() (*UserStory, error) {
	userstory := genericToUserStory(&u)
	userstory.UserStoryDetailGET = u
	return userstory, nil
}

// IssueFiltersDataDetail => https://taigaio.github.io/taiga-doc/dist/api.html#object-userstory-filters-data
type IssueFiltersDataDetail struct {
	AssignedTo []struct {
		Count    int    `json:"count,omitempty"`
		FullName string `json:"full_name,omitempty"`
		ID       int    `json:"id,omitempty"`
	} `json:"assigned_to,omitempty"`
	AssignedUsers []struct {
		Count    int    `json:"count,omitempty"`
		FullName string `json:"full_name,omitempty"`
		ID       int    `json:"id,omitempty"`
	} `json:"assigned_users,omitempty"`
	Epics  []EpicMinimal `json:"epics,omitempty"`
	Owners []struct {
		Count    int    `json:"count,omitempty"`
		FullName string `json:"full_name,omitempty"`
		ID       int    `json:"id,omitempty"`
	} `json:"owners,omitempty"`
	Roles []struct {
		Color string `json:"color,omitempty"`
		Count int    `json:"count,omitempty"`
		ID    int    `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Order int    `json:"order,omitempty"`
	} `json:"roles,omitempty"`
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

/* EXTRA FIELDS */

// UserStoryQueryParams holds fields to be used as URL query parameters to filter the queried objects
//
// To set `OrderBy`, use the methods attached to this struct
type UserStoryQueryParams struct {
	Project            int    `url:"project,omitempty"`
	Milestone          int    `url:"milestone,omitempty"`
	MilestoneIsNull    bool   `url:"milestone__isnull,omitempty"`
	Status             int    `url:"status,omitempty"`
	StatusIsArchived   bool   `url:"status__is_archived,omitempty"`
	Tags               string `url:"tags,omitempty"` // comma-separated strings w/o whitespace
	Watchers           int    `url:"watchers,omitempty"`
	AssignedTo         int    `url:"assigned_to,omitempty"`
	Epic               int    `url:"epic,omitempty"`
	Role               int    `url:"role,omitempty"`
	StatusIsClosed     bool   `url:"status__is_closed,omitempty"`
	IncludeAttachments bool   `url:"include_attachments,omitempty"`
	IncludeTasks       bool   `url:"include_tasks,omitempty"`
	ExcludeStatus      int    `url:"exclude_status,omitempty"`
	ExcludeTags        string `url:"exclude_tags,omitempty"` // comma-separated strings w/o whitespace
	ExcludeAssignedTo  int    `url:"exclude_assigned_to,omitempty"`
	ExcludeRole        int    `url:"exclude_role,omitempty"`
	ExcludeEpic        int    `url:"exclude_epic,omitempty"`
}

// UserStoryNestedTask is returned only when IncludeTasks is set to true in UserStoryQueryParams
type UserStoryNestedTask struct {
	Subject   string `json:"subject"`
	ID        int    `json:"id"`
	Ref       int    `json:"ref"`
	IsBlocked bool   `json:"is_blocked"`
	IsIocaine bool   `json:"is_iocaine"`
	StatusID  int    `json:"status_id"`
	IsClosed  bool   `json:"is_closed"`
}
