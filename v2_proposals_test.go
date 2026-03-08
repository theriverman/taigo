package taigo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newJSONResponse(req *http.Request, statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func newUnitTestClient(t *testing.T, rt roundTripperFunc) *Client {
	t.Helper()

	client := &Client{
		BaseURL:             "http://taiga.test",
		HTTPClient:          &http.Client{Transport: rt},
		AutoRefreshDisabled: true,
	}
	if err := client.Initialise(); err != nil {
		t.Fatalf("initialise client: %v", err)
	}
	return client
}

func TestDisablePaginationHeaderSemantics(t *testing.T) {
	client := &Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}, AutoRefreshDisabled: true}
	if err := client.Initialise(); err != nil {
		t.Fatalf("initialise client: %v", err)
	}

	if got := client.Headers.Get("x-disable-pagination"); got != "True" {
		t.Fatalf("expected x-disable-pagination to be True by default, got %q", got)
	}

	client.DisablePagination(false)
	if got := client.Headers.Get("x-disable-pagination"); got != "" {
		t.Fatalf("expected x-disable-pagination to be removed when pagination is enabled, got %q", got)
	}

	client.DisablePagination(true)
	client.DisablePagination(true)
	if got := len(client.Headers.Values("x-disable-pagination")); got != 1 {
		t.Fatalf("expected exactly one x-disable-pagination header value, got %d", got)
	}
}

func TestProjectsQueryOrderBySerialises(t *testing.T) {
	params := &ProjectsQueryParameters{}
	params.TotalFansLastMonth()

	v, err := query.Values(params)
	if err != nil {
		t.Fatalf("encode query params: %v", err)
	}
	if got := v.Get("order_by"); got != "total_fans_last_month" {
		t.Fatalf("expected order_by=total_fans_last_month, got %q", got)
	}
}

func TestMilestoneQueryCanEncodeExplicitFalse(t *testing.T) {
	params := &MilestonesQueryParams{Closed: BoolPtr(false)}
	v, err := query.Values(params)
	if err != nil {
		t.Fatalf("encode query params: %v", err)
	}
	if got := v.Get("closed"); got != "false" {
		t.Fatalf("expected closed=false, got %q", got)
	}
}

func TestTaskQueryTagsEncodeAsCommaSeparatedValue(t *testing.T) {
	params := &TasksQueryParams{}
	params.SetTags("backend", "api", "v2")
	v, err := query.Values(params)
	if err != nil {
		t.Fatalf("encode query params: %v", err)
	}
	if got := v.Get("tags"); got != "backend,api,v2" {
		t.Fatalf("expected tags=backend,api,v2, got %q", got)
	}
}

func TestWebhookTestUsesTestEndpoint(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []byte

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			gotBody = body
		}
		return newJSONResponse(req, http.StatusOK, `{"id":123}`), nil
	})

	whLog, err := client.Webhook.Test(42)
	if err != nil {
		t.Fatalf("test webhook: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST method, got %s", gotMethod)
	}
	if gotPath != "/api/v1/webhooks/42/test" {
		t.Fatalf("expected test endpoint path, got %s", gotPath)
	}
	if len(bytes.TrimSpace(gotBody)) != 0 {
		t.Fatalf("expected empty request body for webhook test call, got %q", string(gotBody))
	}
	if whLog.ID != 123 {
		t.Fatalf("expected decoded webhook log id 123, got %d", whLog.ID)
	}
}

func TestUsersWatchedAndLikedDecodeArrayPayloads(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/api/v1/users/7/watched":
			return newJSONResponse(req, http.StatusOK, `[{"id":11,"type":"task","project":2}]`), nil
		case "/api/v1/users/7/liked":
			if got := req.URL.Query().Get("type"); got != "task" {
				t.Fatalf("expected liked query type=task, got %q", got)
			}
			if got := req.URL.Query().Get("q"); got != "bug" {
				t.Fatalf("expected liked query q=bug, got %q", got)
			}
			return newJSONResponse(req, http.StatusOK, `[{"id":13,"type":"issue","project":4}]`), nil
		default:
			t.Fatalf("unexpected path hit: %s", req.URL.Path)
			return nil, nil
		}
	})

	watched, err := client.User.GetWatchedContent(7, nil)
	if err != nil {
		t.Fatalf("get watched: %v", err)
	}
	if len(watched) != 1 || watched[0].ID != 11 {
		t.Fatalf("expected 1 watched item with id 11, got %+v", watched)
	}

	liked, err := client.User.GetLikedContent(7, &UsersHighlightedQueryParams{Type: "task", Q: "bug"})
	if err != nil {
		t.Fatalf("get liked: %v", err)
	}
	if len(liked) != 1 || liked[0].ID != 13 {
		t.Fatalf("expected 1 liked item with id 13, got %+v", liked)
	}
}

func TestListWebhookLogsIgnoresProjectFilter(t *testing.T) {
	var gotWebhook string
	var gotProject string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotWebhook = req.URL.Query().Get("webhook")
		gotProject = req.URL.Query().Get("project")
		return newJSONResponse(req, http.StatusOK, `[{"id":99}]`), nil
	})

	logs, err := client.Webhook.ListWebhookLogs(&WebhookQueryParameters{ProjectID: 12, WebhookID: 34})
	if err != nil {
		t.Fatalf("list webhook logs: %v", err)
	}
	if gotWebhook != "34" {
		t.Fatalf("expected webhook query to be retained, got %q", gotWebhook)
	}
	if gotProject != "" {
		t.Fatalf("expected project query to be removed for webhook logs, got %q", gotProject)
	}
	if len(*logs) != 1 || (*logs)[0].ID != 99 {
		t.Fatalf("unexpected webhook logs payload: %+v", *logs)
	}
}

func TestCustomAttributeValueJSONKeys(t *testing.T) {
	issuePayload, err := json.Marshal(IssueCustomAttribValues{Issue: 11})
	if err != nil {
		t.Fatalf("marshal issue custom attributes values: %v", err)
	}
	if !strings.Contains(string(issuePayload), "\"issue\"") || strings.Contains(string(issuePayload), "\"epic\"") {
		t.Fatalf("unexpected issue custom attributes values payload: %s", string(issuePayload))
	}

	taskPayload, err := json.Marshal(TaskCustomAttribValues{Task: 22})
	if err != nil {
		t.Fatalf("marshal task custom attributes values: %v", err)
	}
	if !strings.Contains(string(taskPayload), "\"task\"") || strings.Contains(string(taskPayload), "\"epic\"") {
		t.Fatalf("unexpected task custom attributes values payload: %s", string(taskPayload))
	}

	userStoryPayload, err := json.Marshal(UserStoryCustomAttribValues{UserStory: 33})
	if err != nil {
		t.Fatalf("marshal user story custom attributes values: %v", err)
	}
	if !strings.Contains(string(userStoryPayload), "\"user_story\"") || strings.Contains(string(userStoryPayload), "\"epic\"") {
		t.Fatalf("unexpected user story custom attributes values payload: %s", string(userStoryPayload))
	}
}

func TestRequestServiceReturnsAPIErrorForNon2xx(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		return newJSONResponse(req, http.StatusBadRequest, `{\"detail\":\"bad request\"}`), nil
	})

	_, err := client.Request.Get(client.MakeURL("projects"), &Project{})
	if err == nil {
		t.Fatalf("expected API error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Body, "bad request") {
		t.Fatalf("expected API error body to include backend message, got %q", apiErr.Body)
	}
}
