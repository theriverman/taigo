package main

import (
	"fmt"
	"os"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestTasks(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	us, err := Client.UserStory.Create(&taiga.UserStory{Project: testProjID, Subject: "US for task testing"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = Client.UserStory.Delete(us.ID)
	}()

	task, err := Client.Task.Create(&taiga.Task{Project: testProjID, UserStory: us.ID, Subject: "Task created by Taigo tests"})
	if err != nil {
		t.Fatal(err)
	}

	// List
	tasks, err := Client.Task.List(&taiga.TasksQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	}
	if len(tasks) == 0 {
		t.Errorf("expected at least one task in list")
	}

	// Get
	taskByID, err := Client.Task.Get(task.ID)
	if err != nil {
		t.Error(err)
	}
	if taskByID.ID != task.ID {
		t.Errorf("got %d, want %d", taskByID.ID, task.ID)
	}

	// GetByRef
	taskByRef, err := Client.Task.GetByRef(task.Ref, &taiga.Project{ID: testProjID})
	if err != nil {
		t.Error(err)
	}
	if taskByRef.ID != task.ID {
		t.Errorf("got %d, want %d", taskByRef.ID, task.ID)
	}

	// Edit
	task.Description = "Task edited via integration test"
	editedTask, err := Client.Task.Edit(task)
	if err != nil {
		t.Error(err)
	}
	if editedTask.Description != "Task edited via integration test" {
		t.Errorf("got %q, want %q", editedTask.Description, "Task edited via integration test")
	}

	// Attachment create/get/list
	testFileName := "initial_test_data.json"
	attachment := &taiga.Attachment{Name: testFileName, Description: "Task attachment from integration test"}
	attachment.SetFilePath(fmt.Sprintf("%s%s%s", cwd, string(os.PathSeparator), testFileName))
	createdAttachment, err := Client.Task.CreateAttachment(attachment, task)
	if err != nil {
		t.Error(err)
	}

	attachmentByID, err := Client.Task.GetAttachment(createdAttachment.ID)
	if err != nil {
		t.Error(err)
	}
	if attachmentByID.ID != createdAttachment.ID {
		t.Errorf("got %d, want %d", attachmentByID.ID, createdAttachment.ID)
	}

	attachments, err := Client.Task.ListAttachments(task)
	if err != nil {
		t.Error(err)
	}
	if len(attachments) == 0 {
		t.Errorf("expected at least one attachment")
	}

	// Delete
	if _, err := Client.Task.Delete(task.ID); err != nil {
		t.Error(err)
	}
}
