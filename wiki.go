package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// WikiService is a handle to actions related to Wiki pages
//
// https://taigaio.github.io/taiga-doc/dist/api.html#wiki
type WikiService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-list
func (s *WikiService) List(queryParams *WikiQueryParams) ([]WikiPage, error) {
	url := s.client.MakeURL(s.Endpoint)
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
	}
	var wikiPages []WikiPage
	_, err := s.client.Request.Get(url, &wikiPages)
	if err != nil {
		return nil, err
	}
	return wikiPages, nil
}

// Create -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-create
func (s *WikiService) Create(wikiPage *WikiPage) (*WikiPage, error) {
	url := s.client.MakeURL(s.Endpoint)
	var page WikiPage

	if isEmpty(wikiPage.Project) || isEmpty(wikiPage.Slug) || isEmpty(wikiPage.Content) {
		return nil, errors.New("a mandatory field(project, slug, content) is missing. See API documentation")
	}

	_, err := s.client.Request.Post(url, &wikiPage, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-get
func (s *WikiService) Get(wikiPageID int) (*WikiPage, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(wikiPageID))
	var page WikiPage
	_, err := s.client.Request.Get(url, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetBySlug -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-by-slug
func (s *WikiService) GetBySlug(slug string, projectID int) (*WikiPage, error) {
	paramValues, _ := query.Values(WikiQueryParams{Slug: slug, Project: projectID})
	url := fmt.Sprintf("%s?%s", s.client.MakeURL(s.Endpoint, "by_slug"), paramValues.Encode())
	var page WikiPage
	_, err := s.client.Request.Get(url, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Edit sends a PATCH request to edit a Wiki page -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-edit
func (s *WikiService) Edit(wikiPage *WikiPage) (*WikiPage, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(wikiPage.ID))
	var responseWikiPage WikiPage

	if wikiPage.ID == 0 {
		return nil, errors.New("passed WikiPage does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remotePage, err := s.Get(wikiPage.ID)
	if err != nil {
		return nil, err
	}
	wikiPage.Version = remotePage.Version
	_, err = s.client.Request.Patch(url, &wikiPage, &responseWikiPage)
	if err != nil {
		return nil, err
	}
	return &responseWikiPage, nil
}

// Update is an alias for Edit.
func (s *WikiService) Update(wikiPage *WikiPage) (*WikiPage, error) {
	return s.Edit(wikiPage)
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-delete
func (s *WikiService) Delete(wikiPageID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(wikiPageID))
	return s.client.Request.Delete(url)
}

// Render -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-render
func (s *WikiService) Render(content string, projectID int) (string, error) {
	url := s.client.MakeURL(s.Endpoint, "render")
	payload := WikiRenderPayload{
		Content:   content,
		ProjectID: projectID,
	}
	var renderResp WikiRenderResponse
	_, err := s.client.Request.Post(url, &payload, &renderResp)
	if err != nil {
		return "", err
	}
	return renderResp.Data, nil
}

// CreateAttachment creates a new Wiki attachment -> https://taigaio.github.io/taiga-doc/dist/api.html#wiki-create-attachment
func (s *WikiService) CreateAttachment(attachment *Attachment, wikiPage *WikiPage) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, wikiPage)
}
