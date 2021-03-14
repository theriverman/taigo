package main

import (
	"testing"

	"github.com/theriverman/taigo"
)

func TestProjects(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// List Projects
	projectsList, err := Client.Project.List(nil)
	if err != nil {
		t.Error(err)
	}
	projects, err := projectsList.AsProjects() // Convert the initial ProjectsLIST into a more generic Project
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Total Projects: %d", len(projects))
	}

	// Create new Project
	newProj, err := Client.Project.Create(&taigo.Project{
		Name:        "A new test project",
		Description: "A project for testing purposes only",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Created Project's Name: %s", newProj.Name)
	}

	// Edit new Project
	newProjDescTextEdit := newProj.Description + "_some more text after edit"
	newProj.Description = newProjDescTextEdit
	newProjEdited, err := Client.Project.Edit(newProj)
	if err != nil {
		t.Error(err)
	}
	if newProjEdited.Description != newProjDescTextEdit {
		t.Errorf("got %s, want %s", newProjEdited.Description, newProjDescTextEdit)
	}

	// Delete new Project
	_, err = Client.Project.Delete(newProj.ID)
	if err != nil {
		t.Error(err)
	}

	// Get Project by slug
	projectBySlug, err := Client.Project.GetBySlug(testProjSlug)
	if err != nil {
		t.Error(err)
	}
	if projectBySlug.ID != testProjID {
		t.Errorf("got %d, want %d", projectBySlug.ID, testProjID)
	}

	// Get Project by ID
	projectByID, err := Client.Project.Get(testProjID)
	if err != nil {
		t.Error(err)
	}
	if projectByID.ID != testProjID {
		t.Errorf("got %d, want %d", projectByID.ID, testProjID)
	}

}
