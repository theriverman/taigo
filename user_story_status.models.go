package taigo

// UserStoryStatus -> https://taigaio.github.io/taiga-doc/dist/api.html#user-story-statuses
type UserStoryStatus struct {
	Color      string `json:"color"`
	ID         int    `json:"id"`
	IsArchived bool   `json:"is_archived"`
	IsClosed   bool   `json:"is_closed"`
	Name       string `json:"name"`
	Order      int    `json:"order"`
	ProjectID  int    `json:"project_id"`
	Slug       string `json:"slug"`
	WipLimit   int    `json:"wip_limit"`
}
