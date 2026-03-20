package taigo

import (
	"errors"
	"net/http"
	"testing"
)

func TestNegativeMatrixOfflineAPIErrors(t *testing.T) {
	type matrixCase struct {
		name   string
		invoke func(*Client) error
	}

	cases := []matrixCase{
		{
			name: "request/get",
			invoke: func(c *Client) error {
				_, err := c.Request.Get(c.MakeURL("projects"), &Project{})
				return err
			},
		},
		{
			name: "epics/get",
			invoke: func(c *Client) error {
				_, err := c.Epic.Get(99)
				return err
			},
		},
		{
			name: "userstories/get",
			invoke: func(c *Client) error {
				_, err := c.UserStory.Get(99)
				return err
			},
		},
		{
			name: "tasks/get",
			invoke: func(c *Client) error {
				_, err := c.Task.Get(99)
				return err
			},
		},
		{
			name: "issues/get",
			invoke: func(c *Client) error {
				_, err := c.Issue.Get(99)
				return err
			},
		},
		{
			name: "wiki/get",
			invoke: func(c *Client) error {
				_, err := c.Wiki.Get(99)
				return err
			},
		},
		{
			name: "webhooks/get",
			invoke: func(c *Client) error {
				_, err := c.Webhook.Get(99)
				return err
			},
		},
		{
			name: "applications/get_token",
			invoke: func(c *Client) error {
				_, err := c.Application.GetToken(99)
				return err
			},
		},
		{
			name: "history/task",
			invoke: func(c *Client) error {
				_, err := c.History.Task(99)
				return err
			},
		},
		{
			name: "importers/github_list_users",
			invoke: func(c *Client) error {
				_, err := c.Importer.GithubListUsers(nil)
				return err
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
				return newJSONResponse(req, http.StatusBadRequest, `{"detail":"bad request in negative matrix"}`), nil
			})

			err := tc.invoke(client)
			if err == nil {
				t.Fatalf("expected an API error")
			}

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *APIError, got %T (%v)", err, err)
			}
			if apiErr.StatusCode != http.StatusBadRequest {
				t.Fatalf("status mismatch: got %d want %d", apiErr.StatusCode, http.StatusBadRequest)
			}
			if apiErr.Body == "" {
				t.Fatalf("expected API error body to be populated")
			}
		})
	}
}

func TestNegativeMatrixValidationGuards(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		return newJSONResponse(req, http.StatusOK, `{}`), nil
	})

	if _, err := client.Task.Create(&Task{Project: 1}); err == nil {
		t.Fatalf("expected task create validation error for missing subject")
	}
	if _, err := client.UserStory.Create(&UserStory{Project: 1}); err == nil {
		t.Fatalf("expected user story create validation error for missing subject")
	}
	if _, err := client.Wiki.Create(&WikiPage{Project: 1, Slug: "", Content: "x"}); err == nil {
		t.Fatalf("expected wiki create validation error for missing slug")
	}
}
