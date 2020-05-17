package gotaiga

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Client is the session manager of Taiga Driver
type Client struct {
	Credentials        *Credentials
	APIURL             string                    // set by system
	APIversion         string                    // default: "v1"
	BaseURL            string                    // i.e.: "http://taiga.test" | Same value as `api` in `taiga-front-dist/dist/conf.json`
	Headers            map[string]string         // set by system&user
	HTTPClient         *http.Client              // set by user
	IsLoggedIn         bool                      // set by system
	Logger             *log.Logger               // customisable logging interface
	LoginType          string                    // i.e.: "normal"; "github"; "ldap"
	Token              string                    // set by system; can be set manually
	TokenType          string                    // default=Bearer; options:Bearer,Application
	SelfUser           *UserAuthenticationDetail // User logged in
	defaultProjectID   int                       // Project used in query params by default
	pagination         *Pagination               // Pagination details extracted from the LAST http response
	paginationDisabled bool                      // indicates pagination status

	// Services
	Auth      *AuthService
	Epic      *EpicService
	Issue     *IssueService
	Milestone *MilestoneService
	Project   *ProjectService
	Task      *TaskService
	UserStory *UserStoryService
	User      *UserService
	Webhook   *WebhookService
}

// TODO: Pack Taiga operations into services, such as, ProjectService, EpicService, MilestoneService, etc...
// These services should be available under Client and implement their actions (List, Get, Create, etc...) there!

// SetDefaultProject takes an int and uses that to internally set the default project ID.
func (c *Client) SetDefaultProject(projectID int) error {
	if !(projectID > 0) {
		msg := fmt.Sprintf("Could not set Default Project ID. Provided projectID was: %d", projectID)
		c.Logger.Fatalln(msg)
		return NewError(msg)
	}
	c.defaultProjectID = projectID
	return nil
}

// SetDefaultProjectBySlug takes an slug string and uses that to internally set the default project ID.
func (c *Client) SetDefaultProjectBySlug(projectSlug string) error {
	if projectSlug == "" {
		msg := fmt.Sprintf("Could not set Default Project ID. Provided projectSlug was: %s", projectSlug)
		c.Logger.Fatalln(msg)
		return NewError(msg)
	}
	proj, err := c.Project.GetBySlug(projectSlug)
	if err != nil {
		return err
	}
	c.defaultProjectID = proj.ID
	return nil
}

// GetDefaultProject returns the currently set Default Project ID
func (c *Client) GetDefaultProject() int {
	return c.defaultProjectID
}

// GetDefaultProjectAsQueryParam returns the currently set Default Project ID formatted as a QueryParam
func (c *Client) GetDefaultProjectAsQueryParam() string {
	return fmt.Sprintf("?project=%d", c.defaultProjectID)
}

// ClearDefaultProject resets the currently set Default Project ID to 0 (None)
func (c *Client) ClearDefaultProject() {
	c.defaultProjectID = 0
}

// HasDefaultProject returns true if there's a Default Project ID set
func (c *Client) HasDefaultProject() bool {
	if c.defaultProjectID > 0 {
		return true
	}
	return false
}

// Initialise returns a new Taiga Client which is the entrypoint of the driver
func (c *Client) Initialise(credentials *Credentials) error {
	// Instantiate a logger handle
	if c.Logger == nil {
		c.instantiateLogger()
	}

	// Taiga.Client safety guards
	if len(c.BaseURL) < len("http://") { // compares for a minimum of len("http://")
		return errors.New("BaseURL is not set or invalid")
	}
	if len(c.LoginType) <= 1 {
		return errors.New("LoginType is not set")
	}
	//Set basic token type
	if len(c.TokenType) <= 1 {
		c.TokenType = "Bearer"
	}
	// Set basic headers
	c.Headers = make(map[string]string)
	c.setContentTypeToJSON()

	// Setup APIversion | If user did not set it to anything, default to `v1`
	if c.APIversion == "" {
		c.APIversion = "v1"
	}
	// Compile URL to Taiga API | Example: https://api.taiga.io/api/v1
	c.APIURL = c.BaseURL + "/api/" + c.APIversion

	// Disable pagination ()
	c.DisablePagination(true) // https://taigaio.github.io/taiga-doc/dist/api.html#_pagination

	// Bootstrapping Services
	c.Auth = &AuthService{c}
	c.Epic = &EpicService{c}
	c.Issue = &IssueService{c}
	c.Milestone = &MilestoneService{c}
	c.Project = &ProjectService{c}
	c.Task = &TaskService{c}
	c.UserStory = &UserStoryService{c}
	c.User = &UserService{c}
	c.Webhook = &WebhookService{c}

	user, err := c.Auth.login(credentials)
	if err != nil {
		return err
	}
	c.SelfUser = user

	return nil
}

// DisablePagination controls the value of header `x-disable-pagination`.
func (c *Client) DisablePagination(b bool) {
	var decision string = strings.Title(strconv.FormatBool(b))
	m := map[string]string{
		"x-disable-pagination": decision,
	}
	c.LoadExternalHeaders(m)
	c.paginationDisabled = b
}

// GetPagination returns the Pagination struct created from the last response
func (c *Client) GetPagination() Pagination {
	return *c.pagination
}

// LoadExternalHeaders loads a map of header key/value pairs permemently into `Client.Headers`
func (c *Client) LoadExternalHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers[k] = v
	}
}

func (c *Client) setContentTypeToJSON() {
	c.Headers["Content-Type"] = "application/json"
}

func (c *Client) setContentTypeToFormData() {
	c.Headers["Content-Type"] = "multipart/form-data"
}

func (c *Client) setToken() {
	c.Headers["Authorization"] = c.TokenType + " " + c.Token
}

func (c *Client) loadHeaders(request *http.Request) {
	for k, v := range c.Headers {
		request.Header.Set(k, v)
	}
}

func (c *Client) setContentType(s string) {
	c.Headers["Content-Type"] = s
}

func (c *Client) instantiateLogger() {
	c.Logger = &log.Logger{}
	// Configure Log
	c.Logger.SetOutput(os.Stdout)
	c.Logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	c.Logger.SetPrefix("go-taiga >> ")
}
