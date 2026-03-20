package taigo

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

type contractCase struct {
	name            string
	invoke          func(*Client) error
	wantMethod      string
	wantPath        string
	wantQuery       map[string]string
	strictQuery     bool
	responseStatus  int
	responseBody    string
	expectEmptyBody bool
	assertBody      func(t *testing.T, body []byte)
}

func TestContractMatrixSingleRequestRoutes(t *testing.T) {
	cases := []contractCase{
		{
			name:       "applications/get_token",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/applications/42/token",
			invoke: func(c *Client) error {
				_, err := c.Application.GetToken(42)
				return err
			},
		},
		{
			name:            "memberships/resend_invitation",
			wantMethod:      http.MethodPost,
			wantPath:        "/api/v1/memberships/7/resend_invitation",
			expectEmptyBody: true,
			invoke: func(c *Client) error {
				_, err := c.MembershipInvitation.ResendInvitation(7)
				return err
			},
		},
		{
			name:         "timeline/user",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/timeline/user/3",
			wantQuery:    map[string]string{"page": "2"},
			strictQuery:  true,
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.Timeline.User(3, &TimelineQueryParams{Page: 2})
				return err
			},
		},
		{
			name:         "timeline/project",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/timeline/project/9",
			wantQuery:    map[string]string{"page": "1"},
			strictQuery:  true,
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.Timeline.Project(9, &TimelineQueryParams{Page: 1})
				return err
			},
		},
		{
			name:         "timeline/profile",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/timeline/profile/5",
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.Timeline.Profile(5, nil)
				return err
			},
		},
		{
			name:       "importers/trello/list_projects",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/importers/trello/list_projects",
			invoke: func(c *Client) error {
				_, err := c.Importer.TrelloListProjects(nil)
				return err
			},
		},
		{
			name:        "importers/github/list_users",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/importers/github/list_users",
			wantQuery:   map[string]string{"project": "2"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.Importer.GithubListUsers(&ImporterAuthURLQueryParams{Project: 2})
				return err
			},
		},
		{
			name:       "importers/github/import_project",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/importers/github/import_project",
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				if !strings.Contains(string(body), `"project":"repo-one"`) {
					t.Fatalf("expected import payload to include project, got %s", string(body))
				}
			},
			invoke: func(c *Client) error {
				_, err := c.Importer.GithubImportProject(RawResource{"project": "repo-one"})
				return err
			},
		},
		{
			name:       "importers/jira/import_project",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/importers/jira/import_project",
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				if !strings.Contains(string(body), `"project":"jira-one"`) {
					t.Fatalf("expected import payload to include project, got %s", string(body))
				}
			},
			invoke: func(c *Client) error {
				_, err := c.Importer.JiraImportProject(RawResource{"project": "jira-one"})
				return err
			},
		},
		{
			name:            "history/delete_comment_uses_query",
			wantMethod:      http.MethodPost,
			wantPath:        "/api/v1/history/task/11/delete_comment",
			wantQuery:       map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			strictQuery:     true,
			expectEmptyBody: true,
			invoke: func(c *Client) error {
				_, err := c.History.DeleteTaskComment(11, "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name:        "tasks/get_by_ref",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/tasks/by_ref",
			wantQuery:   map[string]string{"ref": "8", "project": "2"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.Task.GetByRef(8, &Project{ID: 2})
				return err
			},
		},
		{
			name:        "issues/get_by_ref",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/issues/by_ref",
			wantQuery:   map[string]string{"ref": "15", "project": "2"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.Issue.GetByRef(15, &Project{ID: 2})
				return err
			},
		},
		{
			name:        "userstories/get_by_ref",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/userstories/by_ref",
			wantQuery:   map[string]string{"ref": "21", "project": "2"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.UserStory.GetByRef(21, &Project{ID: 2})
				return err
			},
		},
		{
			name:        "wiki/get_by_slug",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/wiki/by_slug",
			wantQuery:   map[string]string{"slug": "my-page", "project": "2"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.Wiki.GetBySlug("my-page", 2)
				return err
			},
		},
		{
			name:         "wiki/render",
			wantMethod:   http.MethodPost,
			wantPath:     "/api/v1/wiki/render",
			responseBody: `{"data":"<p>ok</p>"}`,
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				if !strings.Contains(string(body), `"content":"**hello**"`) {
					t.Fatalf("expected render payload content, got %s", string(body))
				}
				if !strings.Contains(string(body), `"project_id":2`) {
					t.Fatalf("expected render payload project_id, got %s", string(body))
				}
			},
			invoke: func(c *Client) error {
				_, err := c.Wiki.Render("**hello**", 2)
				return err
			},
		},
		{
			name:            "webhooks/test",
			wantMethod:      http.MethodPost,
			wantPath:        "/api/v1/webhooks/9/test",
			expectEmptyBody: true,
			invoke: func(c *Client) error {
				_, err := c.Webhook.Test(9)
				return err
			},
		},
		{
			name:         "webhooklogs/list",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/webhooklogs",
			wantQuery:    map[string]string{"webhook": "9"},
			strictQuery:  true,
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.Webhook.Logs(&WebhookQueryParameters{ProjectID: 2, WebhookID: 9})
				return err
			},
		},
		{
			name:        "projects/get_by_slug",
			wantMethod:  http.MethodGet,
			wantPath:    "/api/v1/projects/by_slug",
			wantQuery:   map[string]string{"slug": "demo"},
			strictQuery: true,
			invoke: func(c *Client) error {
				_, err := c.Project.GetBySlug("demo")
				return err
			},
		},
		{
			name:         "tasks/list_attachments",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/tasks/attachments",
			wantQuery:    map[string]string{"object_id": "8", "project": "4"},
			strictQuery:  true,
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.Task.ListAttachments(&Task{ID: 8, Project: 4})
				return err
			},
		},
		{
			name:         "memberships/list_invitations",
			wantMethod:   http.MethodGet,
			wantPath:     "/api/v1/memberships",
			wantQuery:    map[string]string{"project": "2"},
			strictQuery:  true,
			responseBody: "[]",
			invoke: func(c *Client) error {
				_, err := c.MembershipInvitation.ListInvitations(&MembershipInvitationsQueryParams{Project: 2})
				return err
			},
		},
		{
			name:       "invitations/get_by_token",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/invitations/token-abc",
			invoke: func(c *Client) error {
				_, err := c.MembershipInvitation.GetInvitationByToken("token-abc")
				return err
			},
		},
		{
			name:       "invitations/apply_by_token",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/invitations/token-abc",
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				if !strings.Contains(string(body), `"accepted_terms":true`) {
					t.Fatalf("expected invitation apply payload, got %s", string(body))
				}
			},
			invoke: func(c *Client) error {
				_, err := c.MembershipInvitation.ApplyInvitationByToken("token-abc", RawResource{"accepted_terms": true})
				return err
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			status := tc.responseStatus
			if status == 0 {
				status = http.StatusOK
			}
			body := tc.responseBody
			if body == "" {
				body = "{}"
			}

			client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
				if req.Method != tc.wantMethod {
					t.Fatalf("method mismatch: got %s want %s", req.Method, tc.wantMethod)
				}
				if req.URL.Path != tc.wantPath {
					t.Fatalf("path mismatch: got %s want %s", req.URL.Path, tc.wantPath)
				}
				if tc.wantQuery != nil {
					queryValues := req.URL.Query()
					for key, wantValue := range tc.wantQuery {
						if gotValue := queryValues.Get(key); gotValue != wantValue {
							t.Fatalf("query mismatch for %s: got %q want %q", key, gotValue, wantValue)
						}
					}
					if tc.strictQuery && len(queryValues) != len(tc.wantQuery) {
						t.Fatalf("query key count mismatch: got %d want %d (%v)", len(queryValues), len(tc.wantQuery), queryValues)
					}
				}

				if req.Body != nil {
					payload, _ := io.ReadAll(req.Body)
					if tc.expectEmptyBody && len(bytes.TrimSpace(payload)) != 0 {
						t.Fatalf("expected empty request body, got %q", string(payload))
					}
					if tc.assertBody != nil {
						tc.assertBody(t, payload)
					}
				} else if tc.assertBody != nil {
					tc.assertBody(t, nil)
				}
				return newJSONResponse(req, status, body), nil
			})

			if err := tc.invoke(client); err != nil {
				t.Fatalf("invoke failed: %v", err)
			}
		})
	}
}

func TestContractTaskEditUsesSinglePATCHWithCallerVersion(t *testing.T) {
	step := 0
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		step++
		if req.Method != http.MethodPatch || req.URL.Path != "/api/v1/tasks/11" {
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		body, _ := io.ReadAll(req.Body)
		if !strings.Contains(string(body), `"version":4`) {
			t.Fatalf("expected patch payload version from caller, got %s", string(body))
		}
		return newJSONResponse(req, http.StatusOK, `{"id":11,"version":5}`), nil
	})

	task := &Task{ID: 11, Version: 4, Subject: "updated"}
	if _, err := client.Task.Edit(task); err != nil {
		t.Fatalf("task edit failed: %v", err)
	}
	if step != 1 {
		t.Fatalf("expected exactly 1 request, got %d", step)
	}
}

func TestContractTaskEditRequiresVersion(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
		return nil, nil
	})

	_, err := client.Task.Edit(&Task{ID: 11, Subject: "updated"})
	if err == nil || !strings.Contains(err.Error(), "version is required") {
		t.Fatalf("expected version validation error, got %v", err)
	}
}

func TestContractWikiEditUsesSinglePATCHWithCallerVersion(t *testing.T) {
	step := 0
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		step++
		if req.Method != http.MethodPatch || req.URL.Path != "/api/v1/wiki/5" {
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		body, _ := io.ReadAll(req.Body)
		if !strings.Contains(string(body), `"version":3`) {
			t.Fatalf("expected patch payload version from caller, got %s", string(body))
		}
		return newJSONResponse(req, http.StatusOK, `{"id":5,"version":4}`), nil
	})

	page := &WikiPage{ID: 5, Version: 3, Content: "# updated wiki"}
	if _, err := client.Wiki.Edit(page); err != nil {
		t.Fatalf("wiki edit failed: %v", err)
	}
	if step != 1 {
		t.Fatalf("expected exactly 1 request, got %d", step)
	}
}

func TestContractWikiEditRequiresVersion(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
		return nil, nil
	})

	_, err := client.Wiki.Edit(&WikiPage{ID: 5, Content: "# updated wiki"})
	if err == nil || !strings.Contains(err.Error(), "version is required") {
		t.Fatalf("expected version validation error, got %v", err)
	}
}

func TestContractInvitationTokenFallsBackToPOSTWhenGETNotAllowed(t *testing.T) {
	step := 0
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		step++
		switch step {
		case 1:
			if req.Method != http.MethodGet || req.URL.Path != "/api/v1/invitations/token-fallback" {
				t.Fatalf("unexpected first request: %s %s", req.Method, req.URL.Path)
			}
			return newJSONResponse(req, http.StatusMethodNotAllowed, `{"detail":"method not allowed"}`), nil
		case 2:
			if req.Method != http.MethodPost || req.URL.Path != "/api/v1/invitations/token-fallback" {
				t.Fatalf("unexpected second request: %s %s", req.Method, req.URL.Path)
			}
			var body []byte
			if req.Body != nil {
				body, _ = io.ReadAll(req.Body)
			}
			if len(bytes.TrimSpace(body)) != 0 {
				t.Fatalf("expected empty POST payload in fallback call, got %q", string(body))
			}
			return newJSONResponse(req, http.StatusOK, `{"ok":true}`), nil
		default:
			t.Fatalf("unexpected extra request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	if _, err := client.MembershipInvitation.GetInvitationByToken("token-fallback"); err != nil {
		t.Fatalf("invitation token fallback failed: %v", err)
	}
	if step != 2 {
		t.Fatalf("expected exactly 2 requests, got %d", step)
	}
}
