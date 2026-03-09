package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

type smokeSuite struct {
	client     *taiga.Client
	projectID  int
	projectRef *taiga.Project
}

func newSmokeSuite() *smokeSuite {
	return &smokeSuite{
		client:     Client,
		projectID:  testProjID,
		projectRef: &taiga.Project{ID: testProjID},
	}
}

func (s *smokeSuite) unique(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, strings.ToLower(RandStringBytesMaskImprSrcUnsafe(8)))
}

func (s *smokeSuite) createSupportingUserStory() (*taiga.UserStory, error) {
	return s.client.UserStory.Create(&taiga.UserStory{
		Project: s.projectID,
		Subject: s.unique("smoke-support-us"),
	})
}

type smokeResource struct {
	value   any
	cleanup func() error
}

type smokeCRUDCase struct {
	name          string
	list          func(t *testing.T, s *smokeSuite) error
	create        func(t *testing.T, s *smokeSuite) (smokeResource, error)
	id            func(resource any) int
	get           func(t *testing.T, s *smokeSuite, id int, created any) (any, error)
	update        func(t *testing.T, s *smokeSuite, resource any) (any, error)
	verify        func(t *testing.T, s *smokeSuite, created any, fetched any, updated any)
	delete        func(t *testing.T, s *smokeSuite, resource any) error
	assertDeleted func(t *testing.T, s *smokeSuite, id int) error
}

func TestSmokeCRUDMatrix(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	suite := newSmokeSuite()
	cases := []smokeCRUDCase{
		epicSmokeCase(),
		userStorySmokeCase(),
		taskSmokeCase(),
		issueSmokeCase(),
		wikiSmokeCase(),
		webhookSmokeCase(),
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runSmokeCRUDCase(t, suite, tc)
		})
	}
}

func runSmokeCRUDCase(t *testing.T, suite *smokeSuite, tc smokeCRUDCase) {
	t.Helper()

	if tc.list != nil {
		if err := tc.list(t, suite); err != nil {
			t.Fatalf("list check failed: %v", err)
		}
	}

	created, err := tc.create(t, suite)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.cleanup != nil {
		t.Cleanup(func() {
			if err := ignoreNotFoundError(created.cleanup()); err != nil {
				t.Errorf("cleanup failed: %v", err)
			}
		})
	}

	resourceID := tc.id(created.value)
	if resourceID <= 0 {
		t.Fatalf("invalid resource id extracted from created payload: %d", resourceID)
	}

	fetched := created.value
	if tc.get != nil {
		fetched, err = tc.get(t, suite, resourceID, created.value)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
	}

	updated := fetched
	if tc.update != nil {
		updated, err = tc.update(t, suite, fetched)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
	}

	if tc.verify != nil {
		tc.verify(t, suite, created.value, fetched, updated)
	}

	if tc.delete != nil {
		if err := tc.delete(t, suite, updated); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	}
	if tc.assertDeleted != nil {
		if err := tc.assertDeleted(t, suite, resourceID); err != nil {
			t.Fatalf("post-delete assertion failed: %v", err)
		}
	}
}

func epicSmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "epics",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.Epic.List(&taiga.EpicsQueryParams{Project: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			epic, err := s.client.Epic.Create(&taiga.Epic{
				Project: s.projectID,
				Subject: s.unique("smoke-epic"),
			})
			if err != nil {
				return smokeResource{}, err
			}
			return smokeResource{
				value: epic,
				cleanup: func() error {
					_, err := s.client.Epic.Delete(epic.ID)
					return err
				},
			}, nil
		},
		id: func(resource any) int {
			epic, ok := resource.(*taiga.Epic)
			if !ok || epic == nil {
				return 0
			}
			return epic.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, created any) (any, error) {
			createdEpic := mustCastResource[taiga.Epic](t, created)
			byID, err := s.client.Epic.Get(id)
			if err != nil {
				return nil, err
			}
			byRef, err := s.client.Epic.GetByRef(createdEpic.Ref, s.projectRef)
			if err != nil {
				return nil, err
			}
			if byRef.ID != byID.ID {
				return nil, fmt.Errorf("GetByRef returned id=%d, expected %d", byRef.ID, byID.ID)
			}
			return byID, nil
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			epic := mustCastResource[taiga.Epic](t, resource)
			epic.Description = s.unique("smoke-epic-description")
			return s.client.Epic.Edit(epic)
		},
		verify: func(t *testing.T, _ *smokeSuite, created any, fetched any, updated any) {
			createdEpic := mustCastResource[taiga.Epic](t, created)
			fetchedEpic := mustCastResource[taiga.Epic](t, fetched)
			updatedEpic := mustCastResource[taiga.Epic](t, updated)

			if fetchedEpic.ID != createdEpic.ID {
				t.Fatalf("fetched epic id mismatch: got %d want %d", fetchedEpic.ID, createdEpic.ID)
			}
			if updatedEpic.Description == "" {
				t.Fatalf("expected updated epic description to be non-empty")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			epic := mustCastResource[taiga.Epic](t, resource)
			_, err := s.client.Epic.Delete(epic.ID)
			return err
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.Epic.Get(id)
			return expectResourceMissing(err, "epic")
		},
	}
}

func userStorySmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "userstories",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.UserStory.List(&taiga.UserStoryQueryParams{Project: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			userStory, err := s.client.UserStory.Create(&taiga.UserStory{
				Project: s.projectID,
				Subject: s.unique("smoke-us"),
			})
			if err != nil {
				return smokeResource{}, err
			}
			return smokeResource{
				value: userStory,
				cleanup: func() error {
					_, err := s.client.UserStory.Delete(userStory.ID)
					return err
				},
			}, nil
		},
		id: func(resource any) int {
			userStory, ok := resource.(*taiga.UserStory)
			if !ok || userStory == nil {
				return 0
			}
			return userStory.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, created any) (any, error) {
			createdUS := mustCastResource[taiga.UserStory](t, created)
			byID, err := s.client.UserStory.Get(id)
			if err != nil {
				return nil, err
			}
			byRef, err := s.client.UserStory.GetByRef(createdUS.Ref, s.projectRef)
			if err != nil {
				return nil, err
			}
			if byRef.ID != byID.ID {
				return nil, fmt.Errorf("GetByRef returned id=%d, expected %d", byRef.ID, byID.ID)
			}
			return byID, nil
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			userStory := mustCastResource[taiga.UserStory](t, resource)
			userStory.Subject = s.unique("smoke-us-updated")
			return s.client.UserStory.Edit(userStory)
		},
		verify: func(t *testing.T, _ *smokeSuite, created any, fetched any, updated any) {
			createdUS := mustCastResource[taiga.UserStory](t, created)
			fetchedUS := mustCastResource[taiga.UserStory](t, fetched)
			updatedUS := mustCastResource[taiga.UserStory](t, updated)

			if fetchedUS.ID != createdUS.ID {
				t.Fatalf("fetched user story id mismatch: got %d want %d", fetchedUS.ID, createdUS.ID)
			}
			if updatedUS.Subject == createdUS.Subject {
				t.Fatalf("expected updated user story subject to differ from create step")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			userStory := mustCastResource[taiga.UserStory](t, resource)
			_, err := s.client.UserStory.Delete(userStory.ID)
			return err
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.UserStory.Get(id)
			return expectResourceMissing(err, "user story")
		},
	}
}

func taskSmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "tasks",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.Task.List(&taiga.TasksQueryParams{Project: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			userStory, err := s.createSupportingUserStory()
			if err != nil {
				return smokeResource{}, err
			}
			task, err := s.client.Task.Create(&taiga.Task{
				Project:   s.projectID,
				UserStory: userStory.ID,
				Subject:   s.unique("smoke-task"),
			})
			if err != nil {
				_, _ = s.client.UserStory.Delete(userStory.ID)
				return smokeResource{}, err
			}
			return smokeResource{
				value: task,
				cleanup: func() error {
					_, errTask := s.client.Task.Delete(task.ID)
					if err := ignoreNotFoundError(errTask); err != nil {
						return err
					}
					_, errUS := s.client.UserStory.Delete(userStory.ID)
					return ignoreNotFoundError(errUS)
				},
			}, nil
		},
		id: func(resource any) int {
			task, ok := resource.(*taiga.Task)
			if !ok || task == nil {
				return 0
			}
			return task.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, created any) (any, error) {
			createdTask := mustCastResource[taiga.Task](t, created)
			byID, err := s.client.Task.Get(id)
			if err != nil {
				return nil, err
			}
			byRef, err := s.client.Task.GetByRef(createdTask.Ref, s.projectRef)
			if err != nil {
				return nil, err
			}
			if byRef.ID != byID.ID {
				return nil, fmt.Errorf("GetByRef returned id=%d, expected %d", byRef.ID, byID.ID)
			}
			return byID, nil
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			task := mustCastResource[taiga.Task](t, resource)
			task.Description = s.unique("smoke-task-description")
			return s.client.Task.Edit(task)
		},
		verify: func(t *testing.T, _ *smokeSuite, created any, fetched any, updated any) {
			createdTask := mustCastResource[taiga.Task](t, created)
			fetchedTask := mustCastResource[taiga.Task](t, fetched)
			updatedTask := mustCastResource[taiga.Task](t, updated)

			if fetchedTask.ID != createdTask.ID {
				t.Fatalf("fetched task id mismatch: got %d want %d", fetchedTask.ID, createdTask.ID)
			}
			if updatedTask.Description == "" {
				t.Fatalf("expected updated task description to be non-empty")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			task := mustCastResource[taiga.Task](t, resource)
			_, err := s.client.Task.Delete(task.ID)
			return err
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.Task.Get(id)
			return expectResourceMissing(err, "task")
		},
	}
}

func issueSmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "issues",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.Issue.List(&taiga.IssueQueryParams{Project: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			issue, err := s.client.Issue.Create(&taiga.Issue{
				Project: s.projectID,
				Subject: s.unique("smoke-issue"),
			})
			if err != nil {
				return smokeResource{}, err
			}
			return smokeResource{
				value: issue,
				cleanup: func() error {
					_, err := s.client.Issue.Delete(issue.ID)
					return err
				},
			}, nil
		},
		id: func(resource any) int {
			issue, ok := resource.(*taiga.Issue)
			if !ok || issue == nil {
				return 0
			}
			return issue.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, created any) (any, error) {
			createdIssue := mustCastResource[taiga.Issue](t, created)
			byID, err := s.client.Issue.Get(id)
			if err != nil {
				return nil, err
			}
			byRef, err := s.client.Issue.GetByRef(createdIssue.Ref, s.projectRef)
			if err != nil {
				return nil, err
			}
			if byRef.ID != byID.ID {
				return nil, fmt.Errorf("GetByRef returned id=%d, expected %d", byRef.ID, byID.ID)
			}
			return byID, nil
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			issue := mustCastResource[taiga.Issue](t, resource)
			issue.Description = s.unique("smoke-issue-description")
			return s.client.Issue.Edit(issue)
		},
		verify: func(t *testing.T, _ *smokeSuite, created any, fetched any, updated any) {
			createdIssue := mustCastResource[taiga.Issue](t, created)
			fetchedIssue := mustCastResource[taiga.Issue](t, fetched)
			updatedIssue := mustCastResource[taiga.Issue](t, updated)

			if fetchedIssue.ID != createdIssue.ID {
				t.Fatalf("fetched issue id mismatch: got %d want %d", fetchedIssue.ID, createdIssue.ID)
			}
			if updatedIssue.Description == "" {
				t.Fatalf("expected updated issue description to be non-empty")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			issue := mustCastResource[taiga.Issue](t, resource)
			_, err := s.client.Issue.Delete(issue.ID)
			return err
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.Issue.Get(id)
			return expectResourceMissing(err, "issue")
		},
	}
}

func wikiSmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "wiki",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.Wiki.List(&taiga.WikiQueryParams{Project: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			page, err := s.client.Wiki.Create(&taiga.WikiPage{
				Project: s.projectID,
				Slug:    s.unique("smoke-wiki"),
				Content: "# smoke wiki page",
			})
			if err != nil {
				return smokeResource{}, err
			}
			return smokeResource{
				value: page,
				cleanup: func() error {
					_, err := s.client.Wiki.Delete(page.ID)
					return err
				},
			}, nil
		},
		id: func(resource any) int {
			page, ok := resource.(*taiga.WikiPage)
			if !ok || page == nil {
				return 0
			}
			return page.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, created any) (any, error) {
			createdPage := mustCastResource[taiga.WikiPage](t, created)
			byID, err := s.client.Wiki.Get(id)
			if err != nil {
				return nil, err
			}
			bySlug, err := s.client.Wiki.GetBySlug(createdPage.Slug, s.projectID)
			if err != nil {
				return nil, err
			}
			if bySlug.ID != byID.ID {
				return nil, fmt.Errorf("GetBySlug returned id=%d, expected %d", bySlug.ID, byID.ID)
			}
			return byID, nil
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			page := mustCastResource[taiga.WikiPage](t, resource)
			page.Content = "# " + s.unique("smoke-wiki-updated")
			return s.client.Wiki.Edit(page)
		},
		verify: func(t *testing.T, _ *smokeSuite, created any, fetched any, updated any) {
			createdPage := mustCastResource[taiga.WikiPage](t, created)
			fetchedPage := mustCastResource[taiga.WikiPage](t, fetched)
			updatedPage := mustCastResource[taiga.WikiPage](t, updated)

			if fetchedPage.ID != createdPage.ID {
				t.Fatalf("fetched wiki id mismatch: got %d want %d", fetchedPage.ID, createdPage.ID)
			}
			if updatedPage.Content == createdPage.Content {
				t.Fatalf("expected updated wiki content to differ from create step")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			page := mustCastResource[taiga.WikiPage](t, resource)
			_, err := s.client.Wiki.Delete(page.ID)
			return err
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.Wiki.Get(id)
			return expectResourceMissing(err, "wiki page")
		},
	}
}

func webhookSmokeCase() smokeCRUDCase {
	return smokeCRUDCase{
		name: "webhooks",
		list: func(t *testing.T, s *smokeSuite) error {
			_, err := s.client.Webhook.List(&taiga.WebhookQueryParameters{ProjectID: s.projectID})
			return err
		},
		create: func(t *testing.T, s *smokeSuite) (smokeResource, error) {
			webhook, err := s.client.Webhook.Create(&taiga.Webhook{
				Project: s.projectID,
				Name:    s.unique("smoke-webhook"),
				Key:     s.unique("smoke-webhook-key"),
				URL:     "https://example.com/" + s.unique("taigo"),
			})
			if err != nil {
				return smokeResource{}, err
			}
			return smokeResource{
				value: webhook,
				cleanup: func() error {
					return s.client.Webhook.Delete(webhook.ID)
				},
			}, nil
		},
		id: func(resource any) int {
			webhook, ok := resource.(*taiga.Webhook)
			if !ok || webhook == nil {
				return 0
			}
			return webhook.ID
		},
		get: func(t *testing.T, s *smokeSuite, id int, _ any) (any, error) {
			return s.client.Webhook.Get(id)
		},
		update: func(t *testing.T, s *smokeSuite, resource any) (any, error) {
			webhook := mustCastResource[taiga.Webhook](t, resource)
			webhook.Name = s.unique("smoke-webhook-updated")
			return s.client.Webhook.Update(webhook)
		},
		verify: func(t *testing.T, s *smokeSuite, created any, fetched any, updated any) {
			createdWebhook := mustCastResource[taiga.Webhook](t, created)
			fetchedWebhook := mustCastResource[taiga.Webhook](t, fetched)
			updatedWebhook := mustCastResource[taiga.Webhook](t, updated)

			if fetchedWebhook.ID != createdWebhook.ID {
				t.Fatalf("fetched webhook id mismatch: got %d want %d", fetchedWebhook.ID, createdWebhook.ID)
			}
			if updatedWebhook.Name == createdWebhook.Name {
				t.Fatalf("expected updated webhook name to differ from create step")
			}

			logEntry, err := s.client.Webhook.Test(updatedWebhook.ID)
			if err != nil {
				t.Fatalf("webhook test endpoint failed: %v", err)
			}
			if logEntry.ID == 0 {
				t.Fatalf("webhook test endpoint did not return a log id")
			}
		},
		delete: func(t *testing.T, s *smokeSuite, resource any) error {
			webhook := mustCastResource[taiga.Webhook](t, resource)
			return s.client.Webhook.Delete(webhook.ID)
		},
		assertDeleted: func(t *testing.T, s *smokeSuite, id int) error {
			_, err := s.client.Webhook.Get(id)
			return expectResourceMissing(err, "webhook")
		},
	}
}

func ignoreNotFoundError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *taiga.APIError
	if errors.As(err, &apiErr) && (apiErr.StatusCode == http.StatusNotFound || apiErr.StatusCode == http.StatusGone) {
		return nil
	}
	return err
}

func expectResourceMissing(err error, resourceName string) error {
	if err == nil {
		return fmt.Errorf("%s still retrievable after delete", resourceName)
	}
	var apiErr *taiga.APIError
	if errors.As(err, &apiErr) && (apiErr.StatusCode == http.StatusNotFound || apiErr.StatusCode == http.StatusGone) {
		return nil
	}
	return fmt.Errorf("unexpected %s get error after delete: %w", resourceName, err)
}

func mustCastResource[T any](t *testing.T, resource any) *T {
	t.Helper()

	value, ok := resource.(*T)
	if !ok || value == nil {
		t.Fatalf("unexpected resource type: %T", resource)
	}
	return value
}
