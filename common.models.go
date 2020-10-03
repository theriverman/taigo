package taigo

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AgilePoints is a [string/int] key/value pair to represent agile points in a UserStory, Milestone, etc...
//
// JSON Representation example:
/*
	"points": {
		"1": 12,
		"2": 2,
		"3": 5,
		"4": 5
	},
*/
type AgilePoints map[string]float64

// Points represent the Agile Points configured for the project and set for respective Taiga object
type Points = AgilePoints

// Tags represent the tags slice of respective Taiga object
type Tags [][]string

// Attachment => https://taigaio.github.io/taiga-doc/dist/api.html#object-attachment-detail
type Attachment struct {
	AttachedFile     string    `json:"attached_file,omitempty"`
	CreatedDate      time.Time `json:"created_date,omitempty"`
	Description      string    `json:"description,omitempty"`
	FromComment      bool      `json:"from_comment,omitempty"`
	ID               int       `json:"id,omitempty"`
	IsDeprecated     bool      `json:"is_deprecated,omitempty"`
	ModifiedDate     time.Time `json:"modified_date,omitempty"`
	Name             string    `json:"name,omitempty"`
	ObjectID         int       `json:"object_id,omitempty"`
	Order            int       `json:"order,omitempty"`
	Owner            int       `json:"owner,omitempty"`
	PreviewURL       string    `json:"preview_url,omitempty"`
	Project          int       `json:"project,omitempty"`
	Sha1             string    `json:"sha1,omitempty"`
	Size             int       `json:"size,omitempty"`
	ThumbnailCardURL string    `json:"thumbnail_card_url,omitempty"`
	URL              string    `json:"url,omitempty"`
	filePath         string    // For package-internal use only
}

// GenericObjectAttachment represents an array of minimal attachment details
// This array is filled when the `IncludeAttachments` query parameter is true
type GenericObjectAttachment struct {
	AttachedFile     string `json:"attached_file,omitempty"`
	ID               int    `json:"id,omitempty"`
	ThumbnailCardURL string `json:"thumbnail_card_url,omitempty"`
}

// SetFilePath takes the path to the file be uploaded
func (a *Attachment) SetFilePath(FilePath string) {
	a.filePath = FilePath
}

// attachmentsQueryParams is a helper to transfer and render ObjectID and Project ID as URL query parameters
type attachmentsQueryParams struct {
	ObjectID    int    `url:"object_id,omitempty"`
	Project     int    `url:"project,omitempty"`
	endpointURI string // unexported to exclude it from `go-querystring/query.Values()`
}

// Neighbors represents a read-only field
type Neighbors struct {
	Next     Next     `json:"next"`
	Previous Previous `json:"previous"`
}

// Next represents a read-only field
type Next struct {
	ID      int    `json:"id"`
	Ref     int    `json:"ref"`
	Subject string `json:"subject"`
}

// Previous represents a read-only field
type Previous = Next

// Owner represents the owner of an object
type Owner struct {
	BigPhoto        string `json:"big_photo,omitempty"`
	FullNameDisplay string `json:"full_name_display,omitempty"`
	GravatarID      string `json:"gravatar_id,omitempty"`
	ID              int    `json:"id,omitempty"`
	IsActive        bool   `json:"is_active,omitempty"`
	Photo           string `json:"photo,omitempty"`
	Username        string `json:"username,omitempty"`
}

// TagsColors represent color code and color name combinations for a tag
type TagsColors struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

// UserStoriesCounts represents the number of userStories
type UserStoriesCounts struct {
	Progress int `json:"progress,omitempty"`
	Total    int `json:"total,omitempty"`
}

/*
	READ-ONLY FIELDS

	All  fields ending in `_extra_info` (`assigned_to_extra_info`, `is_private_extra_info`, `owner_extra_info`, `project_extra_info`,
	`status_extra_info`, `status_extra_info`, `user_story_extra_info`…​) are read-only fields

	https://taigaio.github.io/taiga-doc/dist/api.html#_read_only_fields
*/

// AssignedToExtraInfo is a read-only field
type AssignedToExtraInfo = Owner

// IsPrivateExtraInfo is a read-only field
type IsPrivateExtraInfo struct {
	Reason       string `json:"reason,omitempty"`
	CanBeUpdated bool   `json:"can_be_updated,omitempty"`
}

// StatusExtraInfo is a read-only field
type StatusExtraInfo struct {
	Color    string `json:"color"`
	IsClosed bool   `json:"is_closed"`
	Name     string `json:"name"`
}

// OwnerExtraInfo is a read-only field
type OwnerExtraInfo = Owner

// ProjectExtraInfo represents a read-only field
type ProjectExtraInfo struct {
	ID           int    `json:"id"`
	LogoSmallURL string `json:"logo_small_url"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
}

// UserStoryExtraInfo is a read-only field
type UserStoryExtraInfo struct {
	Epics   []EpicMinimal `json:"epics"`
	ID      int           `json:"id"`
	Ref     int           `json:"ref"`
	Subject string        `json:"subject"`
}

// Pagination represents the information returned via headers
//
// https://taigaio.github.io/taiga-doc/dist/api.html#_pagination
type Pagination struct {
	Paginated         bool     // indicating if pagination is being used for the request
	PaginatedBy       int      // number of results per page
	PaginationCount   int      // total number of results
	PaginationCurrent int      // current page
	PaginationNext    *url.URL // next results
	PaginationPrev    *url.URL // previous results
}

// LoadFromHeaders accepts an *http.Response struct and reads the relevant
// pagination headers returned by Taiga
func (p *Pagination) LoadFromHeaders(c *Client, response *http.Response) {
	paginated := response.Header.Get("X-Paginated") // Check if response is paginated
	if paginated == "true" {
		p.Paginated = true
		p.PaginatedBy, _ = strconv.Atoi(response.Header.Get("X-Paginated-By"))
		p.PaginationCount, _ = strconv.Atoi(response.Header.Get("X-Paginated-Count"))
		p.PaginationCurrent, _ = strconv.Atoi(response.Header.Get("X-Paginated-Current"))
		p.PaginationNext, _ = url.Parse(response.Header.Get("X-Pagination-Next"))
		p.PaginationPrev, _ = url.Parse(response.Header.Get("X-Pagination-Prev"))
	} else {
		p.Paginated = false
	}
}
