package main

import (
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestWebhooks(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	// List webhooks
	webhooks, err := Client.Webhook.ListWebhooks(&taiga.WebhookQueryParameters{ProjectID: testProjID})
	if err != nil {
		t.Error(err)
	}
	t.Logf("Total webhooks: %d", len(webhooks))

	// Create webhook
	webhook, err := Client.Webhook.CreateWebhook(&taiga.Webhook{
		Name:    "Taigo Integration Webhook",
		Project: testProjID,
		URL:     "https://example.com/taigo-integration-webhook",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get webhook
	webhookByID, err := Client.Webhook.GetWebhook(&taiga.Webhook{ID: webhook.ID})
	if err != nil {
		t.Error(err)
	}
	if webhookByID.ID != webhook.ID {
		t.Errorf("got %d, want %d", webhookByID.ID, webhook.ID)
	}

	// Edit webhook
	updatedName := "Taigo Integration Webhook Updated"
	webhook.Name = updatedName
	editedWebhook, err := Client.Webhook.EditWebhook(webhook)
	if err != nil {
		t.Error(err)
	}
	if editedWebhook.Name != updatedName {
		t.Errorf("got %q, want %q", editedWebhook.Name, updatedName)
	}

	// Test webhook endpoint
	testLog, err := Client.Webhook.TestWebhook(webhook)
	if err != nil {
		t.Error(err)
	}
	if testLog.ID == 0 {
		t.Errorf("expected webhook test to produce a webhook log")
	}

	// List webhook logs
	logs, err := Client.Webhook.ListWebhookLogs(&taiga.WebhookQueryParameters{WebhookID: webhook.ID})
	if err != nil {
		t.Error(err)
	}
	if len(*logs) == 0 {
		t.Errorf("expected at least one webhook log")
	}

	// Get and resend first webhook log, when available
	if len(*logs) > 0 {
		logItem := (*logs)[0]
		logByID, err := Client.Webhook.GetWebhookLog(logItem.ID)
		if err != nil {
			t.Error(err)
		}
		if logByID.ID != logItem.ID {
			t.Errorf("got %d, want %d", logByID.ID, logItem.ID)
		}

		if _, err := Client.Webhook.ResendWebhookRequest(&taiga.WebhookLog{ID: logItem.ID}); err != nil {
			t.Error(err)
		}
	}

	// Delete webhook
	if err := Client.Webhook.DeleteWebhook(webhook); err != nil {
		t.Error(err)
	}
}
