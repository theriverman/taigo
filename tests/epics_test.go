package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestEpics(t *testing.T) {
	setupClient()

	// List Epics
	epics, err := comTestClient.Epic.List(&taiga.EpicsQueryParams{Project: dummyProjectID})
	if err != nil {
		t.Error(err)
	}
	if len(epics) == 0 {
		t.Error("Returned Epic list was empty.")
	}

	// Create Epic
	epic := &taiga.Epic{
		Subject: "Test Epic by Taigo",
		Project: dummyProjectID,
	}
	epic, err = comTestClient.Epic.Create(epic)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Edit Epic
	epic.Description = "Added some text here via Client.Epic.Edit()"
	epic, err = comTestClient.Epic.Edit(epic)
	if err != nil {
		t.Error(err)
	}

	// Get Epic
	e1, err := comTestClient.Epic.Get(epic.ID)
	if (err != nil) || (e1.ID != epic.ID) {
		t.Error(err)
	}

	// Get Epic by Ref
	e2, err := comTestClient.Epic.GetByRef(epic.Ref, &taiga.Project{ID: dummyProjectID})
	if (err != nil) || (e2.ID != epic.ID) {
		t.Error(err)
	}

	// Delete Epic by ID
	_, err = comTestClient.Epic.Delete(epic.ID)
	if err != nil {
		t.Error(err)
	}

	// Destroy taiga.Client{}
	teardownClient()
}
