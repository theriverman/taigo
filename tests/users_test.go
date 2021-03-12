package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestUsers(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// List Users
	users, err := Client.User.List(&taiga.UsersQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Total Users: %d", len(users))
	}

	// Get /users/me
	me, err := Client.User.Me()
	if err != nil {
		t.Error(err)
	}
	if me.ID != testUserID {
		t.Errorf("got %d, want %d", me.ID, testUserID)
	}

	// Get /users/{{ .testUserID }} and compare the retrieved FullNameDisplay
	adminUser, err := Client.User.Get(testUserID)
	if err != nil {
		t.Error(err)
	}
	if adminUser.FullNameDisplay != "admin" {
		t.Errorf("got %s, want %s", adminUser.FullName, "admin")
	}

	// Patch the retrieved adminUser
	adminUserBioText := "Some text in user's bio"
	adminUser.Bio = adminUserBioText
	adminUser.Email = "" // exclude from payload to avoid "_error_message": "Duplicated email"
	adminUserPatched, err := Client.User.Edit(adminUser)
	if err != nil {
		t.Error(err)
	}
	if adminUserPatched.Bio != adminUserBioText {
		t.Errorf("got %s, want %s", adminUserPatched.Bio, adminUserBioText)
	}
	// Reset to original; check again
	adminUser.Bio = ""
	if adminUserPatched.Bio != adminUserBioText {
		t.Errorf("got %s, want %s", "", "")
	}

}
