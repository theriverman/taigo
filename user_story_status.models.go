package taigo

// UserStoryStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#user-story-statuses
type UserStoryStatus struct {
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

// UserStoryStatusCreateRequest represents payload for creating user story statuses.
type UserStoryStatusCreateRequest struct {
	Color      string `json:"color,omitempty"`
	IsArchived bool   `json:"is_archived,omitempty"`
	IsClosed   bool   `json:"is_closed,omitempty"`
	Name       string `json:"name"`
	Order      int    `json:"order,omitempty"`
	Project    int    `json:"project"`
	WipLimit   int    `json:"wip_limit,omitempty"`
}

// UserStoryStatusEditRequest represents sparse non-destructive updates for user story statuses.
type UserStoryStatusEditRequest struct {
	Color      string `json:"color,omitempty"`
	IsArchived bool   `json:"is_archived,omitempty"`
	IsClosed   bool   `json:"is_closed,omitempty"`
	Name       string `json:"name,omitempty"`
	Order      int    `json:"order,omitempty"`
	Project    int    `json:"project,omitempty"`
	WipLimit   int    `json:"wip_limit,omitempty"`
}

// UserStoryStatusPatch represents explicit PATCH payload for user story statuses.
type UserStoryStatusPatch struct {
	Color      *string `json:"color,omitempty"`
	IsArchived *bool   `json:"is_archived,omitempty"`
	IsClosed   *bool   `json:"is_closed,omitempty"`
	Name       *string `json:"name,omitempty"`
	Order      *int    `json:"order,omitempty"`
	Project    *int    `json:"project,omitempty"`
	WipLimit   *int    `json:"wip_limit,omitempty"`
}
