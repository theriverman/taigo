package taigo

import "time"

// IssueCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#issue-custom-attributes-list
type IssueCustomAttribute struct {
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

// IssueCustomAttributeCreateRequest represents payload for creating issue custom attributes.
type IssueCustomAttributeCreateRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project"`
	Type        string `json:"type"`
}

// IssueCustomAttributeEditRequest represents sparse non-destructive updates.
type IssueCustomAttributeEditRequest struct {
	Description string `json:"description,omitempty"`
	Extra       any    `json:"extra,omitempty"`
	Name        string `json:"name,omitempty"`
	Order       int    `json:"order,omitempty"`
	Project     int    `json:"project,omitempty"`
	Type        string `json:"type,omitempty"`
}

// IssueCustomAttributePatch represents explicit PATCH payload for issue custom attributes.
type IssueCustomAttributePatch struct {
	Description *string `json:"description,omitempty"`
	Extra       any     `json:"extra,omitempty"`
	Name        *string `json:"name,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Project     *int    `json:"project,omitempty"`
	Type        *string `json:"type,omitempty"`
}
