package taigo

import "time"

func genericToTask(anyTaskObject interface{}) *Task {
	object := Task{}
	convertStructViaJSON(&anyTaskObject, &object)
	return &object
}

func genericToTasks(anyTaskObjectSlice interface{}) []Task {
	objects := []Task{}
	convertStructViaJSON(&anyTaskObjectSlice, &objects)
	return objects
}

// Task represents a subset of (TaskDetail, TaskDetailGET, TaskDetailLIST)
type Task struct {
	ID                int
	AssignedTo        int         `json:"assigned_to,omitempty"`
	BlockedNote       string      `json:"blocked_note,omitempty"`
	Description       string      `json:"description,omitempty"`
	ExternalReference interface{} `json:"external_reference,omitempty"`
	IsBlocked         bool        `json:"is_blocked,omitempty"`
	IsClosed          bool        `json:"is_closed,omitempty"`
	IsIocaine         bool        `json:"is_iocaine,omitempty"`
	Milestone         int         `json:"milestone,omitempty"`
	Project           int         `json:"project,omitempty"`
	Ref               int
	Status            int      `json:"status,omitempty"`
	Subject           string   `json:"subject,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	TaskboardOrder    int      `json:"taskboard_order,omitempty"`
	UsOrder           int      `json:"us_order,omitempty"`
	UserStory         int      `json:"user_story,omitempty"`
	Version           int
	Watchers          []int `json:"watchers,omitempty"`
	TaskDetail        *TaskDetail
	TaskDetailGET     *TaskDetailGET
	TaskDetailLIST    *TaskDetailLIST
}

// TaskDetailLIST => https://taigaio.github.io/taiga-doc/dist/api.html#object-task-detail-list
type TaskDetailLIST []struct {
	AssignedTo          int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments         []Attachment        `json:"attachments,omitempty"`
	BlockedNote         string              `json:"blocked_note,omitempty"`
	CreatedDate         time.Time           `json:"created_date,omitempty"`
	DueDate             string              `json:"due_date,omitempty"`
	DueDateReason       string              `json:"due_date_reason,omitempty"`
	DueDateStatus       string              `json:"due_date_status,omitempty"`
	ExternalReference   interface{}         `json:"external_reference,omitempty"`
	FinishedDate        time.Time           `json:"finished_date,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IsBlocked           bool                `json:"is_blocked,omitempty"`
	IsClosed            bool                `json:"is_closed,omitempty"`
	IsIocaine           bool                `json:"is_iocaine,omitempty"`
	IsVoter             bool                `json:"is_voter,omitempty"`
	IsWatcher           bool                `json:"is_watcher,omitempty"`
	Milestone           int                 `json:"milestone,omitempty"`
	MilestoneSlug       string              `json:"milestone_slug,omitempty"`
	ModifiedDate        time.Time           `json:"modified_date,omitempty"`
	Owner               int                 `json:"owner,omitempty"`
	OwnerExtraInfo      OwnerExtraInfo      `json:"owner_extra_info,omitempty"`
	Project             int                 `json:"project,omitempty"`
	ProjectExtraInfo    ProjectExtraInfo    `json:"project_extra_info,omitempty"`
	Ref                 int                 `json:"ref,omitempty"`
	Status              int                 `json:"status,omitempty"`
	StatusExtraInfo     StatusExtraInfo     `json:"status_extra_info,omitempty"`
	Subject             string              `json:"subject,omitempty"`
	Tags                Tags                `json:"tags,omitempty"`
	TaskboardOrder      int64               `json:"taskboard_order,omitempty"`
	TotalComments       int                 `json:"total_comments,omitempty"`
	TotalVoters         int                 `json:"total_voters,omitempty"`
	TotalWatchers       int                 `json:"total_watchers,omitempty"`
	UsOrder             int64               `json:"us_order,omitempty"`
	UserStory           int                 `json:"user_story,omitempty"`
	UserStoryExtraInfo  UserStoryExtraInfo  `json:"user_story_extra_info,omitempty"`
	Version             int                 `json:"version,omitempty"`
	Watchers            []int               `json:"watchers,omitempty"`
}

// AsTask packs the returned TaskDetailLIST into a generic Task struct
func (t *TaskDetailLIST) AsTask() ([]Task, error) {
	tasks := genericToTasks(&t)
	for i := 0; i < len(tasks); i++ {
		tasks[i].TaskDetailLIST = t
	}
	return tasks, nil
}

// TaskDetailGET => https://taigaio.github.io/taiga-doc/dist/api.html#object-task-detail-get
type TaskDetailGET struct {
	AssignedTo           int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo  AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments          []interface{}       `json:"attachments,omitempty"`
	BlockedNote          string              `json:"blocked_note,omitempty"`
	BlockedNoteHTML      string              `json:"blocked_note_html,omitempty"`
	Comment              string              `json:"comment,omitempty"`
	CreatedDate          time.Time           `json:"created_date,omitempty"`
	Description          string              `json:"description,omitempty"`
	DescriptionHTML      string              `json:"description_html,omitempty"`
	DueDate              string              `json:"due_date,omitempty"`
	DueDateReason        string              `json:"due_date_reason,omitempty"`
	DueDateStatus        string              `json:"due_date_status,omitempty"`
	ExternalReference    interface{}         `json:"external_reference,omitempty"`
	FinishedDate         time.Time           `json:"finished_date,omitempty"`
	GeneratedUserStories interface{}         `json:"generated_user_stories,omitempty"`
	ID                   int                 `json:"id,omitempty"`
	IsBlocked            bool                `json:"is_blocked,omitempty"`
	IsClosed             bool                `json:"is_closed,omitempty"`
	IsIocaine            bool                `json:"is_iocaine,omitempty"`
	IsVoter              bool                `json:"is_voter,omitempty"`
	IsWatcher            bool                `json:"is_watcher,omitempty"`
	Milestone            int                 `json:"milestone,omitempty"`
	MilestoneSlug        string              `json:"milestone_slug,omitempty"`
	ModifiedDate         time.Time           `json:"modified_date,omitempty"`
	Neighbors            Neighbors           `json:"neighbors,omitempty"`
	Owner                int                 `json:"owner,omitempty"`
	OwnerExtraInfo       OwnerExtraInfo      `json:"owner_extra_info,omitempty"`
	Project              int                 `json:"project,omitempty"`
	ProjectExtraInfo     ProjectExtraInfo    `json:"project_extra_info,omitempty"`
	Ref                  int                 `json:"ref,omitempty"`
	Status               int                 `json:"status,omitempty"`
	StatusExtraInfo      StatusExtraInfo     `json:"status_extra_info,omitempty"`
	Subject              string              `json:"subject,omitempty"`
	Tags                 Tags                `json:"tags,omitempty"`
	TaskboardOrder       int64               `json:"taskboard_order,omitempty"`
	TotalComments        int                 `json:"total_comments,omitempty"`
	TotalVoters          int                 `json:"total_voters,omitempty"`
	TotalWatchers        int                 `json:"total_watchers,omitempty"`
	UsOrder              int64               `json:"us_order,omitempty"`
	UserStory            int                 `json:"user_story,omitempty"`
	UserStoryExtraInfo   UserStoryExtraInfo  `json:"user_story_extra_info,omitempty"`
	Version              int                 `json:"version,omitempty"`
	Watchers             []int               `json:"watchers,omitempty"`
}

// AsTask packs the returned TaskDetailGET into a generic Task struct
func (t *TaskDetailGET) AsTask() (*Task, error) {
	task := genericToTask(&t)
	task.TaskDetailGET = t
	return task, nil

}

// TaskDetail => https://taigaio.github.io/taiga-doc/dist/api.html#object-task-detail
type TaskDetail struct {
	AssignedTo           int                 `json:"assigned_to,omitempty"`
	AssignedToExtraInfo  AssignedToExtraInfo `json:"assigned_to_extra_info,omitempty"`
	Attachments          []interface{}       `json:"attachments,omitempty"`
	BlockedNote          string              `json:"blocked_note,omitempty"`
	BlockedNoteHTML      string              `json:"blocked_note_html,omitempty"`
	Comment              string              `json:"comment,omitempty"`
	CreatedDate          time.Time           `json:"created_date,omitempty"`
	Description          string              `json:"description,omitempty"`
	DescriptionHTML      string              `json:"description_html,omitempty"`
	DueDate              string              `json:"due_date,omitempty"`
	DueDateReason        string              `json:"due_date_reason,omitempty"`
	DueDateStatus        string              `json:"due_date_status,omitempty"`
	ExternalReference    interface{}         `json:"external_reference,omitempty"`
	FinishedDate         time.Time           `json:"finished_date,omitempty"`
	GeneratedUserStories interface{}         `json:"generated_user_stories,omitempty"`
	ID                   int                 `json:"id,omitempty"`
	IsBlocked            bool                `json:"is_blocked,omitempty"`
	IsClosed             bool                `json:"is_closed,omitempty"`
	IsIocaine            bool                `json:"is_iocaine,omitempty"`
	IsVoter              bool                `json:"is_voter,omitempty"`
	IsWatcher            bool                `json:"is_watcher,omitempty"`
	Milestone            int                 `json:"milestone,omitempty"`
	MilestoneSlug        string              `json:"milestone_slug,omitempty"`
	ModifiedDate         time.Time           `json:"modified_date,omitempty"`
	Neighbors            struct {
		Next struct {
			ID      int    `json:"id,omitempty"`
			Ref     int    `json:"ref,omitempty"`
			Subject string `json:"subject,omitempty"`
		} `json:"next,omitempty"`
		Previous interface{} `json:"previous,omitempty"`
	} `json:"neighbors,omitempty"`
	Owner              int                `json:"owner,omitempty"`
	OwnerExtraInfo     OwnerExtraInfo     `json:"owner_extra_info,omitempty"`
	Project            int                `json:"project,omitempty"`
	ProjectExtraInfo   ProjectExtraInfo   `json:"project_extra_info,omitempty"`
	Ref                int                `json:"ref,omitempty"`
	Status             int                `json:"status,omitempty"`
	StatusExtraInfo    StatusExtraInfo    `json:"status_extra_info,omitempty"`
	Subject            string             `json:"subject,omitempty"`
	Tags               Tags               `json:"tags,omitempty"`
	TaskboardOrder     int64              `json:"taskboard_order,omitempty"`
	TotalComments      int                `json:"total_comments,omitempty"`
	TotalVoters        int                `json:"total_voters,omitempty"`
	TotalWatchers      int                `json:"total_watchers,omitempty"`
	UsOrder            int64              `json:"us_order,omitempty"`
	UserStory          int                `json:"user_story,omitempty"`
	UserStoryExtraInfo UserStoryExtraInfo `json:"user_story_extra_info,omitempty"`
	Version            int                `json:"version,omitempty"`
	Watchers           []int              `json:"watchers,omitempty"`
}

// AsTask packs the returned TaskDetail into a generic Task struct
func (t *TaskDetail) AsTask() (*Task, error) {
	task := genericToTask(&t)
	task.TaskDetail = t
	return task, nil
}

// TaskVoterDetail => https://taigaio.github.io/taiga-doc/dist/api.html#object-task-voter-detail
type TaskVoterDetail struct {
	FullName string `json:"full_name,omitempty"`
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

/* EXTRAS */

// TasksQueryParams holds fields to be used as URL query parameters to filter the queried objects
//
// To set `OrderBy`, use the methods attached to this struct
type TasksQueryParams struct {
	Project           int      `url:"project,omitempty"`
	Status            int      `url:"status,omitempty"`
	Tags              []string `url:"tags,omitempty"`
	UserStory         int      `url:"user_story,omitempty"`
	Role              int      `url:"role,omitempty"`
	Owner             int      `url:"owner,omitempty"`
	Milestone         int      `url:"milestone,omitempty"`
	Watchers          int      `url:"watchers,omitempty"`
	AssignedTo        int      `url:"assigned_to,omitempty"`
	StatusIsClosed    bool     `url:"status__is_closed,omitempty"`
	ExcludeStatus     int      `url:"exclude_status,omitempty"`
	ExcludeTags       int      `url:"exclude_tags,omitempty"`
	ExcludeRole       int      `url:"exclude_role,omitempty"`
	ExcludeOwner      int      `url:"exclude_owner,omitempty"`
	ExcludeAssignedTo int      `url:"exclude_assigned_to,omitempty"`
}

// TaskMinimal represent a small subset of a full Task object
type TaskMinimal struct {
	Color   string  `json:"color"`
	ID      int     `json:"id"`
	Project Project `json:"project"`
	Ref     int     `json:"ref"`
	Subject string  `json:"subject"`
}
