package taigo

import "time"

// TaskCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#task-custom-attributes
type TaskCustomAttribute struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	ProjectID    int         `json:"project_id"`
	Type         string      `json:"type"`
}
