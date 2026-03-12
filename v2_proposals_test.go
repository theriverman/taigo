package taigo

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

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

	if got := client.GetHeader("x-disable-pagination"); got != "True" {
		t.Fatalf("expected x-disable-pagination to be True by default, got %q", got)
	}

	client.DisablePagination(false)
	if got := client.GetHeader("x-disable-pagination"); got != "" {
		t.Fatalf("expected x-disable-pagination to be removed when pagination is enabled, got %q", got)
	}

	client.DisablePagination(true)
	client.DisablePagination(true)
	if got := len(client.GetHeaderValues("x-disable-pagination")); got != 1 {
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

func TestRefreshAuthTokenSelfUpdateRefreshesAuthorizationHeader(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/api/v1/auth/refresh" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return newJSONResponse(req, http.StatusOK, `{"auth_token":"new-token","refresh":"new-refresh"}`), nil
	})

	client.setAuthTokens(TokenBearer, "old-token", "old-refresh")
	if _, err := client.Auth.RefreshAuthToken(true); err != nil {
		t.Fatalf("refresh token failed: %v", err)
	}
	if got := client.GetAuthorizationHeader(); got != "Bearer new-token" {
		t.Fatalf("expected updated authorization header, got %q", got)
	}
}

func TestAppendQueryParamsReturnsEncodingError(t *testing.T) {
	_, err := appendQueryParams("http://taiga.test/api/v1/users", []string{"not-a-struct"})
	if err == nil {
		t.Fatalf("expected query encoding error, got nil")
	}
}

func TestIssueCreatePayloadOmitsReadOnlyDates(t *testing.T) {
	var gotBody string
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		gotBody = string(body)
		return newJSONResponse(req, http.StatusOK, `{"id":1,"project":2,"subject":"issue"}`), nil
	})

	_, err := client.Issue.Create(&Issue{
		Project:     2,
		Subject:     "issue",
		CreatedDate: mustParseTime(t, "2025-01-01T00:00:00Z"),
	})
	if err != nil {
		t.Fatalf("issue create failed: %v", err)
	}
	if strings.Contains(gotBody, "created_date") || strings.Contains(gotBody, "modified_date") {
		t.Fatalf("unexpected read-only date fields in issue create payload: %s", gotBody)
	}
}

func TestAttachmentUploadUsesFromCommentValue(t *testing.T) {
	tmp, err := os.CreateTemp("", "taigo-upload-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString("payload"); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("parse media type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("expected multipart/form-data, got %s", mediaType)
		}
		reader := multipart.NewReader(req.Body, params["boundary"])
		fields := map[string]string{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("read multipart part: %v", err)
			}
			value, _ := io.ReadAll(part)
			if part.FileName() == "" {
				fields[part.FormName()] = string(value)
			}
		}
		if fields["from_comment"] != "true" {
			t.Fatalf("expected from_comment=true, got %q", fields["from_comment"])
		}
		return newJSONResponse(req, http.StatusOK, `{"id":1}`), nil
	})
	client.setAuthTokens(TokenBearer, "token", "refresh")

	attachment := &Attachment{Name: "file", FromComment: true}
	attachment.SetFilePath(tmp.Name())
	_, err = newfileUploadRequest(client, client.MakeURL("tasks", "attachments"), attachment, &Task{ID: 11, Project: 22})
	if err != nil {
		t.Fatalf("attachment upload failed: %v", err)
	}
}

func TestAttachmentUploadReturnsDecodeErrorOnInvalidJSON(t *testing.T) {
	tmp, err := os.CreateTemp("", "taigo-upload-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if err := os.WriteFile(tmp.Name(), []byte("payload"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		return newJSONResponse(req, http.StatusOK, `{"id":`), nil
	})
	attachment := &Attachment{Name: "file"}
	attachment.SetFilePath(tmp.Name())
	_, err = newfileUploadRequest(client, client.MakeURL("tasks", "attachments"), attachment, &Task{ID: 11, Project: 22})
	if err == nil || !strings.Contains(err.Error(), "could not decode attachment response") {
		t.Fatalf("expected attachment decode error, got %v", err)
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parse time %q: %v", value, err)
	}
	return parsed
}

func TestApplicationGetTokenUsesApplicationIDPath(t *testing.T) {
	var gotPath, gotMethod string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		return newJSONResponse(req, http.StatusOK, `{"token":"abc"}`), nil
	})

	token, err := client.Application.GetToken(42)
	if err != nil {
		t.Fatalf("get application token: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("expected GET method, got %s", gotMethod)
	}
	if gotPath != "/api/v1/applications/42/token" {
		t.Fatalf("expected path /api/v1/applications/42/token, got %s", gotPath)
	}
	if (*token)["token"] != "abc" {
		t.Fatalf("expected token payload to decode, got %+v", token)
	}
}

func TestMembershipResendInvitationUsesMembershipActionAndEmptyBody(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []byte

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			gotBody = body
		}
		return newJSONResponse(req, http.StatusOK, `{"id":7}`), nil
	})

	_, err := client.MembershipInvitation.ResendInvitation(7)
	if err != nil {
		t.Fatalf("resend invitation: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST method, got %s", gotMethod)
	}
	if gotPath != "/api/v1/memberships/7/resend_invitation" {
		t.Fatalf("expected memberships resend_invitation path, got %s", gotPath)
	}
	if len(bytes.TrimSpace(gotBody)) != 0 {
		t.Fatalf("expected empty request body, got %q", string(gotBody))
	}
}

func TestTimelineEndpointsUseIDsInPath(t *testing.T) {
	requests := make([]string, 0, 3)

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.URL.RequestURI())
		return newJSONResponse(req, http.StatusOK, `[]`), nil
	})

	if _, err := client.Timeline.User(3, &TimelineQueryParams{Page: 2}); err != nil {
		t.Fatalf("timeline user: %v", err)
	}
	if _, err := client.Timeline.Project(9, nil); err != nil {
		t.Fatalf("timeline project: %v", err)
	}
	if _, err := client.Timeline.Profile(5, nil); err != nil {
		t.Fatalf("timeline profile: %v", err)
	}

	want := []string{
		"/api/v1/timeline/user/3?page=2",
		"/api/v1/timeline/project/9",
		"/api/v1/timeline/profile/5",
	}
	if len(requests) != len(want) {
		t.Fatalf("expected %d requests, got %d (%v)", len(want), len(requests), requests)
	}
	for i := range want {
		if requests[i] != want[i] {
			t.Fatalf("request %d mismatch: want %s, got %s", i, want[i], requests[i])
		}
	}
}

func TestTimelineProjectCanUseMappedDefaultProjectID(t *testing.T) {
	var gotPath string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		return newJSONResponse(req, http.StatusOK, `[]`), nil
	})
	client.Project.ConfigureMappedServices(55)

	if _, err := client.Project.Timeline.Project(0, nil); err != nil {
		t.Fatalf("timeline project with mapped project id: %v", err)
	}
	if gotPath != "/api/v1/timeline/project/55" {
		t.Fatalf("expected mapped timeline project path, got %s", gotPath)
	}
}

func TestImporterEndpointsUseCurrentRoutes(t *testing.T) {
	requests := make([]string, 0, 4)
	methods := make([]string, 0, 4)

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.URL.Path)
		methods = append(methods, req.Method)
		return newJSONResponse(req, http.StatusOK, `{"ok":true}`), nil
	})

	if _, err := client.Importer.TrelloListProjects(nil); err != nil {
		t.Fatalf("trello list projects: %v", err)
	}
	if _, err := client.Importer.GithubListUsers(nil); err != nil {
		t.Fatalf("github list users: %v", err)
	}
	if _, err := client.Importer.GithubImportProject(RawResource{"project": "x"}); err != nil {
		t.Fatalf("github import project: %v", err)
	}
	if _, err := client.Importer.JiraImportProject(RawResource{"project": "y"}); err != nil {
		t.Fatalf("jira import project: %v", err)
	}

	wantPaths := []string{
		"/api/v1/importers/trello/list_projects",
		"/api/v1/importers/github/list_users",
		"/api/v1/importers/github/import_project",
		"/api/v1/importers/jira/import_project",
	}
	wantMethods := []string{
		http.MethodGet,
		http.MethodGet,
		http.MethodPost,
		http.MethodPost,
	}
	for i := range wantPaths {
		if requests[i] != wantPaths[i] {
			t.Fatalf("request %d path mismatch: want %s, got %s", i, wantPaths[i], requests[i])
		}
		if methods[i] != wantMethods[i] {
			t.Fatalf("request %d method mismatch: want %s, got %s", i, wantMethods[i], methods[i])
		}
	}
}

func TestInitialiseSetsDefaultHTTPClientWhenUnset(t *testing.T) {
	client := &Client{
		BaseURL:             "http://example.com",
		AutoRefreshDisabled: true,
	}
	if err := client.Initialise(); err != nil {
		t.Fatalf("initialise client: %v", err)
	}
	if client.HTTPClient == nil {
		t.Fatalf("expected HTTPClient to be initialized")
	}
	if client.HTTPClient.Timeout != DefaultHTTPTimeout {
		t.Fatalf("expected default HTTP timeout %s, got %s", DefaultHTTPTimeout, client.HTTPClient.Timeout)
	}
}

func TestCloseDisablesAutoRefresh(t *testing.T) {
	client := &Client{}
	client.Close()
	if !client.AutoRefreshDisabled {
		t.Fatalf("expected AutoRefreshDisabled to be true after Close")
	}
}

func TestTaskPatchCanSendZeroAndFalseValues(t *testing.T) {
	var gotPath string
	var gotMethod string
	var gotBody string
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		body, _ := io.ReadAll(req.Body)
		gotBody = string(body)
		return newJSONResponse(req, http.StatusOK, `{"id":11,"version":2}`), nil
	})

	isBlocked := false
	milestone := 0
	subject := ""
	_, err := client.Task.Patch(11, &TaskPatch{
		IsBlocked: &isBlocked,
		Milestone: &milestone,
		Subject:   &subject,
		Version:   1,
	})
	if err != nil {
		t.Fatalf("task patch failed: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/v1/tasks/11" {
		t.Fatalf("unexpected request %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"is_blocked":false`) {
		t.Fatalf("expected explicit false in payload, got %s", gotBody)
	}
	if !strings.Contains(gotBody, `"milestone":0`) {
		t.Fatalf("expected explicit zero milestone in payload, got %s", gotBody)
	}
	if !strings.Contains(gotBody, `"subject":""`) {
		t.Fatalf("expected explicit empty subject in payload, got %s", gotBody)
	}
}

func TestMappedTaskListInjectsDefaultProjectWithoutMutatingQuery(t *testing.T) {
	var gotProject string
	var gotStatus string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotProject = req.URL.Query().Get("project")
		gotStatus = req.URL.Query().Get("status")
		return newJSONResponse(req, http.StatusOK, `[]`), nil
	})
	client.Project.ConfigureMappedServices(55)

	q := &TasksQueryParams{Status: 3}
	if _, err := client.Project.Task.List(q); err != nil {
		t.Fatalf("list tasks with mapped service: %v", err)
	}
	if gotProject != "55" {
		t.Fatalf("expected injected project=55, got %q", gotProject)
	}
	if gotStatus != "3" {
		t.Fatalf("expected status query to be preserved, got %q", gotStatus)
	}
	if q.Project != 0 {
		t.Fatalf("expected caller query to remain unchanged, got project=%d", q.Project)
	}
}

func TestMappedTaskListRespectsExplicitProjectInQuery(t *testing.T) {
	var gotProject string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotProject = req.URL.Query().Get("project")
		return newJSONResponse(req, http.StatusOK, `[]`), nil
	})
	client.Project.ConfigureMappedServices(55)

	q := &TasksQueryParams{Project: 77}
	if _, err := client.Project.Task.List(q); err != nil {
		t.Fatalf("list tasks with explicit project: %v", err)
	}
	if gotProject != "77" {
		t.Fatalf("expected explicit project=77 to win, got %q", gotProject)
	}
}

func TestRequestServiceCapturesPaginationHeaders(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		resp := newJSONResponse(req, http.StatusOK, `[]`)
		resp.Header.Set("X-Paginated", "true")
		resp.Header.Set("X-Paginated-By", "50")
		resp.Header.Set("X-Paginated-Count", "150")
		resp.Header.Set("X-Paginated-Current", "2")
		resp.Header.Set("X-Pagination-Next", "http://taiga.test/api/v1/tasks?page=3")
		resp.Header.Set("X-Pagination-Prev", "http://taiga.test/api/v1/tasks?page=1")
		return resp, nil
	})
	client.DisablePagination(false)

	var tasks []Task
	if _, err := client.Request.Get(client.MakeURL("tasks"), &tasks); err != nil {
		t.Fatalf("request get tasks: %v", err)
	}
	p := client.GetPagination()
	if !p.Paginated {
		t.Fatalf("expected paginated=true")
	}
	if p.PaginatedBy != 50 || p.PaginationCount != 150 || p.PaginationCurrent != 2 {
		t.Fatalf("unexpected pagination values: %+v", p)
	}
	if p.PaginationNext == nil || p.PaginationNext.String() == "" {
		t.Fatalf("expected PaginationNext to be populated")
	}
	if p.PaginationPrev == nil || p.PaginationPrev.String() == "" {
		t.Fatalf("expected PaginationPrev to be populated")
	}
}

func TestResolveProjectRejectsNilQueryParams(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		return newJSONResponse(req, http.StatusOK, `{}`), nil
	})

	if _, err := client.Resolver.ResolveProject(nil); err == nil {
		t.Fatalf("expected ResolveProject(nil) to return error")
	}
}

func TestCloneOperationsDoNotMutateSourceObjects(t *testing.T) {
	type capturedRequest struct {
		path string
		body string
	}
	requests := make([]capturedRequest, 0, 2)

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		payload, _ := io.ReadAll(req.Body)
		requests = append(requests, capturedRequest{
			path: req.URL.Path,
			body: string(payload),
		})
		switch req.URL.Path {
		case "/api/v1/epics":
			return newJSONResponse(req, http.StatusCreated, `{"id":9001,"project":7,"subject":"clone-epic"}`), nil
		case "/api/v1/userstories":
			return newJSONResponse(req, http.StatusCreated, `{"id":9002,"project":7,"subject":"clone-us"}`), nil
		default:
			return newJSONResponse(req, http.StatusOK, `{}`), nil
		}
	})

	epic := &Epic{ID: 101, Ref: 22, Version: 5, Project: 7, Subject: "clone-epic"}
	if _, err := epic.Clone(client.Epic); err != nil {
		t.Fatalf("clone epic: %v", err)
	}
	if epic.ID != 101 || epic.Ref != 22 || epic.Version != 5 {
		t.Fatalf("epic was mutated by Clone: %+v", epic)
	}

	us := &UserStory{ID: 102, Ref: 33, Version: 6, Project: 7, Subject: "clone-us"}
	if _, err := client.UserStory.Clone(us); err != nil {
		t.Fatalf("clone user story: %v", err)
	}
	if us.ID != 102 || us.Ref != 33 || us.Version != 6 {
		t.Fatalf("user story was mutated by Clone: %+v", us)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 create requests for clone operations, got %d", len(requests))
	}
	for _, req := range requests {
		if strings.Contains(req.body, "\"id\":101") || strings.Contains(req.body, "\"id\":102") {
			t.Fatalf("clone request unexpectedly contained original ID in payload: %s", req.body)
		}
		if strings.Contains(req.body, "\"ref\":22") || strings.Contains(req.body, "\"ref\":33") {
			t.Fatalf("clone request unexpectedly contained original Ref in payload: %s", req.body)
		}
		if strings.Contains(req.body, "\"version\":5") || strings.Contains(req.body, "\"version\":6") {
			t.Fatalf("clone request unexpectedly contained original Version in payload: %s", req.body)
		}
	}
}

func TestHistoryDeleteCommentUsesQueryParameter(t *testing.T) {
	var gotPath, gotID, gotMethod string
	var gotBody []byte

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotID = req.URL.Query().Get("id")
		gotMethod = req.Method
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			gotBody = body
		}
		return newJSONResponse(req, http.StatusOK, `{"deleted":true}`), nil
	})

	_, err := client.History.DeleteTaskComment(11, "6ca8f268-fcab-4be0-8927-040cf5fdd95a")
	if err != nil {
		t.Fatalf("delete task comment: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST method, got %s", gotMethod)
	}
	if gotPath != "/api/v1/history/task/11/delete_comment" {
		t.Fatalf("unexpected delete_comment path: %s", gotPath)
	}
	if gotID != "6ca8f268-fcab-4be0-8927-040cf5fdd95a" {
		t.Fatalf("expected comment id in query string, got %q", gotID)
	}
	if len(bytes.TrimSpace(gotBody)) != 0 {
		t.Fatalf("expected empty delete_comment body, got %q", string(gotBody))
	}
}

func TestRequestServiceReturnsReadableResponseBody(t *testing.T) {
	const payload = `{"id":1,"name":"demo"}`

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		return newJSONResponse(req, http.StatusOK, payload), nil
	})

	var out RawResource
	resp, err := client.Request.Get(client.MakeURL("projects", "1"), &out)
	if err != nil {
		t.Fatalf("request get: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body after request helper decode: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if string(body) != payload {
		t.Fatalf("expected response body to remain readable, got %q", string(body))
	}
	if out["name"] != "demo" {
		t.Fatalf("expected decoded response payload, got %+v", out)
	}
}

func TestTaskListAttachmentsBuildsValidQueryURL(t *testing.T) {
	var gotPath, gotObjectID, gotProject string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotObjectID = req.URL.Query().Get("object_id")
		gotProject = req.URL.Query().Get("project")
		return newJSONResponse(req, http.StatusOK, `[{"id":1}]`), nil
	})

	attachments, err := client.Task.ListAttachments(&Task{ID: 8, Project: 4})
	if err != nil {
		t.Fatalf("list task attachments: %v", err)
	}
	if gotPath != "/api/v1/tasks/attachments" {
		t.Fatalf("expected attachments path without malformed query separator, got %s", gotPath)
	}
	if gotObjectID != "8" || gotProject != "4" {
		t.Fatalf("unexpected attachment query params: object_id=%q project=%q", gotObjectID, gotProject)
	}
	if len(attachments) != 1 || attachments[0].ID != 1 {
		t.Fatalf("unexpected attachments payload: %+v", attachments)
	}
}

func TestMakeURLNormalizesEndpointParts(t *testing.T) {
	client := &Client{APIURL: "http://taiga.test/api/v1"}
	got := client.MakeURL("/projects/", "/42/", "userstories")
	if got != "http://taiga.test/api/v1/projects/42/userstories" {
		t.Fatalf("unexpected URL: %s", got)
	}
}

func TestAppendQueryParamsMergesWithExistingQuery(t *testing.T) {
	type qp struct {
		Project int `url:"project,omitempty"`
	}
	got, err := appendQueryParams("http://taiga.test/api/v1/tasks?status=3", &qp{Project: 9})
	if err != nil {
		t.Fatalf("append query params: %v", err)
	}
	if !strings.Contains(got, "status=3") || !strings.Contains(got, "project=9") {
		t.Fatalf("expected merged query string, got %s", got)
	}
}

func TestSetAuthTokensClearsAuthorizationWhenTokenEmpty(t *testing.T) {
	client := &Client{BaseURL: "http://example.com", AutoRefreshDisabled: true}
	if err := client.Initialise(); err != nil {
		t.Fatalf("initialise client: %v", err)
	}
	client.SetAuthTokens(TokenBearer, "token-1", "refresh-1")
	if got := client.GetAuthorizationHeader(); got != "Bearer token-1" {
		t.Fatalf("expected auth header to be set, got %q", got)
	}
	client.SetAuthTokens(TokenBearer, "", "refresh-2")
	if got := client.GetAuthorizationHeader(); got != "" {
		t.Fatalf("expected auth header to be cleared, got %q", got)
	}
}

func TestEpicBulkCreateUsesBulkCreateEndpoint(t *testing.T) {
	var gotPath, gotMethod, gotBody string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		body, _ := io.ReadAll(req.Body)
		gotBody = string(body)
		return newJSONResponse(req, http.StatusCreated, `[{"id":11}]`), nil
	})
	client.Project.ConfigureMappedServices(77)

	_, err := client.Project.Epic.BulkCreate(&EpicBulkCreatePayload{BulkEpics: "Epic A\nEpic B"})
	if err != nil {
		t.Fatalf("bulk create epics: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST method, got %s", gotMethod)
	}
	if gotPath != "/api/v1/epics/bulk_create" {
		t.Fatalf("expected /epics/bulk_create path, got %s", gotPath)
	}
	if !strings.Contains(gotBody, `"project":77`) {
		t.Fatalf("expected default mapped project in payload, got %s", gotBody)
	}
	if !strings.Contains(gotBody, `"bulk_epics":"Epic A`) || !strings.Contains(gotBody, `Epic B"`) {
		t.Fatalf("expected bulk_epics payload, got %s", gotBody)
	}
}

func TestEpicGetFiltersDataUsesProjectQuery(t *testing.T) {
	var gotPath, gotProject string

	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotProject = req.URL.Query().Get("project")
		return newJSONResponse(req, http.StatusOK, `{"statuses":[]}`), nil
	})
	client.Project.ConfigureMappedServices(66)

	if _, err := client.Project.Epic.GetFiltersData(0); err != nil {
		t.Fatalf("get epic filters data: %v", err)
	}
	if gotPath != "/api/v1/epics/filters_data" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotProject != "66" {
		t.Fatalf("expected mapped project query, got %q", gotProject)
	}
}

func TestTaskGetRejectsNonPositiveID(t *testing.T) {
	client := newUnitTestClient(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unexpected request for invalid task ID")
		return nil, nil
	})

	if _, err := client.Task.Get(0); err == nil {
		t.Fatalf("expected validation error for taskID=0")
	}
}
