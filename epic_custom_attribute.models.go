package taigo

import "time"

// EpicCustomAttribute -> https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attribute-detail
type EpicCustomAttribute struct {
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
