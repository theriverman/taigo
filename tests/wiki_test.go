package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestWiki(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	randomSlug := "wiki-test-" + strings.ToLower(RandStringBytesMaskImprSrcUnsafe(8))
	created, err := Client.Wiki.Create(&taiga.WikiPage{
		Project: testProjID,
		Slug:    randomSlug,
		Content: "# Wiki page from integration tests",
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	pages, err := Client.Wiki.List(&taiga.WikiQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	}
	if len(pages) == 0 {
		t.Errorf("expected at least one wiki page")
	}

	// Get
	pageByID, err := Client.Wiki.Get(created.ID)
	if err != nil {
		t.Error(err)
	}
	if pageByID.ID != created.ID {
		t.Errorf("got %d, want %d", pageByID.ID, created.ID)
	}

	// GetBySlug
	pageBySlug, err := Client.Wiki.GetBySlug(randomSlug, testProjID)
	if err != nil {
		t.Error(err)
	}
	if pageBySlug.ID != created.ID {
		t.Errorf("got %d, want %d", pageBySlug.ID, created.ID)
	}

	// Edit
	created.Content = "# Updated wiki page from integration tests"
	edited, err := Client.Wiki.Edit(created)
	if err != nil {
		t.Error(err)
	}
	if edited.Content != "# Updated wiki page from integration tests" {
		t.Errorf("got %q, want %q", edited.Content, "# Updated wiki page from integration tests")
	}

	// Render
	renderedHTML, err := Client.Wiki.Render("**hello wiki**", testProjID)
	if err != nil {
		t.Error(err)
	}
	if renderedHTML == "" {
		t.Errorf("expected rendered html to be non-empty")
	}

	// Attachment
	testFileName := "initial_test_data.json"
	attachment := &taiga.Attachment{Name: testFileName, Description: "Wiki attachment from integration test"}
	attachment.SetFilePath(fmt.Sprintf("%s%s%s", cwd, string(os.PathSeparator), testFileName))
	createdAttachment, err := Client.Wiki.CreateAttachment(attachment, created)
	if err != nil {
		t.Error(err)
	}
	if createdAttachment.Name != testFileName {
		t.Errorf("got %q, want %q", createdAttachment.Name, testFileName)
	}

	// Delete
	if _, err := Client.Wiki.Delete(created.ID); err != nil {
		t.Error(err)
	}
}
