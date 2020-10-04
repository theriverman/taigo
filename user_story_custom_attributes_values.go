package taigo

// UserStoryCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type UserStoryCustomAttribValues struct {
	TgObjCAVDBase
	Epic int `json:"epic,omitempty"`
}
