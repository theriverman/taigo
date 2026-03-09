package taigo

// ProjectTemplateDetail is a raw DTO for project template detail endpoints.
type ProjectTemplateDetail = RawResource

// ProjectTemplateDetailService is a read-only handle to project template details.
type ProjectTemplateDetailService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// Get returns project template detail by template ID.
func (s *ProjectTemplateDetailService) Get(projectTemplateID int) (*ProjectTemplateDetail, error) {
	return getRawResource(s.client, s.Endpoint, projectTemplateID)
}
