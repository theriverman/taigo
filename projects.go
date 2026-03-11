package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// ProjectService is a handle to actions related to Projects
//
// https://taigaio.github.io/taiga-doc/dist/api.html#projects
type ProjectService struct {
	client *Client
	// defaultProjectID int
	Endpoint string
	// Mapped services for simple access
	areMappedServicesConfigured bool
	Auth                        *AuthService
	Epic                        *EpicService
	Issue                       *IssueService
	Milestone                   *MilestoneService
	Resolver                    *ResolverService
	Stats                       *StatsService
	Task                        *TaskService
	UserStory                   *UserStoryService
	User                        *UserService
	Webhook                     *WebhookService
	Wiki                        *WikiService
	Point                       *PointService
	Priority                    *PriorityService
	Severity                    *SeverityService
	IssueType                   *IssueTypeService
	EpicStatus                  *EpicStatusService
	IssueStatus                 *IssueStatusService
	TaskStatus                  *TaskStatusService
	UserStoryStatus             *UserStoryStatusService
	EpicCustomAttribute         *EpicCustomAttributeService
	IssueCustomAttribute        *IssueCustomAttributeService
	TaskCustomAttribute         *TaskCustomAttributeService
	UserStoryCustomAttribute    *UserStoryCustomAttributeService
	Application                 *ApplicationService
	ApplicationToken            *ApplicationTokenService
	Search                      *SearchService
	UserStorage                 *UserStorageService
	ProjectTemplate             *ProjectTemplateService
	ProjectTemplateDetail       *ProjectTemplateDetailService
	MembershipInvitation        *MembershipInvitationService
	WikiLink                    *WikiLinkService
	History                     *HistoryService
	NotifyPolicy                *NotifyPolicyService
	Contact                     *ContactService
	Feedback                    *FeedbackService
	ExportImport                *ExportImportService
	Timeline                    *TimelineService
	Locale                      *LocaleService
	Importer                    *ImporterService
	ContribPlugin               *ContribPluginService
	ObjectsSummary              *ObjectsSummaryService
}

type projectCreatePayload struct {
	CreationTemplate          int    `json:"creation_template,omitempty"`
	Description               string `json:"description"`
	IsBacklogActivated        bool   `json:"is_backlog_activated,omitempty"`
	IsIssuesActivated         bool   `json:"is_issues_activated,omitempty"`
	IsKanbanActivated         bool   `json:"is_kanban_activated,omitempty"`
	IsPrivate                 bool   `json:"is_private,omitempty"`
	IsWikiActivated           bool   `json:"is_wiki_activated,omitempty"`
	Name                      string `json:"name"`
	Videoconferences          string `json:"videoconferences,omitempty"`
	VideoconferencesExtraData string `json:"videoconferences_extra_data,omitempty"`
}

type projectEditPayload struct {
	CreationTemplate          int    `json:"creation_template,omitempty"`
	Description               string `json:"description,omitempty"`
	IsBacklogActivated        bool   `json:"is_backlog_activated,omitempty"`
	IsIssuesActivated         bool   `json:"is_issues_activated,omitempty"`
	IsKanbanActivated         bool   `json:"is_kanban_activated,omitempty"`
	IsPrivate                 bool   `json:"is_private,omitempty"`
	IsWikiActivated           bool   `json:"is_wiki_activated,omitempty"`
	Name                      string `json:"name,omitempty"`
	Videoconferences          string `json:"videoconferences,omitempty"`
	VideoconferencesExtraData string `json:"videoconferences_extra_data,omitempty"`
}

// ConfigureMappedServices maps all services to the *ProjectService with a selected project preconfigured
func (s *ProjectService) ConfigureMappedServices(ProjectID int) {
	s.Auth = &AuthService{s.client, ProjectID, "auth"}
	s.Epic = &EpicService{s.client, ProjectID, "epics"}
	s.Issue = &IssueService{s.client, ProjectID, "issues"}
	s.Milestone = &MilestoneService{s.client, ProjectID, "milestones"}
	s.Resolver = &ResolverService{s.client, ProjectID, "resolver"}
	s.Stats = &StatsService{s.client, ProjectID, "stats"}
	s.Task = &TaskService{s.client, ProjectID, "tasks"}
	s.UserStory = &UserStoryService{s.client, ProjectID, "userstories"}
	s.User = &UserService{s.client, ProjectID, "users"}
	s.Webhook = &WebhookService{s.client, ProjectID, "webhooks", "webhooklogs"}
	s.Wiki = &WikiService{s.client, ProjectID, "wiki"}
	s.Point = &PointService{s.client, ProjectID, "points"}
	s.Priority = &PriorityService{s.client, ProjectID, "priorities"}
	s.Severity = &SeverityService{s.client, ProjectID, "severities"}
	s.IssueType = &IssueTypeService{s.client, ProjectID, "issue-types"}
	s.EpicStatus = &EpicStatusService{s.client, ProjectID, "epic-statuses"}
	s.IssueStatus = &IssueStatusService{s.client, ProjectID, "issue-statuses"}
	s.TaskStatus = &TaskStatusService{s.client, ProjectID, "task-statuses"}
	s.UserStoryStatus = &UserStoryStatusService{s.client, ProjectID, "userstory-statuses"}
	s.EpicCustomAttribute = &EpicCustomAttributeService{s.client, ProjectID, "epic-custom-attributes"}
	s.IssueCustomAttribute = &IssueCustomAttributeService{s.client, ProjectID, "issue-custom-attributes"}
	s.TaskCustomAttribute = &TaskCustomAttributeService{s.client, ProjectID, "task-custom-attributes"}
	s.UserStoryCustomAttribute = &UserStoryCustomAttributeService{s.client, ProjectID, "userstory-custom-attributes"}
	s.Application = &ApplicationService{s.client, ProjectID, "applications"}
	s.ApplicationToken = &ApplicationTokenService{s.client, ProjectID, "application-tokens"}
	s.Search = &SearchService{s.client, ProjectID, "search"}
	s.UserStorage = &UserStorageService{s.client, ProjectID, "user-storage"}
	s.ProjectTemplate = &ProjectTemplateService{s.client, ProjectID, "project-templates"}
	s.ProjectTemplateDetail = &ProjectTemplateDetailService{s.client, ProjectID, "project-templates"}
	s.MembershipInvitation = &MembershipInvitationService{s.client, ProjectID, "memberships", "invitations"}
	s.WikiLink = &WikiLinkService{s.client, ProjectID, "wiki-links"}
	s.History = &HistoryService{s.client, ProjectID, "history"}
	s.NotifyPolicy = &NotifyPolicyService{s.client, ProjectID, "notify-policies"}
	s.Contact = &ContactService{s.client, ProjectID, "contact"}
	s.Feedback = &FeedbackService{s.client, ProjectID, "feedback"}
	s.ExportImport = &ExportImportService{s.client, ProjectID, "exporter", "importer"}
	s.Timeline = &TimelineService{s.client, ProjectID, "timeline"}
	s.Locale = &LocaleService{s.client, ProjectID, "locales"}
	s.Importer = &ImporterService{s.client, ProjectID, "importers"}
	s.ContribPlugin = &ContribPluginService{s.client, ProjectID, "contrib-plugins"}
	s.ObjectsSummary = &ObjectsSummaryService{s.client, ProjectID, "objects-summary"}

	s.areMappedServicesConfigured = true
}

// AreMappedServicesConfigured returns true if project-related mapped services have been configured
func (s *ProjectService) AreMappedServicesConfigured() bool {
	return s.areMappedServicesConfigured
}

// List -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-list
//
// The results can be filtered by passing in a ProjectListQueryFilter struct
func (s *ProjectService) List(queryParameters *ProjectsQueryParameters) (*ProjectsList, error) {
	/*
		The results can be filtered using the following parameters:
		  * Member
		  * Members
		  * IsLookingForPeople
		  * IsFeatured
		  * IsBacklogActivated
		  * IsKanbanActivated

		The results can be ordered using the order_by parameter with the values:
		  * memberships__user_order
		  * total_fans
		  * total_fans_last_week
		  * total_fans_last_month
		  * total_fans_last_year
		  * total_activity
		  * total_activity_last_week
		  * total_activity_last_month
		  * total_activity_last_year
	*/

	url := s.client.MakeURL(s.Endpoint)
	if queryParameters != nil {
		var err error
		url, err = appendQueryParams(url, queryParameters)
		if err != nil {
			return nil, err
		}
	}
	var projects ProjectsList

	_, err := s.client.Request.Get(url, &projects)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

// Create -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-create
// Required fields: name, description
func (s *ProjectService) Create(project *Project) (*Project, error) {
	if err := requireNonNil("project", project); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var p ProjectDetail
	// Check for required fields
	// name, description
	if isEmpty(project.Name) || isEmpty(project.Description) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	payload := projectCreatePayload{
		CreationTemplate:          project.CreationTemplate,
		Description:               project.Description,
		IsBacklogActivated:        project.IsBacklogActivated,
		IsIssuesActivated:         project.IsIssuesActivated,
		IsKanbanActivated:         project.IsKanbanActivated,
		IsPrivate:                 project.IsPrivate,
		IsWikiActivated:           project.IsWikiActivated,
		Name:                      project.Name,
		Videoconferences:          project.Videoconferences,
		VideoconferencesExtraData: project.VideoconferencesExtraData,
	}

	_, err := s.client.Request.Post(url, &payload, &p)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-get
func (s *ProjectService) Get(projectID int) (*Project, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(projectID))
	var p ProjectDetail

	_, err := s.client.Request.Get(url, &p)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// GetBySlug -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-get-by-slug
func (s *ProjectService) GetBySlug(slug string) (*Project, error) {
	queryParams := struct {
		Slug string `url:"slug"`
	}{Slug: slug}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "by_slug"), &queryParams)
	if err != nil {
		return nil, err
	}
	var p ProjectDetail

	_, err = s.client.Request.Get(url, &p)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// Edit edits an Project via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#projects-edit
// Available Meta: ProjectDetail
func (s *ProjectService) Edit(project *Project) (*Project, error) {
	if err := requireNonNil("project", project); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(project.ID))
	var p ProjectDetail

	if project.ID == 0 {
		return nil, errors.New("passed Project does not have an ID yet. Does it exist?")
	}

	payload := projectEditPayload{
		CreationTemplate:          project.CreationTemplate,
		Description:               project.Description,
		IsBacklogActivated:        project.IsBacklogActivated,
		IsIssuesActivated:         project.IsIssuesActivated,
		IsKanbanActivated:         project.IsKanbanActivated,
		IsPrivate:                 project.IsPrivate,
		IsWikiActivated:           project.IsWikiActivated,
		Name:                      project.Name,
		Videoconferences:          project.Videoconferences,
		VideoconferencesExtraData: project.VideoconferencesExtraData,
	}

	_, err := s.client.Request.Patch(url, &payload, &p)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// Update is an alias for Edit.
func (s *ProjectService) Update(project *Project) (*Project, error) {
	return s.Edit(project)
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#projects-delete
func (s *ProjectService) Delete(projectID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(projectID))
	return s.client.Request.Delete(url)
}
