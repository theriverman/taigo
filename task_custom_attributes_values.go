package taigo

// TaskCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-task-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type TaskCustomAttribValues struct {
	TgObjCAVDBase
	Task int `json:"task,omitempty"`
}

// TaskCustomAttributesValues -> https://docs.taiga.io/api.html#object-task-custom-attributes-values-detail
type TaskCustomAttributesValues struct {
	AttributesValues map[string]string `json:"attributes_values"`
	Version          int               `json:"version"`
	Task             int               `json:"task"`
}
