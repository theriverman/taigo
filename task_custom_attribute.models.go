package taigo

import "time"

// TaskCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#task-custom-attributes
type TaskCustomAttribute struct {
	CreatedDate  time.Time `json:"created_date"`
	Description  string    `json:"description"`
	Extra        any       `json:"extra"`
	ID           int       `json:"id"`
	ModifiedDate time.Time `json:"modified_date"`
	Name         string    `json:"name"`
	Order        int       `json:"order"`
	ProjectID    int       `json:"project_id"`
	Type         string    `json:"type"`
}

// TaskCustomAttributeCreateRequest represents payload for creating task custom attributes.
type TaskCustomAttributeCreateRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project"`
	Type        string `json:"type"`
}

// TaskCustomAttributeEditRequest represents sparse non-destructive updates.
type TaskCustomAttributeEditRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name,omitempty"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project,omitempty"`
	Type        string `json:"type,omitempty"`
}

// TaskCustomAttributePatch represents explicit PATCH payload for task custom attributes.
type TaskCustomAttributePatch struct {
	Description *string `json:"description,omitempty"`
	Extra       any     `json:"extra,omitempty"`
	Name        *string `json:"name,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Project     *int    `json:"project,omitempty"`
	Type        *string `json:"type,omitempty"`
}
