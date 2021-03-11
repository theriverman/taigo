package main

import (
	"fmt"
	"os"
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestIssues(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// Create Issue
	subject := "Test Issue"
	issue, err := Client.Issue.Create(&taiga.Issue{Project: testProjID, Subject: subject})
	if (err != nil) || (issue.Subject != subject) {
		t.Error(err)
		t.FailNow()
	}

	// List Issues
	issues, err := Client.Issue.List(&taiga.IssueQueryParams{})
	if err != nil {
		t.Error(err)
	}
	if len(issues) != 1 {
		t.Errorf("got %q, want %q", len(issues), 1)
	}

	// Edit Issue
	issue.Description = "Added some text here via Client.Issue.Edit()"
	issuePatched, err := Client.Issue.Edit(issue)
	if err != nil {
		t.Error(err)
	}
	if issuePatched.Version != 2 {
		t.Errorf("got %q, want %q", issuePatched.Version, 2)
	}

	// Get Issue
	issueGet, err := Client.Issue.Get(issue.ID)
	if err != nil {
		t.Error(err)
	}
	if issueGet.Subject != subject {
		t.Errorf("got %q, want %q", issueGet.Subject, subject)
	}

	// Create an Issue Attachment
	attachment := &taiga.Attachment{
		Name:        "A random project file",
		Description: "This is a test file uploaded via TAIGO",
	}
	testFileName := "initial_test_data.json"
	attachment.SetFilePath(fmt.Sprintf("%s%s%s", cwd, string(os.PathSeparator), testFileName))
	attachmentDetails, err := Client.Issue.CreateAttachment(attachment, issue)
	if err != nil {
		t.Error(err)
	}

	if attachmentDetails.Name != testFileName {
		t.Errorf("got %q, want %q", attachmentDetails.Name, testFileName)
	}

}
