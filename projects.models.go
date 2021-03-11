package taigo

import (
	"time"
)

func genericToProject(anyProjectObject interface{}) *Project {
	payloadProject := Project{}
	convertStructViaJSON(&anyProjectObject, &payloadProject)
	return &payloadProject
}

func genericToProjects(anyProjectObjectSlice interface{}) []Project {
	payloadProjectsSlice := []Project{}
	convertStructViaJSON(&anyProjectObjectSlice, &payloadProjectsSlice)
	return payloadProjectsSlice
}

// ProjectPoints represents the registered Agile Points to a project
type ProjectPoints struct {
	Order     int      `json:"order"`
	Name      string   `json:"name"`
	ID        int      `json:"id"`
	ProjectID int      `json:"project_id"`
	Value     *float64 `json:"value"`
}

// ProjectTagsColors is a [string/string] key/value pair to represent project-wide Tag/Colour combinations
// JSON Representation example:
// {
// 	"tags_colors": {
// 		"high": "#D35163",
// 		"normal": "#78D351"
//  },
// }
type ProjectTagsColors map[string]string

// IsValueNil returns true if ProjectPoints.Value is nil
func (pp ProjectPoints) IsValueNil() bool {
	if pp.Value == nil {
		return true
	}
	return false
}

// Project is a subset of all possible Project type variants
//
// https://taigaio.github.io/taiga-doc/dist/api.html#projects-create
type Project struct {
	ID                        int     `json:"id"`
	Slug                      string  `json:"slug"`
	CreationTemplate          int     `json:"creation_template"`
	Description               string  `json:"description"`
	IsBacklogActivated        bool    `json:"is_backlog_activated"`
	IsIssuesActivated         bool    `json:"is_issues_activated"`
	IsKanbanActivated         bool    `json:"is_kanban_activated"`
	IsPrivate                 bool    `json:"is_private"`
	IsWikiActivated           bool    `json:"is_wiki_activated"`
	Name                      string  `json:"name"`
	TotalMilestones           int     `json:"total_milestones"`
	TotalStoryPoints          float64 `json:"total_story_points"`
	Videoconferences          string  `json:"videoconferences"`
	VideoconferencesExtraData string  `json:"videoconferences_extra_data"`
	ProjectsLIST              *ProjectsList
	ProjectDETAIL             *ProjectDetail
}

// AsProject packs the returned ProjectDETAIL into a generic Project struct
func (p *ProjectDetail) AsProject() (*Project, error) {
	project := genericToProject(&p)
	project.ProjectDETAIL = p
	return project, nil
}

// AsProjects packs the returned ProjectsLIST into a generic Project []struct
func (p *ProjectsList) AsProjects() ([]Project, error) {
	projects := genericToProjects(&p)
	for i := 0; i < len(projects); i++ {
		projects[i].ProjectsLIST = p
	}
	return projects, nil
}

// ProjectDetail -> https://taigaio.github.io/taiga-doc/dist/api.html#object-project-detail
type ProjectDetail struct {
	AnonPermissions           []string                             `json:"anon_permissions"`
	BlockedCode               string                               `json:"blocked_code"`
	CreatedDate               time.Time                            `json:"created_date"`
	CreationTemplate          int                                  `json:"creation_template"`
	DefaultEpicStatus         int                                  `json:"default_epic_status"`
	DefaultIssueStatus        int                                  `json:"default_issue_status"`
	DefaultIssueType          int                                  `json:"default_issue_type"`
	DefaultPoints             float64                              `json:"default_points"`
	DefaultPriority           int                                  `json:"default_priority"`
	DefaultSeverity           int                                  `json:"default_severity"`
	DefaultTaskStatus         int                                  `json:"default_task_status"`
	DefaultUsStatus           int                                  `json:"default_us_status"`
	Description               string                               `json:"description"`
	EpicCustomAttributes      []EpicCustomAttributeDefinition      `json:"epic_custom_attributes"`
	EpicStatuses              []epicStatus                         `json:"epic_statuses"`
	EpicsCsvUUID              string                               `json:"epics_csv_uuid"`
	IAmAdmin                  bool                                 `json:"i_am_admin"`
	IAmMember                 bool                                 `json:"i_am_member"`
	IAmOwner                  bool                                 `json:"i_am_owner"`
	ID                        int                                  `json:"id"`
	IsBacklogActivated        bool                                 `json:"is_backlog_activated"`
	IsContactActivated        bool                                 `json:"is_contact_activated"`
	IsEpicsActivated          bool                                 `json:"is_epics_activated"`
	IsFan                     bool                                 `json:"is_fan"`
	IsFeatured                bool                                 `json:"is_featured"`
	IsIssuesActivated         bool                                 `json:"is_issues_activated"`
	IsKanbanActivated         bool                                 `json:"is_kanban_activated"`
	IsLookingForPeople        bool                                 `json:"is_looking_for_people"`
	IsOutOfOwnerLimits        bool                                 `json:"is_out_of_owner_limits"`
	IsPrivate                 bool                                 `json:"is_private"`
	IsPrivateExtraInfo        IsPrivateExtraInfo                   `json:"is_private_extra_info"`
	IsWatcher                 bool                                 `json:"is_watcher"`
	IsWikiActivated           bool                                 `json:"is_wiki_activated"`
	IssueCustomAttributes     []IssueCustomAttributeDefinition     `json:"issue_custom_attributes"`
	IssueDuedates             []issueDueDate                       `json:"issue_duedates"`
	IssueStatuses             []issueStatus                        `json:"issue_statuses"`
	IssueTypes                []issueType                          `json:"issue_types"`
	IssuesCsvUUID             string                               `json:"issues_csv_uuid"`
	LogoBigURL                string                               `json:"logo_big_url"`
	LogoSmallURL              string                               `json:"logo_small_url"`
	LookingForPeopleNote      string                               `json:"looking_for_people_note"`
	MaxMemberships            int                                  `json:"max_memberships"`
	Members                   []members                            `json:"members"`
	Milestones                []milestone                          `json:"milestones"`
	ModifiedDate              time.Time                            `json:"modified_date"`
	MyHomepage                interface{}                          `json:"my_homepage"`
	MyPermissions             []string                             `json:"my_permissions"`
	Name                      string                               `json:"name"`
	NotifyLevel               int                                  `json:"notify_level"`
	Owner                     Owner                                `json:"owner"`
	Points                    []ProjectPoints                      `json:"points"`
	Priorities                []priority                           `json:"priorities"`
	PublicPermissions         []string                             `json:"public_permissions"`
	Roles                     []roles                              `json:"roles"`
	Severities                []severity                           `json:"severities"`
	Slug                      string                               `json:"slug"`
	Tags                      []string                             `json:"tags"`
	TagsColors                ProjectTagsColors                    `json:"tags_colors"`
	TaskCustomAttributes      []TaskCustomAttributeDefinition      `json:"task_custom_attributes"`
	TaskDuedates              []taskDueDates                       `json:"task_duedates"`
	TaskStatuses              []taskStatus                         `json:"task_statuses"`
	TasksCsvUUID              string                               `json:"tasks_csv_uuid"`
	TotalActivity             int                                  `json:"total_activity"`
	TotalActivityLastMonth    int                                  `json:"total_activity_last_month"`
	TotalActivityLastWeek     int                                  `json:"total_activity_last_week"`
	TotalActivityLastYear     int                                  `json:"total_activity_last_year"`
	TotalClosedMilestones     int                                  `json:"total_closed_milestones"`
	TotalFans                 int                                  `json:"total_fans"`
	TotalFansLastMonth        int                                  `json:"total_fans_last_month"`
	TotalFansLastWeek         int                                  `json:"total_fans_last_week"`
	TotalFansLastYear         int                                  `json:"total_fans_last_year"`
	TotalMemberships          int                                  `json:"total_memberships"`
	TotalMilestones           int                                  `json:"total_milestones"`
	TotalStoryPoints          float64                              `json:"total_story_points"`
	TotalWatchers             int                                  `json:"total_watchers"`
	TotalsUpdatedDatetime     time.Time                            `json:"totals_updated_datetime"`
	TransferToken             string                               `json:"transfer_token"`
	UsDuedates                []userStoryDueDate                   `json:"us_duedates"`
	UsStatuses                []userStoryStatus                    `json:"us_statuses"`
	UserstoriesCsvUUID        string                               `json:"userstories_csv_uuid"`
	UserstoryCustomAttributes []UserStoryCustomAttributeDefinition `json:"userstory_custom_attributes"`
	Videoconferences          string                               `json:"videoconferences"`
	VideoconferencesExtraData string                               `json:"videoconferences_extra_data"`
}

// EpicCustomAttributeDefinition != EpicCustomAttribute
type EpicCustomAttributeDefinition struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	ProjectID    int         `json:"project_id"`
	Type         string      `json:"type"`
}

// epicStatus != EpicStatus
type epicStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

// IssueCustomAttributeDefinition != IssueCustomAttribute
type IssueCustomAttributeDefinition struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	ProjectID    int         `json:"project_id"`
	Type         string      `json:"type"`
}

type issueDueDate struct {
	ByDefault bool   `json:"by_default"`
	Color     string `json:"color"`
	DaysToDue int    `json:"days_to_due"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

// issueStatus != IssueStatus
type issueStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

type issueType struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

type members struct {
	Color           string `json:"color"`
	FullName        string `json:"full_name"`
	FullNameDisplay string `json:"full_name_display"`
	GravatarID      string `json:"gravatar_id"`
	ID              int    `json:"id"`
	IsActive        bool   `json:"is_active"`
	Photo           string `json:"photo"`
	Role            int    `json:"role"`
	RoleName        string `json:"role_name"`
	Username        string `json:"username"`
}
type milestone struct {
	Closed bool   `json:"closed"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

type priority struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

type roles struct {
	Computable  bool     `json:"computable"`
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Order       int      `json:"order"`
	Permissions []string `json:"permissions"`
	ProjectID   int      `json:"project_id"`
	Slug        string   `json:"slug"`
}

type severity struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

// TaskCustomAttributeDefinition != TaskCustomAttribute
type TaskCustomAttributeDefinition struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	ProjectID    int         `json:"project_id"`
	Type         string      `json:"type"`
}

type taskDueDates struct {
	ByDefault bool   `json:"by_default"`
	Color     string `json:"color"`
	DaysToDue int    `json:"days_to_due"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

// taskStatus != TaskStatus
type taskStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

type userStoryDueDate struct {
	ByDefault bool   `json:"by_default"`
	Color     string `json:"color"`
	DaysToDue int    `json:"days_to_due"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
}

// userStoryStatus != UserStoryStatus
type userStoryStatus struct {
	Color      string `json:"color"`
	ID         int    `json:"id"`
	IsArchived bool   `json:"is_archived"`
	IsClosed   bool   `json:"is_closed"`
	Name       string `json:"name"`
	Order      int    `json:"order"`
	ProjectID  int    `json:"project_id"`
	Slug       string `json:"slug"`
	WipLimit   int    `json:"wip_limit"`
}

// UserStoryCustomAttributeDefinition != UserStoryCustomAttribute
type UserStoryCustomAttributeDefinition struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	ProjectID    int         `json:"project_id"`
	Type         string      `json:"type"`
}

// ProjectsList -> https://taigaio.github.io/taiga-doc/dist/api.html#object-project-list-entry
type ProjectsList []struct {
	IsEpicsActivated          bool              `json:"is_epics_activated"`
	IsIssuesActivated         bool              `json:"is_issues_activated"`
	LogoSmallURL              string            `json:"logo_small_url"`
	LookingForPeopleNote      string            `json:"looking_for_people_note"`
	TotalActivityLastMonth    int               `json:"total_activity_last_month"`
	DefaultEpicStatus         int               `json:"default_epic_status"`
	DefaultSeverity           int               `json:"default_severity"`
	IsKanbanActivated         bool              `json:"is_kanban_activated"`
	Videoconferences          string            `json:"videoconferences"`
	ModifiedDate              time.Time         `json:"modified_date"`
	Name                      string            `json:"name"`
	IsLookingForPeople        bool              `json:"is_looking_for_people"`
	Description               string            `json:"description"`
	TotalClosedMilestones     int               `json:"total_closed_milestones"`
	DefaultUsStatus           int               `json:"default_us_status"`
	TotalFansLastMonth        int               `json:"total_fans_last_month"`
	TotalMilestones           int               `json:"total_milestones"`
	MyPermissions             []string          `json:"my_permissions"`
	Members                   []int             `json:"members"`
	Owner                     Owner             `json:"owner"`
	NotifyLevel               int               `json:"notify_level"`
	TagsColors                ProjectTagsColors `json:"tags_colors"`
	IsWikiActivated           bool              `json:"is_wiki_activated"`
	VideoconferencesExtraData string            `json:"videoconferences_extra_data"`
	CreatedDate               time.Time         `json:"created_date"`
	TotalWatchers             int               `json:"total_watchers"`
	IAmAdmin                  bool              `json:"i_am_admin"`
	DefaultIssueStatus        int               `json:"default_issue_status"`
	CreationTemplate          int               `json:"creation_template"`
	TotalStoryPoints          float64           `json:"total_story_points"`
	AnonPermissions           []string          `json:"anon_permissions"`
	TotalFans                 int               `json:"total_fans"`
	IsBacklogActivated        bool              `json:"is_backlog_activated"`
	ID                        int               `json:"id"`
	BlockedCode               string            `json:"blocked_code"`
	IsPrivate                 bool              `json:"is_private"`
	IsWatcher                 bool              `json:"is_watcher"`
	PublicPermissions         []string          `json:"public_permissions"`
	IsFan                     bool              `json:"is_fan"`
	TotalFansLastWeek         int               `json:"total_fans_last_week"`
	TotalActivityLastYear     int               `json:"total_activity_last_year"`
	DefaultPriority           int               `json:"default_priority"`
	IsContactActivated        bool              `json:"is_contact_activated"`
	Slug                      string            `json:"slug"`
	LogoBigURL                string            `json:"logo_big_url"`
	IsFeatured                bool              `json:"is_featured"`
	IAmOwner                  bool              `json:"i_am_owner"`
	TotalActivityLastWeek     int               `json:"total_activity_last_week"`
	Tags                      []string          `json:"tags"`
	DefaultIssueType          int               `json:"default_issue_type"`
	TotalsUpdatedDatetime     time.Time         `json:"totals_updated_datetime"`
	TotalActivity             int               `json:"total_activity"`
	IAmMember                 bool              `json:"i_am_member"`
	TotalFansLastYear         int               `json:"total_fans_last_year"`
	DefaultPoints             float64           `json:"default_points"`
	DefaultTaskStatus         int               `json:"default_task_status"`
}

// ProjectModulesConfiguration -> https://taigaio.github.io/taiga-doc/dist/api.html#object-project-modules-detail
type ProjectModulesConfiguration struct {
	Bitbucket struct {
		Secret         string   `json:"secret"`
		ValidOriginIps []string `json:"valid_origin_ips"`
		WebhooksURL    string   `json:"webhooks_url"`
	} `json:"bitbucket"`
	Github struct {
		Secret      string `json:"secret"`
		WebhooksURL string `json:"webhooks_url"`
	} `json:"github"`
	Gitlab struct {
		Secret         string   `json:"secret"`
		ValidOriginIps []string `json:"valid_origin_ips"`
		WebhooksURL    string   `json:"webhooks_url"`
	} `json:"gitlab"`
	Gogs struct {
		Secret      string `json:"secret"`
		WebhooksURL string `json:"webhooks_url"`
	} `json:"gogs"`
}

// ProjectStatsDetail -> https://taigaio.github.io/taiga-doc/dist/api.html#object-project-detail
type ProjectStatsDetail struct {
	AssignedPoints        float64     `json:"assigned_points"`
	AssignedPointsPerRole AgilePoints `json:"assigned_points_per_role"`
	ClosedPoints          int         `json:"closed_points"`
	ClosedPointsPerRole   struct {
	} `json:"closed_points_per_role"`
	DefinedPoints        float64     `json:"defined_points"`
	DefinedPointsPerRole AgilePoints `json:"defined_points_per_role"`
	Milestones           []struct {
		ClientIncrement int     `json:"client-increment"`
		Evolution       float64 `json:"evolution"`
		Name            string  `json:"name"`
		Optimal         float64 `json:"optimal"`
		TeamIncrement   int     `json:"team-increment"`
	} `json:"milestones"`
	Name            string  `json:"name"`
	Speed           int     `json:"speed"`
	TotalMilestones int     `json:"total_milestones"`
	TotalPoints     float64 `json:"total_points"`
}

/* EXTRAS */

// ProjectsQueryParameters holds fields to be used as URL query parameters to filter the queried objects
//
// To set `OrderBy`, use the methods attached to this struct
type ProjectsQueryParameters struct {
	Member             int    `url:"member,omitempty"`
	Members            []int  `url:"members,omitempty"`
	IsLookingForPeople bool   `url:"is_looking_for_people,omitempty"`
	IsFeatured         bool   `url:"is_featured,omitempty"`
	IsBacklogActivated bool   `url:"is_backlog_activated,omitempty"`
	IsKanbanActivated  bool   `url:"is_kanban_activated,omitempty"`
	orderBy            string `url:"order_by,omitempty"` // Can be set via struct methods
}

// MembershipsUserOrder => Order by the project order specified by the user
func (queryParams *ProjectsQueryParameters) MembershipsUserOrder() {
	queryParams.orderBy = "memberships__user_order"
}

// TotalFans => Order by total fans for the project
func (queryParams *ProjectsQueryParameters) TotalFans() {
	queryParams.orderBy = "total_fans"
}

// TotalFansLastWeek => Order by number of new fans in the last week
func (queryParams *ProjectsQueryParameters) TotalFansLastWeek() {
	queryParams.orderBy = "total_fans_last_week"
}

// TotalFansLastMonth => Order by number of new fans in the last month
func (queryParams *ProjectsQueryParameters) TotalFansLastMonth() {
	queryParams.orderBy = "total_fans_last_month"
}

// TotalFansLastYear => Order by number of new fans in the last year
func (queryParams *ProjectsQueryParameters) TotalFansLastYear() {
	queryParams.orderBy = "total_fans_last_year"
}

// TotalActivity => Order by number of history entries for the project
func (queryParams *ProjectsQueryParameters) TotalActivity() {
	queryParams.orderBy = "total_activity"
}

// TotalActivityLastWeek => Order by number of history entries generated in the last week
func (queryParams *ProjectsQueryParameters) TotalActivityLastWeek() {
	queryParams.orderBy = "total_activity_last_week"
}

// TotalActivityLastMonth => Order by number of history entries generated in the last month
func (queryParams *ProjectsQueryParameters) TotalActivityLastMonth() {
	queryParams.orderBy = "total_activity_last_month"
}

// TotalActivityLastYear => Order by number of history entries generated in the last year
func (queryParams *ProjectsQueryParameters) TotalActivityLastYear() {
	queryParams.orderBy = "total_activity_last_year"
}
