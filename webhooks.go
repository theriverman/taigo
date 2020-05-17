package gotaiga

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointWebhooksURI = "/webhooks"
const endpointWebhookLogs = "/webhooklogs"

// WebhookService is a handle to actions related to Webhooks
//
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks
type WebhookService struct {
	client *Client
}

// ListWebhooks returns all Webhooks
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-list
func (s *WebhookService) ListWebhooks(queryParameters *WebhookQueryParameters) ([]Webhook, error) {
	url := s.client.APIURL + endpointWebhooksURI
	if queryParameters != nil {
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	var webhooks []Webhook

	err := getRequest(s.client, &webhooks, url)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// CreateWebhook creates a new Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-create
func (s *WebhookService) CreateWebhook(webhook *Webhook) (*Webhook, error) {
	url := s.client.APIURL + endpointWebhooksURI
	var newWebhook Webhook

	err := postRequest(s.client, &newWebhook, url, webhook)
	if err != nil {
		return nil, err
	}
	return &newWebhook, nil
}

// GetWebhook returns a Webhook by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-get
func (s *WebhookService) GetWebhook(webhook *Webhook) (*Webhook, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointWebhooksURI, webhook.ID)
	var respWebhook Webhook
	err := getRequest(s.client, &respWebhook, url)
	if err != nil {
		return nil, err
	}
	return &respWebhook, nil
}

// EditWebhook sends a PATCH request to edit a Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-edit
func (s *WebhookService) EditWebhook(webhook *Webhook) (*Webhook, error) {
	var responseWebhook Webhook
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointWebhooksURI, webhook.ID)

	err := patchRequest(s.client, &responseWebhook, url, &webhook)
	if err != nil {
		return nil, err
	}
	return &responseWebhook, nil
}

// DeleteWebhook sends a DELETE request to delete a Webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-delete
func (s *WebhookService) DeleteWebhook(webhook *Webhook) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointWebhooksURI, webhook.ID)
	err := deleteRequest(s.client, url)
	if err != nil {
		return err
	}
	return nil
}

// TestWebhook sends an empty POST request to test a webhook
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooks-test
func (s *WebhookService) TestWebhook(webhook *Webhook) (*WebhookLog, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointWebhooksURI, webhook.ID)
	var responseWebhookLog WebhookLog
	err := postRequest(s.client, &responseWebhookLog, url, webhook)
	if err != nil {
		return nil, err
	}
	return &responseWebhookLog, nil
}

// ListWebhookLogs returns all Webhook logs
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-list
func (s *WebhookService) ListWebhookLogs(queryParameters *WebhookQueryParameters) (*[]WebhookLog, error) {
	url := s.client.APIURL + endpointWebhookLogs
	if queryParameters != nil {
		queryParameters.ProjectID = 0 // dropping projectID because not required here
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	var webhookLogs []WebhookLog

	err := getRequest(s.client, &webhookLogs, url)
	if err != nil {
		return nil, err
	}
	return &webhookLogs, nil
}

// GetWebhookLog returns a WebhookLog by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-get
func (s *WebhookService) GetWebhookLog(webhook *Webhook) (*WebhookLog, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointWebhookLogs, webhook.ID)
	var respWebhookLog WebhookLog
	err := getRequest(s.client, &respWebhookLog, url)
	if err != nil {
		return nil, err
	}
	return &respWebhookLog, nil
}

// ResendWebhookRequest resends the request from a Webhook Log by ID
// https://taigaio.github.io/taiga-doc/dist/api.html#webhooklogs-resend
func (s *WebhookService) ResendWebhookRequest(webhookLog *WebhookLog) (*WebhookLog, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d/resend", endpointWebhookLogs, webhookLog.ID)
	var respWebhookLog WebhookLog
	err := postRequest(s.client, &respWebhookLog, url, webhookLog)
	if err != nil {
		return nil, err
	}
	return &respWebhookLog, nil
}
