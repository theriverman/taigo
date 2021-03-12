package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestProjects(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// List Users
	users, err := Client.User.List(&taiga.UsersQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Total Users: %d", len(users))
	}

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
