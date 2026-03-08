package taigo

import "time"

// IssueCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#issue-custom-attributes-list
type IssueCustomAttribute struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	Project      int         `json:"project"`
	Type         string      `json:"type"`
}
