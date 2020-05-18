package taigo

import "time"

// WikiPage -> https://taigaio.github.io/taiga-doc/dist/api.html#object-wiki-detail
type WikiPage struct {
	TaigaBaseObject
	Content          string    `json:"content"`
	CreatedDate      time.Time `json:"created_date"`
	Editions         int       `json:"editions"`
	HTML             string    `json:"html"`
	ID               int       `json:"id"`
	IsWatcher        bool      `json:"is_watcher"`
	LastModifier     int       `json:"last_modifier"`
	ModifiedDate     time.Time `json:"modified_date"`
	Owner            int       `json:"owner"`
	Project          int       `json:"project"`
	ProjectExtraInfo struct {
		ID           int    `json:"id"`
		LogoSmallURL string `json:"logo_small_url"`
		Name         string `json:"name"`
		Slug         string `json:"slug"`
	} `json:"project_extra_info"`
	Slug          string `json:"slug"`
	TotalWatchers int    `json:"total_watchers"`
	Version       int    `json:"version"`
}

// GetID returns the ID
func (tgObj *WikiPage) GetID() int {
	return tgObj.ID
}

// GetVersion return the version
func (tgObj *WikiPage) GetVersion() int {
	return tgObj.Version
}

// GetProject returns the project ID
func (tgObj *WikiPage) GetProject() int {
	return tgObj.Project
}
