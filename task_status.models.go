package taigo

// TaskStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#task-statuses
type TaskStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

// TaskStatusCreateRequest represents payload for creating task statuses.
type TaskStatusCreateRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project"`
}

// TaskStatusEditRequest represents sparse non-destructive updates for task statuses.
type TaskStatusEditRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name,omitempty"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project,omitempty"`
}

// TaskStatusPatch represents explicit PATCH payload for task statuses.
type TaskStatusPatch struct {
	Color    *string `json:"color,omitempty"`
	IsClosed *bool   `json:"is_closed,omitempty"`
	Name     *string `json:"name,omitempty"`
	Order    *int    `json:"order,omitempty"`
	Project  *int    `json:"project,omitempty"`
}
