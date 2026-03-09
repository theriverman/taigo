package taigo

// ExportImportService is a handle to export/import actions.
type ExportImportService struct {
	client           *Client
	defaultProjectID int
	ExporterEndpoint string
	ImporterEndpoint string
}

// ExportAsync -> https://docs.taiga.io/api.html#export-export-project
func (s *ExportImportService) ExportAsync(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.ExporterEndpoint, "async")
}

// ExportStatus -> https://docs.taiga.io/api.html#export-check-export-status
func (s *ExportImportService) ExportStatus(exportID string) (*RawResource, error) {
	return getRawResourceAtPath(s.client, s.ExporterEndpoint, exportID)
}

// ImportLoadDump -> https://docs.taiga.io/api.html#import-import-project
func (s *ExportImportService) ImportLoadDump(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.ImporterEndpoint, "load_dump")
}
