package gotaiga

var wikiURI = "/wiki"

// WikiCreateAttachment creates a new Wiki attachment => https://taigaio.github.io/taiga-doc/dist/api.html#wiki-create-attachment
func WikiCreateAttachment(c *Client, attachment *Attachment, filePath string) (*Attachment, error) {
	url := c.APIURL + wikiURI + "/attachments"
	attachment.filePath = filePath
	attachment, err := newfileUploadRequest(c, url, attachment)
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
