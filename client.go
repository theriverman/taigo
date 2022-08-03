package taigo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	// Token Refresh Helpers
	TokenRefreshTicker *time.Ticker
	tokenRefreshDone   chan bool
}

// MakeURL accepts an Endpoint URL and returns a compiled absolute URL
//
// For example:
//	* If the given endpoint URLs are [epics, attachments]
//	* If the BaseURL is https://api.taiga.io
//	* It returns https://api.taiga.io/api/v1/epics/attachments
//  * Suffixes are appended to the URL joined by a slash (/)
func (c *Client) MakeURL(EndpointParts ...string) string {
	return c.APIURL + "/" + strings.Join(EndpointParts, "/")
}

// Initialise returns a new Taiga Client which is the entrypoint of the driver
// Initialise() is automatically called by the `AuthByCredentials` and `AuthByToken` methods.
// If you, for some reason, would like to manually set the Client.Token field, then Initialise() must be called manually!
func (c *Client) Initialise() error {
	// Skip if already Initialised
	if c.isInitialised {
		return nil
	}
	// Taiga.Client safety guards
	if len(c.BaseURL) < len("http://") { // compares for a minimum of len("http://")
		return fmt.Errorf("BaseURL is not set or invalid")

	}
	//Set basic token type
	if len(c.TokenType) <= 1 {
		c.TokenType = "Bearer"
	}
	// Set basic headers
	c.Headers = &http.Header{}
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

	c.isInitialised = true

	// Token Refresh Routine
	if c.AutoRefreshDisabled {
		if c.Verbose {
			log.Println("automatic token refresh subroutine will not be started because AutoRefreshDisabled = false")
		}
		return nil
	}
	if c.AutoRefreshTickerDuration == 0 {
		/*
			https://github.com/kaleidos-ventures/taiga-back/blob/0be90e6a661de51bf9e95744322060f33dafa347/taiga/auth/settings.py#L50
			According to the base settings in Taiga, the default `REFRESH_TOKEN_LIFETIME` is `timedelta(days=1)`.
			Let's play safe, and take only half of that, so we refresh our tokens every 12 hours.
		*/
		c.AutoRefreshTickerDuration = 12 * time.Hour
	}
	if c.Verbose {
		log.Printf("AutoRefreshTickerDuration: %s\n", c.AutoRefreshTickerDuration)
	}
	c.tokenRefreshDone = make(chan bool)
	if c.TokenRefreshTicker == nil {
		c.TokenRefreshTicker = time.NewTicker(c.AutoRefreshTickerDuration)
	}
	if c.RefreshTokenRoutine == nil {
		c.RefreshTokenRoutine = defaultTokenRefreshRoutine
	}
	c.RefreshTokenRoutine(c, c.TokenRefreshTicker) // calling the Token Refresh Routine
	return nil
}

func (c *Client) DisableAutomaticTokenRefresh() {
	c.TokenRefreshTicker.Stop()
	c.tokenRefreshDone <- true
	if c.Verbose {
		log.Println("automatic token refresh has been disabled")
	}
}

// AuthByCredentials authenticates to Taiga using the provided basic credentials
func (c *Client) AuthByCredentials(credentials *Credentials) error {
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
	c.TokenType = tokenType
	c.Token = token
	c.RefreshToken = refreshToken
	c.setToken() // Add to headers

	var err error
	c.Self, err = c.User.Me()
	if err != nil {
		return fmt.Errorf("authentication failed: %s", err)
	}
	return nil
}

// DisablePagination controls the value of header `x-disable-pagination`.
func (c *Client) DisablePagination(b bool) {
	m := map[string]string{
		"x-disable-pagination": cases.Title(language.BritishEnglish).String(strconv.FormatBool(b)),
	}
	c.LoadExternalHeaders(m)
	c.paginationDisabled = b
}

// GetPagination returns the Pagination struct created from the last response
func (c *Client) GetPagination() Pagination {
	return *c.pagination
}

// GetAuthorizationHeader returns the formatted value of Authorization key from Headers
func (c *Client) GetAuthorizationHeader() string {
	return c.Headers.Get("Authorization")
}

// LoadExternalHeaders loads a map of header key/value pairs permemently into `Client.Headers`
func (c *Client) LoadExternalHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers.Add(k, v)
	}
}

func (c *Client) setContentTypeToJSON() {
	c.Headers.Add("Content-Type", "application/json")
}

func (c *Client) setToken() {
	c.Headers.Del("Authorization") // avoid header duplication
	c.Headers.Add("Authorization", c.TokenType+" "+c.Token)
}

// loadHeaders takes an http.Request and maps locally stored Header values to its .Header
func (c *Client) loadHeaders(request *http.Request) {
	for key, values := range *c.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
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
				if refreshData, err := c.Auth.RefreshAuthToken(true); err != nil {
					log.Println(err)
					return
				} else {
					c.Token = refreshData.AuthToken
					c.RefreshToken = refreshData.Refresh
					c.setToken()
				}
			}
		}
	}()
}
