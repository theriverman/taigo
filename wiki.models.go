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

// WikiQueryParams holds fields to be used as URL query parameters to filter wiki pages.
type WikiQueryParams struct {
	Project int    `url:"project,omitempty"`
	Slug    string `url:"slug,omitempty"`
}

// WikiRenderPayload is the request payload for wiki markdown rendering.
type WikiRenderPayload struct {
	Content   string `json:"content"`
	ProjectID int    `json:"project_id"`
}

// WikiRenderResponse is returned by POST /wiki/render.
type WikiRenderResponse struct {
	Data string `json:"data"`
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
