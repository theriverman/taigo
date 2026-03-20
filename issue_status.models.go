package taigo

// IssueStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#issue-statuses
type IssueStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

// IssueStatusCreateRequest represents payload for creating issue statuses.
type IssueStatusCreateRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project"`
}

// IssueStatusEditRequest represents sparse non-destructive updates for issue statuses.
type IssueStatusEditRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name,omitempty"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project,omitempty"`
}

// IssueStatusPatch represents explicit PATCH payload for issue statuses.
type IssueStatusPatch struct {
	Color    *string `json:"color,omitempty"`
	IsClosed *bool   `json:"is_closed,omitempty"`
	Name     *string `json:"name,omitempty"`
	Order    *int    `json:"order,omitempty"`
	Project  *int    `json:"project,omitempty"`
}
