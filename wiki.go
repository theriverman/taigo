package taigo

// WikiService is a handle to actions related to Tasks
//
// https://taigaio.github.io/taiga-doc/dist/api.html#tasks
type WikiService struct {
	client   *Client
	Endpoint string
}

// CreateAttachment creates a new Wiki attachment -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-create-attachment
func (s *WikiService) CreateAttachment(attachment *Attachment, wikiPage *WikiPage) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, wikiPage)
}
