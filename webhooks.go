package taigo

import (
	"fmt"
	"strconv"

	"github.com/google/go-querystring/query"
)

// WebhookService is a handle to actions related to Webhooks
//
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks
type WebhookService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
	EndpointLogs     string
}

// ListWebhooks returns all Webhooks
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-list
func (s *WebhookService) ListWebhooks(queryParams *WebhookQueryParameters) ([]Webhook, error) {
	url := s.client.MakeURL(s.Endpoint)
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
	}
	var webhooks []Webhook

	_, err := s.client.Request.Get(url, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// CreateWebhook creates a new Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-create
func (s *WebhookService) CreateWebhook(webhook *Webhook) (*Webhook, error) {
	url := s.client.MakeURL(s.Endpoint)
	var wh Webhook

	_, err := s.client.Request.Post(url, &webhook, &wh)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

// GetWebhook returns a Webhook by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-get
func (s *WebhookService) GetWebhook(webhook *Webhook) (*Webhook, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID))
	var wh Webhook
	_, err := s.client.Request.Get(url, &wh)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

// EditWebhook sends a PATCH request to edit a Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-edit
func (s *WebhookService) EditWebhook(webhook *Webhook) (*Webhook, error) {
	var responseWebhook Webhook
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID))
	_, err := s.client.Request.Patch(url, &webhook, &responseWebhook)
	if err != nil {
		return nil, err
	}
	return &responseWebhook, nil
}

// DeleteWebhook sends a DELETE request to delete a Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-delete
func (s *WebhookService) DeleteWebhook(webhook *Webhook) error {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID))
	_, err := s.client.Request.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

// TestWebhook sends an empty POST request to test a webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-test
func (s *WebhookService) TestWebhook(webhook *Webhook) (*WebhookLog, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID))
	var whLog WebhookLog
	_, err := s.client.Request.Post(url, &webhook, &whLog)
	if err != nil {
		return nil, err
	}
	return &whLog, nil
}

// ListWebhookLogs returns all Webhook logs
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-list
func (s *WebhookService) ListWebhookLogs(queryParameters *WebhookQueryParameters) (*[]WebhookLog, error) {
	url := s.client.MakeURL(s.EndpointLogs)
	if queryParameters != nil {
		queryParameters.ProjectID = 0 // dropping projectID because not required here
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	var whLogs []WebhookLog

	_, err := s.client.Request.Get(url, &whLogs)
	if err != nil {
		return nil, err
	}
	return &whLogs, nil
}

// GetWebhookLog returns a WebhookLog by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-get
func (s *WebhookService) GetWebhookLog(webhook *Webhook) (*WebhookLog, error) {
	url := s.client.MakeURL(s.EndpointLogs, strconv.Itoa(webhook.ID))
	var whLog WebhookLog
	_, err := s.client.Request.Get(url, &whLog)
	if err != nil {
		return nil, err
	}
	return &whLog, nil
}

// ResendWebhookRequest resends the request from a Webhook Log by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-resend
func (s *WebhookService) ResendWebhookRequest(webhookLog *WebhookLog) (*WebhookLog, error) {
	url := s.client.MakeURL(s.EndpointLogs, strconv.Itoa(webhookLog.ID), "resend")
	var whLog WebhookLog
	_, err := s.client.Request.Post(url, &webhookLog, &whLog)
	if err != nil {
		return nil, err
	}
	return &whLog, nil
}
