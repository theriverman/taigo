package taigo

// TaskCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-task-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type TaskCustomAttribValues struct {
	TgObjCAVDBase
	Task int `json:"task,omitempty"`
}
