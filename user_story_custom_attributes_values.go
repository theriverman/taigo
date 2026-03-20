package taigo

// UserStoryCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-user-story-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type UserStoryCustomAttribValues struct {
	TgObjCAVDBase
	UserStory int `json:"user_story,omitempty"`
}
