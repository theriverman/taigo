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
