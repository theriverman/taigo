package taigo

import (
	"log"
	"time"
)

func genericToIssue(anyIssueObject interface{}) *Issue {
	payloadIssue := Issue{}
	convertStructViaJSON(&anyIssueObject, &payloadIssue)
	return &payloadIssue
}

func genericToIssues(anyIssueObjectSlice interface{}) []Issue {
	payloadIssuesSlice := []Issue{}
	err := convertStructViaJSON(&anyIssueObjectSlice, &payloadIssuesSlice)
	if err != nil {
		log.Println(err)
	}
	return payloadIssuesSlice
}

// Issue represents the mandatory fields of an Issue only
type Issue struct {
	TaigaBaseObject
	ID              int       `json:"id,omitempty"`
	Ref             int       `json:"ref,omitempty"`
	Version         int       `json:"version,omitempty"`
	AssignedTo      int       `json:"assigned_to,omitempty"`
	BlockedNote     string    `json:"blocked_note,omitempty"`
	Description     string    `json:"description,omitempty"`
	IsBlocked       bool      `json:"is_blocked,omitempty"`
	IsClosed        bool      `json:"is_closed,omitempty"`
	Milestone       int       `json:"milestone,omitempty"`
	Owner           int       `json:"owner,omitempty"`
	Priority        int       `json:"priority,omitempty"`
	Project         int       `json:"project"`
	Severity        int       `json:"severity,omitempty"`
	Status          int       `json:"status,omitempty"`
	Subject         string    `json:"subject"`
	Tags            Tags      `json:"tags,omitempty"`
	Type            int       `json:"type,omitempty"`
	Watchers        []int     `json:"watchers,omitempty"`
	CreatedDate     time.Time `json:"created_date,omitempty"`
	ModifiedDate    time.Time `json:"modified_date,omitempty"`
	FinishedDate    time.Time `json:"finished_date,omitempty"`
	DueDate         string    `json:"due_date,omitempty"`
	DueDateReason   string    `json:"due_date_reason,omitempty"`
	DueDateStatus   string    `json:"due_date_status,omitempty"`
	IssueDetail     *IssueDetail
	IssueDetailGET  *IssueDetailGET
	IssueDetailLIST *IssueDetailLIST
}

// GetID returns the ID
func (i *Issue) GetID() int {
	return i.ID
}

// GetRef returns the Ref
func (i *Issue) GetRef() int {
	return i.Ref
}

// GetVersion return the version
func (i *Issue) GetVersion() int {
	return i.Version
}

// GetSubject returns the subject
func (i *Issue) GetSubject() string {
	return i.Subject
}

// GetProject returns the project ID
func (i *Issue) GetProject() int {
	return i.Project
}

// IssueDetailLIST -> Issue detail (LIST)
//
// https://taigaio.github.io/taiga-doc/dist/api.html#object-issue-detail-list
type IssueDetailLIST []struct {
	AssignedTo          int                       `json:"assigned_to"`
	AssignedToExtraInfo AssignedToExtraInfo       `json:"assigned_to_extra_info"`
	Attachments         []GenericObjectAttachment `json:"attachments"`
	BlockedNote         string                    `json:"blocked_note"`
	CreatedDate         time.Time                 `json:"created_date"`
	ModifiedDate        time.Time                 `json:"modified_date"`
	FinishedDate        time.Time                 `json:"finished_date"`
	DueDate             string                    `json:"due_date"`
	DueDateReason       string                    `json:"due_date_reason"`
	DueDateStatus       string                    `json:"due_date_status"`
	ExternalReference   []int                     `json:"external_reference"`
	ID                  int                       `json:"id"`
	IsBlocked           bool                      `json:"is_blocked"`
	IsClosed            bool                      `json:"is_closed"`
	IsVoter             bool                      `json:"is_voter"`
	IsWatcher           bool                      `json:"is_watcher"`
	Milestone           int                       `json:"milestone"`
	Owner               int                       `json:"owner"`
	OwnerExtraInfo      OwnerExtraInfo            `json:"owner_extra_info"`
	Priority            int                       `json:"priority"`
	Project             int                       `json:"project"`
	ProjectExtraInfo    ProjectExtraInfo          `json:"project_extra_info"`
	Ref                 int                       `json:"ref"`
	Severity            int                       `json:"severity"`
	Status              int                       `json:"status"`
	StatusExtraInfo     StatusExtraInfo           `json:"status_extra_info"`
	Subject             string                    `json:"subject"`
	Tags                Tags                      `json:"tags"`
	TotalVoters         int                       `json:"total_voters"`
	TotalWatchers       int                       `json:"total_watchers"`
	Type                int                       `json:"type"`
	Version             int                       `json:"version"`
	Watchers            []int                     `json:"watchers"`
}

// AsIssues packs the returned IssueDetailLIST into a generic Issue struct
func (issueL *IssueDetailLIST) AsIssues() ([]Issue, error) {
	issues := genericToIssues(&issueL)
	for i := 0; i < len(issues); i++ {
		issues[i].IssueDetailLIST = issueL
	}
	return issues, nil
}

// IssueDetailGET -> Issue detail (GET)
// https://taigaio.github.io/taiga-doc/dist/api.html#object-issue-detail-get
type IssueDetailGET struct {
	AssignedTo           int                       `json:"assigned_to"`
	AssignedToExtraInfo  AssignedToExtraInfo       `json:"assigned_to_extra_info"`
	Attachments          []GenericObjectAttachment `json:"attachments"`
	BlockedNote          string                    `json:"blocked_note"`
	BlockedNoteHTML      string                    `json:"blocked_note_html"`
	Comment              string                    `json:"comment"`
	CreatedDate          time.Time                 `json:"created_date"`
	Description          string                    `json:"description"`
	DescriptionHTML      string                    `json:"description_html"`
	DueDate              string                    `json:"due_date"`
	DueDateReason        string                    `json:"due_date_reason"`
	DueDateStatus        string                    `json:"due_date_status"`
	ExternalReference    []int                     `json:"external_reference"`
	FinishedDate         time.Time                 `json:"finished_date"`
	GeneratedUserStories []int                     `json:"generated_user_stories"`
	ID                   int                       `json:"id"`
	IsBlocked            bool                      `json:"is_blocked"`
	IsClosed             bool                      `json:"is_closed"`
	IsVoter              bool                      `json:"is_voter"`
	IsWatcher            bool                      `json:"is_watcher"`
	Milestone            int                       `json:"milestone"`
	ModifiedDate         time.Time                 `json:"modified_date"`
	Neighbors            Neighbors                 `json:"neighbors"`
	Owner                int                       `json:"owner"`
	OwnerExtraInfo       OwnerExtraInfo            `json:"owner_extra_info"`
	Priority             int                       `json:"priority"`
	Project              int                       `json:"project"`
	ProjectExtraInfo     ProjectExtraInfo          `json:"project_extra_info"`
	Ref                  int                       `json:"ref"`
	Severity             int                       `json:"severity"`
	Status               int                       `json:"status"`
	StatusExtraInfo      StatusExtraInfo           `json:"status_extra_info"`
	Subject              string                    `json:"subject"`
	Tags                 Tags                      `json:"tags"`
	TotalVoters          int                       `json:"total_voters"`
	TotalWatchers        int                       `json:"total_watchers"`
	Type                 int                       `json:"type"`
	Version              int                       `json:"version"`
	Watchers             []int                     `json:"watchers"`
}

// AsIssue packs the returned IssueDetailGET into a generic Issue struct
func (issueD *IssueDetailGET) AsIssue() (*Issue, error) {
	issue := genericToIssue(&issueD)
	issue.IssueDetailGET = issueD
	return issue, nil
}

// IssueDetail -> Issue detail
// https://taigaio.github.io/taiga-doc/dist/api.html#object-issue-detail
type IssueDetail struct {
	AssignedTo           int                       `json:"assigned_to"`
	AssignedToExtraInfo  AssignedToExtraInfo       `json:"assigned_to_extra_info"`
	Attachments          []GenericObjectAttachment `json:"attachments"`
	BlockedNote          string                    `json:"blocked_note"`
	BlockedNoteHTML      string                    `json:"blocked_note_html"`
	Comment              string                    `json:"comment"`
	CreatedDate          time.Time                 `json:"created_date"`
	Description          string                    `json:"description"`
	DescriptionHTML      string                    `json:"description_html"`
	DueDate              string                    `json:"due_date"`
	DueDateReason        string                    `json:"due_date_reason"`
	DueDateStatus        string                    `json:"due_date_status"`
	ExternalReference    []int                     `json:"external_reference"`
	FinishedDate         time.Time                 `json:"finished_date"`
	GeneratedUserStories []int                     `json:"generated_user_stories"`
	ID                   int                       `json:"id"`
	IsBlocked            bool                      `json:"is_blocked"`
	IsClosed             bool                      `json:"is_closed"`
	IsVoter              bool                      `json:"is_voter"`
	IsWatcher            bool                      `json:"is_watcher"`
	Milestone            int                       `json:"milestone"`
	ModifiedDate         time.Time                 `json:"modified_date"`
	Neighbors            Neighbors                 `json:"neighbors"`
	Owner                int                       `json:"owner"`
	OwnerExtraInfo       OwnerExtraInfo            `json:"owner_extra_info"`
	Priority             int                       `json:"priority"`
	Project              int                       `json:"project"`
	ProjectExtraInfo     ProjectExtraInfo          `json:"project_extra_info"`
	Ref                  int                       `json:"ref"`
	Severity             int                       `json:"severity"`
	Status               int                       `json:"status"`
	StatusExtraInfo      StatusExtraInfo           `json:"status_extra_info"`
	Subject              string                    `json:"subject"`
	Tags                 Tags                      `json:"tags"`
	TotalVoters          int                       `json:"total_voters"`
	TotalWatchers        int                       `json:"total_watchers"`
	Type                 int                       `json:"type"`
	Version              int                       `json:"version"`
	Watchers             []int                     `json:"watchers"`
}

// AsIssue packs the returned IssueDetailGET into a generic Issue struct
func (issueD *IssueDetail) AsIssue() (*Issue, error) {
	issue := genericToIssue(&issueD)
	issue.IssueDetail = issueD
	return issue, nil
}

// IssueQueryParams holds fields to be used as URL query parameters to filter the queried objects
type IssueQueryParams struct {
	Project            int    `url:"project,omitempty"`
	Milestone          int    `url:"milestone,omitempty"`
	MilestoneIsNull    bool   `url:"milestone__isnull,omitempty"`
	Status             int    `url:"status,omitempty"`
	StatusIsArchived   bool   `url:"status__is_archived,omitempty"`
	Tags               string `url:"tags,omitempty"` // Strings separated by comma `,`
	Watchers           int    `url:"watchers,omitempty"`
	AssignedTo         int    `url:"assigned_to,omitempty"`
	Epic               int    `url:"epic,omitempty"`
	Role               int    `url:"role,omitempty"`
	StatusIsClosed     bool   `url:"status__is_closed,omitempty"`
	Type               int    `url:"type,omitempty"`
	Severity           int    `url:"severity,omitempty"`
	Priority           int    `url:"priority,omitempty"`
	Owner              int    `url:"owner,omitempty"`
	ExcludeStatus      int    `url:"exclude_status,omitempty"`
	ExcludeTags        int    `url:"exclude_tags,omitempty"` // Strings separated by comma `,`
	ExcludeAssignedTo  int    `url:"exclude_assigned_to,omitempty"`
	ExcludeRole        int    `url:"exclude_role,omitempty"`
	ExcludeEpic        int    `url:"exclude_epic,omitempty"`
	ExcludeSeverity    int    `url:"exclude_severity,omitempty"`
	ExcludePriority    int    `url:"exclude_priority,omitempty"`
	ExcludeOwner       int    `url:"exclude_owner,omitempty"`
	ExcludeType        int    `url:"exclude_type,omitempty"`
	IncludeAttachments bool   `url:"include_attachments,omitempty"`
}
