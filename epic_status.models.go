package taigo

// EpicStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#epic-statuses
type EpicStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}

// EpicStatusCreateRequest represents payload for creating epic statuses.
type EpicStatusCreateRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project"`
}

// EpicStatusEditRequest represents sparse non-destructive updates for epic statuses.
type EpicStatusEditRequest struct {
	Color    string `json:"color,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
	Name     string `json:"name,omitempty"`
	Order    int    `json:"order,omitempty"`
	Project  int    `json:"project,omitempty"`
}

// EpicStatusPatch represents explicit PATCH payload for epic statuses.
type EpicStatusPatch struct {
	Color    *string `json:"color,omitempty"`
	IsClosed *bool   `json:"is_closed,omitempty"`
	Name     *string `json:"name,omitempty"`
	Order    *int    `json:"order,omitempty"`
	Project  *int    `json:"project,omitempty"`
}
