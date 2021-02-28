package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestEpics(t *testing.T) {
	setupClient()

	// List Epics
	_, err := Client.Epic.List(&taiga.EpicsQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
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

	// Delete Epic by ID
	_, err = Client.Epic.Delete(epic.ID)
	if err != nil {
		t.Error(err)
	}

	// Destroy taiga.Client{}
	teardownClient()
}
