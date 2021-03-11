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
	}
	if len(users) != 1 {
		t.Errorf("got %q, want %q", len(users), 1)
	}

	// List Projects
	projectsList, err := Client.Project.List(nil)
	if err != nil {
		t.Error(err)
	}
	projects, err := projectsList.AsProjects() // Convert the initial ProjectsLIST into a more generic Project
	if err != nil {
		t.Error(err)
	}
	if len(projects) != 1 {
		t.Errorf("got %q, want %q", len(users), 1)
	}

	// Get Project by slug
	projectBySlug, err := Client.Project.GetBySlug(testProjSlug)
	if err != nil {
		t.Error(err)
	}
	if projectBySlug.ID != testProjID {
		t.Errorf("got %q, want %q", projectBySlug.ID, testProjID)
	}

	// Get Project by ID
	projectByID, err := Client.Project.Get(testProjID)
	if err != nil {
		t.Error(err)
	}
	if projectByID.ID != testProjID {
		t.Errorf("got %q, want %q", projectByID.ID, testProjID)
	}

}
