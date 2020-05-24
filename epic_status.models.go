package taigo

// EpicStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#epic-statuses
type EpicStatus struct {
	Color    string `json:"color"`
	ID       int    `json:"id"`
	IsClosed bool   `json:"is_closed"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	Project  int    `json:"project"`
	Slug     string `json:"slug"`
}
