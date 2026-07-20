package taigo

// EpicCustomAttributeValues -> http://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attributes-values-detail
// You must populate TgObjCAVDBase.AttributesValues with your custom struct representing the actual CAVD.
type EpicCustomAttributeValues struct {
	TgObjCAVDBase
	Epic int `json:"epic,omitempty"`
}

// EpicCustomAttributesValues -> https://docs.taiga.io/api.html#object-epic-custom-attributes-values-detail
type EpicCustomAttributesValues struct {
	AttributesValues map[string]string `json:"attributes_values"`
	Version          int               `json:"version"`
	Epic             int               `json:"epic"`
}
