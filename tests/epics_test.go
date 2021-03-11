package main

import (
	"fmt"
	"os"
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestEpics(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// List Epics
	epics, err := Client.Epic.List(&taiga.EpicsQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Total Epics: %d", len(epics))
	}

	// Create Epic
	epic := &taiga.Epic{
		Subject: "Test Epic by Taigo",
		Project: testProjID,
	}
	epic, err = Client.Epic.Create(epic)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Edit Epic
	epic.Description = "Added some text here via Client.Epic.Edit()"
	epic, err = Client.Epic.Edit(epic)
	if err != nil {
		t.Error(err)
	}

	// Get Epic
	e1, err := Client.Epic.Get(epic.ID)
	if (err != nil) || (e1.ID != epic.ID) {
		t.Error(err)
	}

	// Get Epic by Ref
	e2, err := Client.Epic.GetByRef(epic.Ref, &taiga.Project{ID: testProjID})
	if (err != nil) || (e2.ID != epic.ID) {
		t.Error(err)
	}

	// Edit Epic
	newEpicSubject := "This is the updated Subject"
	epicCopyBase := *e2
	epicCopyBase.Subject = newEpicSubject
	epicPatched, err := Client.Epic.Edit(&epicCopyBase)
	if err != nil {
		t.Error(err)
	}
	if epicPatched.Subject != newEpicSubject {
		t.Errorf("got %q, want %q", epicPatched.Subject, newEpicSubject)
	}

	/*
		Testing `ListRelatedUserStories` & `CreateRelatedUserStory`
		* An Epic is needed, so we create one
		* A User Story is needed, so we create one
		* We connect this UserStory to our Epic with `CreateRelatedUserStory`
		* We list the related USs with `ListRelatedUserStories` which should return a total of 1 US
	*/
	epicForUs, err := Client.Epic.Create(&taiga.Epic{Project: testProjID, Subject: "A regular Epic"})
	if err != nil {
		t.Error(err)
	}
	usToBeRelated, err := Client.UserStory.Create(&taiga.UserStory{Project: testProjID, Subject: "A US related to an Epic"})
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if _, err := Client.UserStory.Delete(usToBeRelated.ID); err != nil {
			t.Error(err)
		}
	}()

	// Create a Related UserStory
	_, err = Client.Epic.CreateRelatedUserStory(epicForUs.ID, usToBeRelated.ID)
	if err != nil {
		t.Error(err)
	}

	// List Related User Stories
	relatedUsList, err := Client.Epic.ListRelatedUserStories(epicForUs.ID)
	if err != nil {
		t.Error(err)
	}
	totalNoOfUs := len(relatedUsList)
	if totalNoOfUs != 1 {
		t.Errorf("got %q, want %q", totalNoOfUs, 1)
	}

	// Create an Epic Attachment
	attachment := &taiga.Attachment{
		Name:        "A random project file",
		Description: "This is a test file uploaded via TAIGO",
	}
	testFileName := "initial_test_data.json"
	attachment.SetFilePath(fmt.Sprintf("%s%s%s", cwd, string(os.PathSeparator), testFileName))
	attachmentDetails, err := Client.Epic.CreateAttachment(attachment, epicForUs)
	if err != nil {
		t.Error(err)
	}

	if attachmentDetails.Name != testFileName {
		t.Errorf("got %q, want %q", attachmentDetails.Name, testFileName)
	}

	// Delete Epic by ID
	for _, e := range []taiga.Epic{*epic, *epicForUs} {
		_, err = Client.Epic.Delete(e.ID)
		if err != nil {
			t.Error(err)
		}
	}

}
