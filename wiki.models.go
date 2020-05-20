package taigo

import "time"

// WikiPage -> https://taigaio.github.io/taiga-doc/dist/api.html#object-wiki-detail
type WikiPage struct {
	TaigaBaseObject
	Content          string           `json:"content"`
	CreatedDate      time.Time        `json:"created_date"`
	Editions         int              `json:"editions"`
	HTML             string           `json:"html"`
	ID               int              `json:"id"`
	IsWatcher        bool             `json:"is_watcher"`
	LastModifier     int              `json:"last_modifier"`
	ModifiedDate     time.Time        `json:"modified_date"`
	Owner            int              `json:"owner"`
	Project          int              `json:"project"`
	ProjectExtraInfo ProjectExtraInfo `json:"project_extra_info"`
	Slug             string           `json:"slug"`
	TotalWatchers    int              `json:"total_watchers"`
	Version          int              `json:"version"`
}

// GetID returns the ID
func (w *WikiPage) GetID() int {
	return w.ID
}

// GetVersion return the version
func (w *WikiPage) GetVersion() int {
	return w.Version
}

// GetProject returns the project ID
func (w *WikiPage) GetProject() int {
	return w.Project
}
