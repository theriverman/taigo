package taigo

// TaskCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type TaskCustomAttribValues struct {
	TgObjCAVDBase
	Epic int `json:"epic,omitempty"`
}
