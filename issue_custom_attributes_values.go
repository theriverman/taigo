package taigo

// IssueCustomAttribValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-issue-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD
type IssueCustomAttribValues struct {
	TgObjCAVDBase
	Issue int `json:"issue,omitempty"`
}

// IssueCustomAttributesValues -> https://docs.taiga.io/api.html#object-issue-custom-attributes-values-detail
type IssueCustomAttributesValues struct {
	AttributesValues map[string]string `json:"attributes_values"`
	Version          int               `json:"version"`
	Issue            int               `json:"issue"`
}
