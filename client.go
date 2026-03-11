package taigo

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TokenBearer is the standard token type for authentication in Taiga
const TokenBearer string = "Bearer"

// TokenApplication is the Token type used for external apps
// These tokens are associated to an existing user and an Application. They can be manually created via the Django ADMIN or programatically via API
// They work in the same way than standard Taiga authentication tokens but the "Authorization" header change slightly.
const TokenApplication string = "Application"

// Client is the session manager of Taiga Driver
type Client struct {
	Credentials               *Credentials
	APIURL                    string                               // set by system
	APIversion                string                               // default: "v1"
	BaseURL                   string                               // i.e.: "http://taiga.test" | Same value as `api` in `taiga-front-dist/dist/conf.json`
	Headers                   *http.Header                         // mostly set by system
	HTTPClient                *http.Client                         // set by user
	Token                     string                               // set by system; can be set manually
	TokenType                 string                               // default=Bearer; options:Bearer,Application
	RefreshToken              string                               // set by system; can be set manually
	RefreshTokenRoutine       func(c *Client, ticker *time.Ticker) // routine periodically refreshing the token
	Self                      *User                                // User logged in
	pagination                *Pagination                          // Pagination details extracted from the LAST http response
	paginationDisabled        bool                                 // indicates pagination status
	isInitialised             bool                                 // indicates if taiga.Client has been initialised already
	Verbose                   bool                                 // internal Taigo events are logged in a more verbose fashion
	AutoRefreshDisabled       bool                                 // if true before initialisation, RefreshTokenRoutine never gets called
	AutoRefreshTickerDuration time.Duration                        // time.Duration between two token refresh requests

	// Core Services
	Request *RequestService

	// Taiga Services
	Auth      *AuthService
	Epic      *EpicService
	Issue     *IssueService
	Milestone *MilestoneService
	Project   *ProjectService
	Resolver  *ResolverService
	Stats     *StatsService
	Task      *TaskService
	UserStory *UserStoryService
	User      *UserService
	Webhook   *WebhookService
	Wiki      *WikiService
	Point     *PointService
	Priority  *PriorityService
	Severity  *SeverityService
	IssueType *IssueTypeService

	EpicStatus      *EpicStatusService
	IssueStatus     *IssueStatusService
	TaskStatus      *TaskStatusService
	UserStoryStatus *UserStoryStatusService

	EpicCustomAttribute      *EpicCustomAttributeService
	IssueCustomAttribute     *IssueCustomAttributeService
	TaskCustomAttribute      *TaskCustomAttributeService
	UserStoryCustomAttribute *UserStoryCustomAttributeService
	Application              *ApplicationService
	ApplicationToken         *ApplicationTokenService
	Search                   *SearchService
	UserStorage              *UserStorageService
	ProjectTemplate          *ProjectTemplateService
	ProjectTemplateDetail    *ProjectTemplateDetailService
	MembershipInvitation     *MembershipInvitationService
	WikiLink                 *WikiLinkService
	History                  *HistoryService
	NotifyPolicy             *NotifyPolicyService
	Contact                  *ContactService
	Feedback                 *FeedbackService
	ExportImport             *ExportImportService
	Timeline                 *TimelineService
	Locale                   *LocaleService
	Importer                 *ImporterService
	ContribPlugin            *ContribPluginService
	ObjectsSummary           *ObjectsSummaryService

	// Token Refresh Helpers
	TokenRefreshTicker *time.Ticker
	tokenRefreshDone   chan bool

	initMu  sync.Mutex
	stateMu sync.RWMutex
}

// MakeURL accepts an Endpoint URL and returns a compiled absolute URL
//
// For example:
//   - If the given endpoint URLs are [epics, attachments]
//   - If the BaseURL is https://api.taiga.io
//   - It returns https://api.taiga.io/api/v1/epics/attachments
//   - Suffixes are appended to the URL joined by a slash (/)
func (c *Client) MakeURL(EndpointParts ...string) string {
	return c.APIURL + "/" + strings.Join(EndpointParts, "/")
}

// Initialise returns a new Taiga Client which is the entrypoint of the driver
// Initialise() is automatically called by the `AuthByCredentials` and `AuthByToken` methods.
// If you, for some reason, would like to manually set the Client.Token field, then Initialise() must be called manually!
func (c *Client) Initialise() error {
	c.initMu.Lock()
	defer c.initMu.Unlock()

	// Skip if already Initialised
	if c.isInitialised {
		return nil
	}
	// Taiga.Client safety guards
	c.BaseURL = strings.TrimRight(c.BaseURL, "/")
	parsedBaseURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil || parsedBaseURL.Scheme == "" || parsedBaseURL.Host == "" {
		return fmt.Errorf("BaseURL is not set or invalid")
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	// Set basic token type
	if len(c.TokenType) <= 1 {
		c.TokenType = TokenBearer
	}
	// Set basic headers
	c.stateMu.Lock()
	c.Headers = &http.Header{}
	c.stateMu.Unlock()
	c.setContentTypeToJSON() // Default Header = {"Content-Type": "application/json"}

	// Setup APIversion | If user did not set it to anything, default to `v1`
	if c.APIversion == "" {
		c.APIversion = "v1"
	}
	// Compile URL to Taiga API | Example: https://api.taiga.io/api/v1
	c.APIURL = c.BaseURL + "/api/" + c.APIversion

	// Disable pagination ()
	c.DisablePagination(true) // https://taigaio.github.io/taiga-doc/dist/api.html#_pagination

	// Bootstrapping Services
	c.Request = &RequestService{c}
	c.Auth = &AuthService{c, 0, "auth"}
	c.Epic = &EpicService{c, 0, "epics"}
	c.Issue = &IssueService{c, 0, "issues"}
	c.Milestone = &MilestoneService{c, 0, "milestones"}
	c.Project = &ProjectService{client: c, Endpoint: "projects"}
	c.Resolver = &ResolverService{c, 0, "resolver"}
	c.Stats = &StatsService{c, 0, "stats"}
	c.Task = &TaskService{c, 0, "tasks"}
	c.UserStory = &UserStoryService{c, 0, "userstories"}
	c.User = &UserService{c, 0, "users"}
	c.Webhook = &WebhookService{c, 0, "webhooks", "webhooklogs"}
	c.Wiki = &WikiService{c, 0, "wiki"}
	c.Point = &PointService{c, 0, "points"}
	c.Priority = &PriorityService{c, 0, "priorities"}
	c.Severity = &SeverityService{c, 0, "severities"}
	c.IssueType = &IssueTypeService{c, 0, "issue-types"}
	c.EpicStatus = &EpicStatusService{c, 0, "epic-statuses"}
	c.IssueStatus = &IssueStatusService{c, 0, "issue-statuses"}
	c.TaskStatus = &TaskStatusService{c, 0, "task-statuses"}
	c.UserStoryStatus = &UserStoryStatusService{c, 0, "userstory-statuses"}
	c.EpicCustomAttribute = &EpicCustomAttributeService{c, 0, "epic-custom-attributes"}
	c.IssueCustomAttribute = &IssueCustomAttributeService{c, 0, "issue-custom-attributes"}
	c.TaskCustomAttribute = &TaskCustomAttributeService{c, 0, "task-custom-attributes"}
	c.UserStoryCustomAttribute = &UserStoryCustomAttributeService{c, 0, "userstory-custom-attributes"}
	c.Application = &ApplicationService{c, 0, "applications"}
	c.ApplicationToken = &ApplicationTokenService{c, 0, "application-tokens"}
	c.Search = &SearchService{c, 0, "search"}
	c.UserStorage = &UserStorageService{c, 0, "user-storage"}
	c.ProjectTemplate = &ProjectTemplateService{c, 0, "project-templates"}
	c.ProjectTemplateDetail = &ProjectTemplateDetailService{c, 0, "project-templates"}
	c.MembershipInvitation = &MembershipInvitationService{c, 0, "memberships", "invitations"}
	c.WikiLink = &WikiLinkService{c, 0, "wiki-links"}
	c.History = &HistoryService{c, 0, "history"}
	c.NotifyPolicy = &NotifyPolicyService{c, 0, "notify-policies"}
	c.Contact = &ContactService{c, 0, "contact"}
	c.Feedback = &FeedbackService{c, 0, "feedback"}
	c.ExportImport = &ExportImportService{c, 0, "exporter", "importer"}
	c.Timeline = &TimelineService{c, 0, "timeline"}
	c.Locale = &LocaleService{c, 0, "locales"}
	c.Importer = &ImporterService{c, 0, "importers"}
	c.ContribPlugin = &ContribPluginService{c, 0, "contrib-plugins"}
	c.ObjectsSummary = &ObjectsSummaryService{c, 0, "objects-summary"}

	c.isInitialised = true

	// Token Refresh Routine
	c.stateMu.RLock()
	autoRefreshDisabled := c.AutoRefreshDisabled
	autoRefreshTickerDuration := c.AutoRefreshTickerDuration
	verbose := c.Verbose
	c.stateMu.RUnlock()
	if autoRefreshDisabled {
		if verbose {
			log.Println("automatic token refresh subroutine will not be started because AutoRefreshDisabled = true")
		}
		return nil
	}
	if autoRefreshTickerDuration == 0 {
		/*
			https://github.com/kaleidos-ventures/taiga-back/blob/0be90e6a661de51bf9e95744322060f33dafa347/taiga/auth/settings.py#L50
			According to the base settings in Taiga, the default `REFRESH_TOKEN_LIFETIME` is `timedelta(days=1)`.
			Let's play safe, and take only half of that, so we refresh our tokens every 12 hours.
		*/
		autoRefreshTickerDuration = 12 * time.Hour
		c.stateMu.Lock()
		c.AutoRefreshTickerDuration = autoRefreshTickerDuration
		c.stateMu.Unlock()
	}
	if verbose {
		log.Printf("AutoRefreshTickerDuration: %s\n", autoRefreshTickerDuration)
	}
	// Buffered channel avoids blocking on repeated disable calls.
	c.stateMu.Lock()
	c.tokenRefreshDone = make(chan bool, 1)
	if c.TokenRefreshTicker == nil {
		c.TokenRefreshTicker = time.NewTicker(autoRefreshTickerDuration)
	}
	if c.RefreshTokenRoutine == nil {
		c.RefreshTokenRoutine = defaultTokenRefreshRoutine
	}
	ticker := c.TokenRefreshTicker
	refreshRoutine := c.RefreshTokenRoutine
	c.stateMu.Unlock()
	refreshRoutine(c, ticker) // calling the Token Refresh Routine
	return nil
}

func (c *Client) DisableAutomaticTokenRefresh() {
	c.stateMu.Lock()
	c.AutoRefreshDisabled = true
	ticker := c.TokenRefreshTicker
	done := c.tokenRefreshDone
	c.stateMu.Unlock()

	if ticker != nil {
		ticker.Stop()
	}
	if done != nil {
		select {
		case done <- true:
		default:
		}
	}
	if c.Verbose {
		log.Println("automatic token refresh has been disabled")
	}
}

// AuthByCredentials authenticates to Taiga using the provided basic credentials
func (c *Client) AuthByCredentials(credentials *Credentials) error {
	if err := requireNonNil("credentials", credentials); err != nil {
		return err
	}
	if err := c.Initialise(); err != nil {
		return err
	}

	if len(credentials.Type) <= 1 {
		return fmt.Errorf("LoginType is not set")
	}

	user, err := c.Auth.login(credentials)
	if err != nil {
		return err
	}

	c.Self = user.AsUser()
	return nil
}

// AuthByToken authenticates to Taiga using provided Token by requesting users/me
func (c *Client) AuthByToken(tokenType, token, refreshToken string) error {
	if err := c.Initialise(); err != nil {
		return err
	}
	if tokenType == "" {
		c.stateMu.RLock()
		tokenType = c.TokenType
		c.stateMu.RUnlock()
		if tokenType == "" {
			tokenType = TokenBearer
		}
	}
	c.setAuthTokens(tokenType, token, refreshToken) // Also adds Authorization header.

	var err error
	c.Self, err = c.User.Me()
	if err != nil {
		return fmt.Errorf("authentication failed: %s", err)
	}
	return nil
}

// DisablePagination controls the value of header `x-disable-pagination`.
func (c *Client) DisablePagination(b bool) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.ensureHeadersLocked()
	const headerName = "x-disable-pagination"
	if b {
		c.Headers.Set(headerName, cases.Title(language.BritishEnglish).String(strconv.FormatBool(true)))
	} else {
		// Taiga checks only header presence; remove the header to enable pagination.
		c.Headers.Del(headerName)
	}
	c.paginationDisabled = b
}

// GetPagination returns the Pagination struct created from the last response
func (c *Client) GetPagination() Pagination {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	if c.pagination == nil {
		return Pagination{}
	}
	return *c.pagination
}

// GetAuthorizationHeader returns the formatted value of Authorization key from Headers
func (c *Client) GetAuthorizationHeader() string {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	if c.Headers == nil {
		return ""
	}
	return c.Headers.Get("Authorization")
}

// LoadExternalHeaders loads a map of header key/value pairs permemently into `Client.Headers`
func (c *Client) LoadExternalHeaders(headers map[string]string) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.ensureHeadersLocked()
	for k, v := range headers {
		c.Headers.Set(k, v)
	}
}

func (c *Client) setContentTypeToJSON() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.ensureHeadersLocked()
	c.Headers.Set("Content-Type", "application/json")
}

func (c *Client) setToken() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.ensureHeadersLocked()
	if c.TokenType == "" {
		c.TokenType = TokenBearer
	}
	c.Headers.Set("Authorization", c.TokenType+" "+c.Token)
}

// loadHeaders takes an http.Request and maps locally stored Header values to its .Header
func (c *Client) loadHeaders(request *http.Request) {
	c.stateMu.RLock()
	headers := http.Header{}
	if c.Headers != nil {
		headers = c.Headers.Clone()
	}
	c.stateMu.RUnlock()
	for key, values := range headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
}

func (c *Client) setPagination(pagination *Pagination) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.pagination = pagination
}

func (c *Client) currentTokens() (token string, refreshToken string) {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.Token, c.RefreshToken
}

func (c *Client) setAuthTokens(tokenType, token, refreshToken string) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if tokenType != "" {
		c.TokenType = tokenType
	}
	if c.TokenType == "" {
		c.TokenType = TokenBearer
	}
	c.Token = token
	c.RefreshToken = refreshToken
	c.ensureHeadersLocked()
	c.Headers.Set("Authorization", c.TokenType+" "+c.Token)
}

func (c *Client) ensureHeadersLocked() {
	if c.Headers == nil {
		c.Headers = &http.Header{}
	}
}

func defaultTokenRefreshRoutine(c *Client, ticker *time.Ticker) {
	go func() {
		for {
			select {
			case <-c.tokenRefreshDone:
				if c.Verbose {
					log.Println("TokenRefreshRoutine has been stopped")
				}
				return
			case t := <-ticker.C:
				if c.Verbose {
					log.Println("TokenRefreshRoutine tick at", t, "-> Refreshing the stored tokens")
				}
				if _, err := c.Auth.RefreshAuthToken(true); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
}
