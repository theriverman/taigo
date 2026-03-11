package taigo

import (
	"strconv"
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

// List returns all webhooks for the current query scope.
func (s *WebhookService) List(queryParams *WebhookQueryParameters) ([]Webhook, error) {
	return s.ListWebhooks(queryParams)
}

// Create creates a webhook.
func (s *WebhookService) Create(webhook *Webhook) (*Webhook, error) {
	return s.CreateWebhook(webhook)
}

// Get returns a webhook by ID.
func (s *WebhookService) Get(webhookID int) (*Webhook, error) {
	return s.GetWebhook(&Webhook{ID: webhookID})
}

// Edit edits a webhook.
func (s *WebhookService) Edit(webhook *Webhook) (*Webhook, error) {
	return s.EditWebhook(webhook)
}

// Delete deletes a webhook by ID.
func (s *WebhookService) Delete(webhookID int) error {
	return s.DeleteWebhook(&Webhook{ID: webhookID})
}

// Test triggers the webhook test endpoint for a webhook ID.
func (s *WebhookService) Test(webhookID int) (*WebhookLog, error) {
	return s.TestWebhook(&Webhook{ID: webhookID})
}

// Logs returns webhook logs.
func (s *WebhookService) Logs(queryParameters *WebhookQueryParameters) (*[]WebhookLog, error) {
	return s.ListWebhookLogs(queryParameters)
}

// Log returns a webhook log by ID.
func (s *WebhookService) Log(webhookLogID int) (*WebhookLog, error) {
	return s.GetWebhookLog(webhookLogID)
}

// Resend resends a webhook request from a log ID.
func (s *WebhookService) Resend(webhookLogID int) (*WebhookLog, error) {
	return s.ResendWebhookRequest(&WebhookLog{ID: webhookLogID})
}

// ListWebhooks returns all Webhooks
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-list
func (s *WebhookService) ListWebhooks(queryParams *WebhookQueryParameters) ([]Webhook, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var webhooks []Webhook

	_, err = s.client.Request.Get(url, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// CreateWebhook creates a new Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-create
func (s *WebhookService) CreateWebhook(webhook *Webhook) (*Webhook, error) {
	if err := requireNonNil("webhook", webhook); err != nil {
		return nil, err
	}
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
	if err := requireNonNil("webhook", webhook); err != nil {
		return nil, err
	}
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
	if err := requireNonNil("webhook", webhook); err != nil {
		return nil, err
	}
	var responseWebhook Webhook
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID))
	_, err := s.client.Request.Patch(url, &webhook, &responseWebhook)
	if err != nil {
		return nil, err
	}
	return &responseWebhook, nil
}

// Update is an alias for EditWebhook.
func (s *WebhookService) Update(webhook *Webhook) (*Webhook, error) {
	return s.EditWebhook(webhook)
}

// DeleteWebhook sends a DELETE request to delete a Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-delete
func (s *WebhookService) DeleteWebhook(webhook *Webhook) error {
	if err := requireNonNil("webhook", webhook); err != nil {
		return err
	}
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
	if err := requireNonNil("webhook", webhook); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(webhook.ID), "test")
	var whLog WebhookLog
	_, err := s.client.Request.Post(url, nil, &whLog)
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
		// webhooklogs endpoint supports filtering by `webhook`, not by `project`.
		qp := *queryParameters
		qp.ProjectID = 0
		var err error
		url, err = appendQueryParams(url, &qp)
		if err != nil {
			return nil, err
		}
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
func (s *WebhookService) GetWebhookLog(webhookLogID int) (*WebhookLog, error) {
	url := s.client.MakeURL(s.EndpointLogs, strconv.Itoa(webhookLogID))
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
	if err := requireNonNil("webhookLog", webhookLog); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.EndpointLogs, strconv.Itoa(webhookLog.ID), "resend")
	var whLog WebhookLog
	_, err := s.client.Request.Post(url, nil, &whLog)
	if err != nil {
		return nil, err
	}
	return &whLog, nil
}
