package taigo

// IssueCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-issue-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type IssueCustomAttribValues struct {
	TgObjCAVDBase
	Issue int `json:"issue,omitempty"`
}
