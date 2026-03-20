package taigo

import (
	"testing"

	"github.com/google/go-querystring/query"
)

type queryMatrixCase struct {
	name   string
	params any
	want   map[string]string
}

func TestQueryFilterMatrixEncoding(t *testing.T) {
	projectsQuery := &ProjectsQueryParameters{
		Member:     7,
		IsFeatured: BoolPtr(false),
	}
	projectsQuery.TotalActivityLastMonth()

	taskQuery := &TasksQueryParams{
		Project:            2,
		StatusIsClosed:     BoolPtr(false),
		IncludeAttachments: BoolPtr(true),
	}
	taskQuery.SetTags("backend", "api")

	cases := []queryMatrixCase{
		{
			name:   "projects/query",
			params: projectsQuery,
			want: map[string]string{
				"member":      "7",
				"is_featured": "false",
				"order_by":    "total_activity_last_month",
			},
		},
		{
			name:   "milestones/query",
			params: &MilestonesQueryParams{Project: 2, Closed: BoolPtr(false)},
			want: map[string]string{
				"project": "2",
				"closed":  "false",
			},
		},
		{
			name: "epics/query",
			params: &EpicsQueryParams{
				Project:            2,
				AssignedTo:         5,
				IncludeAttachments: BoolPtr(true),
				StatusIsClosed:     BoolPtr(false),
			},
			want: map[string]string{
				"project":             "2",
				"assigned_to":         "5",
				"include_attachments": "true",
				"status__is_closed":   "false",
			},
		},
		{
			name: "userstories/query",
			params: &UserStoryQueryParams{
				Project:            2,
				AssignedTo:         5,
				IncludeTasks:       BoolPtr(true),
				StatusIsClosed:     BoolPtr(false),
				IncludeAttachments: BoolPtr(false),
			},
			want: map[string]string{
				"project":             "2",
				"assigned_to":         "5",
				"include_tasks":       "true",
				"status__is_closed":   "false",
				"include_attachments": "false",
			},
		},
		{
			name:   "tasks/query",
			params: taskQuery,
			want: map[string]string{
				"project":             "2",
				"status__is_closed":   "false",
				"include_attachments": "true",
				"tags":                "backend,api",
			},
		},
		{
			name: "issues/query",
			params: &IssueQueryParams{
				Project:            2,
				Owner:              3,
				StatusIsClosed:     BoolPtr(true),
				IncludeAttachments: BoolPtr(false),
			},
			want: map[string]string{
				"project":             "2",
				"owner":               "3",
				"status__is_closed":   "true",
				"include_attachments": "false",
			},
		},
		{
			name:   "wiki/query",
			params: &WikiQueryParams{Project: 2, Slug: "my-page"},
			want: map[string]string{
				"project": "2",
				"slug":    "my-page",
			},
		},
		{
			name:   "webhooks/query",
			params: &WebhookQueryParameters{ProjectID: 2, WebhookID: 9},
			want: map[string]string{
				"project": "2",
				"webhook": "9",
			},
		},
		{
			name:   "memberships/query",
			params: &MembershipsQueryParams{Project: 2, Role: 4},
			want: map[string]string{
				"project": "2",
				"role":    "4",
			},
		},
		{
			name:   "membership_invitations/query",
			params: &MembershipInvitationsQueryParams{Project: 2},
			want: map[string]string{
				"project": "2",
			},
		},
		{
			name:   "timeline/query",
			params: &TimelineQueryParams{Page: 3},
			want: map[string]string{
				"page": "3",
			},
		},
		{
			name:   "importer_auth/query",
			params: &ImporterAuthURLQueryParams{Project: 2},
			want: map[string]string{
				"project": "2",
			},
		},
		{
			name:   "users_highlighted/query",
			params: &UsersHighlightedQueryParams{Type: "task", Q: "bug"},
			want: map[string]string{
				"type": "task",
				"q":    "bug",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			values, err := query.Values(tc.params)
			if err != nil {
				t.Fatalf("query.Values returned error: %v", err)
			}
			for key, want := range tc.want {
				if got := values.Get(key); got != want {
					t.Fatalf("query key %q mismatch: got %q want %q", key, got, want)
				}
			}
		})
	}
}

func TestQueryFilterMatrixOmitsNilPointers(t *testing.T) {
	values, err := query.Values(&MilestonesQueryParams{Project: 2})
	if err != nil {
		t.Fatalf("query.Values returned error: %v", err)
	}
	if got := values.Get("closed"); got != "" {
		t.Fatalf("expected closed to be omitted when nil, got %q", got)
	}
}
