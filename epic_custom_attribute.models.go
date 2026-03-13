package taigo

import "time"

// EpicCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attribute-detail
type EpicCustomAttribute struct {
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

// EpicCustomAttributeCreateRequest represents payload for creating epic custom attributes.
type EpicCustomAttributeCreateRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project"`
	Type        string `json:"type"`
}

// EpicCustomAttributeEditRequest represents sparse non-destructive updates.
type EpicCustomAttributeEditRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name,omitempty"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project,omitempty"`
	Type        string `json:"type,omitempty"`
}

// EpicCustomAttributePatch represents explicit PATCH payload for epic custom attributes.
type EpicCustomAttributePatch struct {
	Description *string `json:"description,omitempty"`
	Extra       any     `json:"extra,omitempty"`
	Name        *string `json:"name,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Project     *int    `json:"project,omitempty"`
	Type        *string `json:"type,omitempty"`
}
