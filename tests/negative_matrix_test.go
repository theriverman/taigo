package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestNegativeMatrixLive(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	type negativeCase struct {
		name          string
		invoke        func() error
		wantStatuses  []int
		wantSubstring string
	}

	cases := []negativeCase{
		{
			name: "epic/get_non_existing",
			invoke: func() error {
				_, err := Client.Epic.Get(99999999)
				return err
			},
			wantStatuses: []int{http.StatusNotFound},
		},
		{
			name: "userstory/get_non_existing",
			invoke: func() error {
				_, err := Client.UserStory.Get(99999999)
				return err
			},
			wantStatuses: []int{http.StatusNotFound},
		},
		{
			name: "task/get_non_existing",
			invoke: func() error {
				_, err := Client.Task.Get(99999999)
				return err
			},
			wantStatuses: []int{http.StatusNotFound},
		},
		{
			name: "issue/get_non_existing",
			invoke: func() error {
				_, err := Client.Issue.Get(99999999)
				return err
			},
			wantStatuses: []int{http.StatusNotFound},
		},
		{
			name: "task/create_missing_subject_validation",
			invoke: func() error {
				_, err := Client.Task.Create(&taiga.Task{Project: testProjID})
				return err
			},
		},
		{
			name: "userstory/create_missing_subject_validation",
			invoke: func() error {
				_, err := Client.UserStory.Create(&taiga.UserStory{Project: testProjID})
				return err
			},
		},
		{
			name: "wiki/create_missing_slug_validation",
			invoke: func() error {
				_, err := Client.Wiki.Create(&taiga.WikiPage{Project: testProjID, Content: "# content"})
				return err
			},
		},
		{
			name: "webhook/create_missing_key_validation",
			invoke: func() error {
				_, err := Client.Webhook.Create(&taiga.Webhook{
					Project: testProjID,
					Name:    "negative-webhook",
					URL:     "https://example.com/negative-webhook",
				})
				return err
			},
			wantSubstring: "key is required",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.invoke()
			if err == nil {
				t.Fatalf("expected an error in negative test")
			}
			if tc.wantSubstring != "" && !strings.Contains(err.Error(), tc.wantSubstring) {
				t.Fatalf("expected error containing %q, got %v", tc.wantSubstring, err)
			}
			if len(tc.wantStatuses) == 0 {
				return
			}
			if err := assertAPIErrorStatusIn(err, tc.wantStatuses...); err != nil {
				t.Fatalf("unexpected error type/status: %v", err)
			}
		})
	}
}

func assertAPIErrorStatusIn(err error, statuses ...int) error {
	if err == nil {
		return fmt.Errorf("missing error")
	}
	var apiErr *taiga.APIError
	if !errors.As(err, &apiErr) {
		return fmt.Errorf("expected *taiga.APIError, got %T (%v)", err, err)
	}
	for _, status := range statuses {
		if apiErr.StatusCode == status {
			return nil
		}
	}
	return fmt.Errorf("unexpected API status %d, expected one of %v", apiErr.StatusCode, statuses)
}
