package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

type cleanupStep struct {
	name string
	fn   func() error
}

type cleanupStack struct {
	steps []cleanupStep
}

func (c *cleanupStack) add(name string, fn func() error) {
	c.steps = append(c.steps, cleanupStep{name: name, fn: fn})
}

func (c *cleanupStack) run(t *testing.T) {
	t.Helper()
	for i := len(c.steps) - 1; i >= 0; i-- {
		step := c.steps[i]
		if err := ignoreMissingLiveError(step.fn()); err != nil {
			t.Errorf("cleanup step %q failed: %v", step.name, err)
		}
	}
}

func TestWorkflowMatrixLive(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	cleanups := &cleanupStack{}
	t.Cleanup(func() { cleanups.run(t) })

	runID := RandStringBytesMaskImprSrcUnsafe(8)
	workflowPrefix := "workflow-" + runID

	epic, err := Client.Epic.Create(&taiga.Epic{
		Project: testProjID,
		Subject: workflowPrefix + "-epic",
	})
	if err != nil {
		t.Fatalf("create epic failed: %v", err)
	}
	cleanups.add("delete epic", func() error {
		_, err := Client.Epic.Delete(epic.ID)
		return err
	})

	userStory, err := Client.UserStory.Create(&taiga.UserStory{
		Project: testProjID,
		Subject: workflowPrefix + "-userstory",
	})
	if err != nil {
		t.Fatalf("create user story failed: %v", err)
	}
	cleanups.add("delete user story", func() error {
		_, err := Client.UserStory.Delete(userStory.ID)
		return err
	})

	task, err := Client.Task.Create(&taiga.Task{
		Project:   testProjID,
		UserStory: userStory.ID,
		Subject:   workflowPrefix + "-task",
	})
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	cleanups.add("delete task", func() error {
		_, err := Client.Task.Delete(task.ID)
		return err
	})

	issue, err := Client.Issue.Create(&taiga.Issue{
		Project: testProjID,
		Subject: workflowPrefix + "-issue",
	})
	if err != nil {
		t.Fatalf("create issue failed: %v", err)
	}
	cleanups.add("delete issue", func() error {
		_, err := Client.Issue.Delete(issue.ID)
		return err
	})

	wikiSlug := workflowPrefix + "-wiki"
	wikiPage, err := Client.Wiki.Create(&taiga.WikiPage{
		Project: testProjID,
		Slug:    wikiSlug,
		Content: "# " + workflowPrefix + " wiki",
	})
	if err != nil {
		t.Fatalf("create wiki page failed: %v", err)
	}
	cleanups.add("delete wiki", func() error {
		_, err := Client.Wiki.Delete(wikiPage.ID)
		return err
	})

	webhook, err := Client.Webhook.Create(&taiga.Webhook{
		Project: testProjID,
		Name:    workflowPrefix + "-webhook",
		Key:     workflowPrefix + "-key",
		URL:     "https://example.com/" + workflowPrefix + "/webhook",
	})
	if err != nil {
		t.Fatalf("create webhook failed: %v", err)
	}
	cleanups.add("delete webhook", func() error {
		return Client.Webhook.Delete(webhook.ID)
	})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}
	attachmentFile := filepath.Join(cwd, "initial_test_data.json")

	taskAttachment := &taiga.Attachment{Name: workflowPrefix + "-task-attachment"}
	taskAttachment.SetFilePath(attachmentFile)
	if _, err := Client.Task.CreateAttachment(taskAttachment, task); err != nil {
		t.Fatalf("create task attachment failed: %v", err)
	}

	wikiAttachment := &taiga.Attachment{Name: workflowPrefix + "-wiki-attachment"}
	wikiAttachment.SetFilePath(attachmentFile)
	if _, err := Client.Wiki.CreateAttachment(wikiAttachment, wikiPage); err != nil {
		t.Fatalf("create wiki attachment failed: %v", err)
	}

	if _, err := Client.History.UserStory(userStory.ID); err != nil {
		t.Fatalf("history userstory failed: %v", err)
	}
	if _, err := Client.History.Task(task.ID); err != nil {
		t.Fatalf("history task failed: %v", err)
	}
	if _, err := Client.History.Issue(issue.ID); err != nil {
		t.Fatalf("history issue failed: %v", err)
	}
	if _, err := Client.History.Wiki(wikiPage.ID); err != nil {
		t.Fatalf("history wiki failed: %v", err)
	}
	if _, err := Client.Timeline.Project(testProjID, &taiga.TimelineQueryParams{Page: 1}); err != nil {
		t.Fatalf("timeline project failed: %v", err)
	}

	epic.Description = workflowPrefix + "-epic-description"
	if _, err := Client.Epic.Update(epic); err != nil {
		t.Fatalf("update epic failed: %v", err)
	}

	userStorySubject := workflowPrefix + "-userstory-updated"
	if _, err := Client.UserStory.Patch(userStory.ID, &taiga.UserStoryPatch{
		Version: userStory.Version,
		Subject: &userStorySubject,
	}); err != nil {
		t.Fatalf("update user story failed: %v", err)
	}

	task.Description = workflowPrefix + "-task-description"
	if _, err := Client.Task.Update(task); err != nil {
		t.Fatalf("update task failed: %v", err)
	}

	issueDescription := workflowPrefix + "-issue-description"
	if _, err := Client.Issue.Patch(issue.ID, &taiga.IssuePatch{
		Version:     issue.Version,
		Description: &issueDescription,
	}); err != nil {
		t.Fatalf("update issue failed: %v", err)
	}

	wikiPage.Content = "# " + workflowPrefix + " wiki updated"
	if _, err := Client.Wiki.Update(wikiPage); err != nil {
		t.Fatalf("update wiki failed: %v", err)
	}

	webhook.Name = workflowPrefix + "-webhook-updated"
	if _, err := Client.Webhook.Update(webhook); err != nil {
		t.Fatalf("update webhook failed: %v", err)
	}

	if _, err := Client.Task.Delete(task.ID); err != nil {
		t.Fatalf("delete task failed: %v", err)
	}
	if _, err := Client.UserStory.Delete(userStory.ID); err != nil {
		t.Fatalf("delete user story failed: %v", err)
	}
	if _, err := Client.Issue.Delete(issue.ID); err != nil {
		t.Fatalf("delete issue failed: %v", err)
	}
	if _, err := Client.Wiki.Delete(wikiPage.ID); err != nil {
		t.Fatalf("delete wiki failed: %v", err)
	}
	if err := Client.Webhook.Delete(webhook.ID); err != nil {
		t.Fatalf("delete webhook failed: %v", err)
	}
	if _, err := Client.Epic.Delete(epic.ID); err != nil {
		t.Fatalf("delete epic failed: %v", err)
	}

	assertMissingAfterDelete(t, "task", func() error { _, err := Client.Task.Get(task.ID); return err })
	assertMissingAfterDelete(t, "user story", func() error { _, err := Client.UserStory.Get(userStory.ID); return err })
	assertMissingAfterDelete(t, "issue", func() error { _, err := Client.Issue.Get(issue.ID); return err })
	assertMissingAfterDelete(t, "wiki", func() error { _, err := Client.Wiki.Get(wikiPage.ID); return err })
	assertMissingAfterDelete(t, "webhook", func() error { _, err := Client.Webhook.Get(webhook.ID); return err })
	assertMissingAfterDelete(t, "epic", func() error { _, err := Client.Epic.Get(epic.ID); return err })
}

func assertMissingAfterDelete(t *testing.T, resourceName string, getFn func() error) {
	t.Helper()
	err := getFn()
	if err == nil {
		t.Fatalf("expected %s to be missing after delete", resourceName)
	}
	if err := assertAPIErrorStatusIn(err, http.StatusNotFound, http.StatusGone); err != nil {
		t.Fatalf("%s missing check failed: %v", resourceName, err)
	}
}

func ignoreMissingLiveError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *taiga.APIError
	if errors.As(err, &apiErr) && (apiErr.StatusCode == http.StatusNotFound || apiErr.StatusCode == http.StatusGone) {
		return nil
	}
	return fmt.Errorf("non-ignorable cleanup error: %w", err)
}
