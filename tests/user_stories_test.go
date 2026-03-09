package main

import (
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestUserStories(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	// List
	userStories, err := Client.UserStory.List(&taiga.UserStoryQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	}
	t.Logf("Total User Stories: %d", len(userStories))

	// Create
	userStory, err := Client.UserStory.Create(&taiga.UserStory{Project: testProjID, Subject: "US created by Taigo tests"})
	if err != nil {
		t.Fatal(err)
	}

	// Get
	usByID, err := Client.UserStory.Get(userStory.ID)
	if err != nil {
		t.Error(err)
	}
	if usByID.ID != userStory.ID {
		t.Errorf("got %d, want %d", usByID.ID, userStory.ID)
	}

	// GetByRef
	usByRef, err := Client.UserStory.GetByRef(userStory.Ref, &taiga.Project{ID: testProjID})
	if err != nil {
		t.Error(err)
	}
	if usByRef.ID != userStory.ID {
		t.Errorf("got %d, want %d", usByRef.ID, userStory.ID)
	}

	// Edit
	newSubject := "US edited by Taigo tests"
	userStory.Subject = newSubject
	editedUS, err := Client.UserStory.Edit(userStory)
	if err != nil {
		t.Error(err)
	}
	if editedUS.Subject != newSubject {
		t.Errorf("got %q, want %q", editedUS.Subject, newSubject)
	}

	// Related task flow
	relatedTask, err := userStory.CreateRelatedTask(Client, taiga.Task{Subject: "Task related to created US"})
	if err != nil {
		t.Error(err)
	}

	relatedTasks, err := userStory.GetRelatedTasks(Client)
	if err != nil {
		t.Error(err)
	}
	if len(relatedTasks) == 0 {
		t.Errorf("expected related tasks to include the created one")
	}

	if _, err := Client.Task.Delete(relatedTask.ID); err != nil {
		t.Error(err)
	}

	// Delete
	if _, err := Client.UserStory.Delete(userStory.ID); err != nil {
		t.Error(err)
	}
}
