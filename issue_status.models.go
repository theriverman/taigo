package taigo

// IssueStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#issue-statuses
type IssueStatus struct {
	Color     string `json:"color"`
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	ProjectID int    `json:"project_id"`
	Slug      string `json:"slug"`
}
