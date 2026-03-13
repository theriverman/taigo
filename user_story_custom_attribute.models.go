package taigo

import "time"

// UserStoryCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#user-story-custom-attributes-list
type UserStoryCustomAttribute struct {
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

// UserStoryCustomAttributeCreateRequest represents payload for creating user story custom attributes.
type UserStoryCustomAttributeCreateRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project"`
	Type        string `json:"type"`
}

// UserStoryCustomAttributeEditRequest represents sparse non-destructive updates.
type UserStoryCustomAttributeEditRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name,omitempty"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project,omitempty"`
	Type        string `json:"type,omitempty"`
}

// UserStoryCustomAttributePatch represents explicit PATCH payload for user story custom attributes.
type UserStoryCustomAttributePatch struct {
	Description *string `json:"description,omitempty"`
	Extra       any     `json:"extra,omitempty"`
	Name        *string `json:"name,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Project     *int    `json:"project,omitempty"`
	Type        *string `json:"type,omitempty"`
}
