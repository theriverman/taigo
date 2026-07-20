package taigo

// UserStoryCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-user-story-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type UserStoryCustomAttribValues struct {
	TgObjCAVDBase
	UserStory int `json:"user_story,omitempty"`
}

// UserStoryCustomAttributesValues -> https://docs.taiga.io/api.html#object-userstory-custom-attributes-values-detail
type UserStoryCustomAttributesValues struct {
	AttributesValues map[string]string `json:"attributes_values"`
	Version          int               `json:"version"`
	UserStory        int               `json:"user_story"`
}
