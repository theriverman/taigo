package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestMilestones(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// Create a milestone(sprint)
	randomString := RandStringBytesMaskImprSrcUnsafe(12) // we need unique milestone names
	milestone, err := Client.Milestone.Create(&taiga.Milestone{
		Name:            "A test milestone_" + randomString,
		Project:         testProjID,
		EstimatedStart:  "2020-02-20",
		EstimatedFinish: "2020-02-22",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(milestone.Name)
	}

	// List milestones
	milestones, msTotalInfo, err := Client.Milestone.List(&taiga.MilestonesQueryParams{Project: testProjID})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Total Epics: %d", len(milestones))
		t.Logf("TotalOpenedMilestones: %d", msTotalInfo.TaigaInfoTotalOpenedMilestones)
		t.Logf("TotalClosedMilestones: %d", msTotalInfo.TaigaInfoTotalClosedMilestones)
	}

	// Get a milestone
	ms, err := Client.Milestone.Get(milestone.ID)
	if err != nil {
		t.Error(err)
	}
	if ms.Name != "A test milestone_"+randomString {
		t.Errorf("got %s, want %s", ms.Name, "A test milestone_"+randomString)
	}

	// Edit a milestone
	msPatchInput := ms
	msPatchInput.Name = "A test milestone_" + randomString + "_EDITED"
	msEdited, err := Client.Milestone.Edit(msPatchInput)
	if err != nil {
		t.Error(err)
	}
	if msEdited.Name != "A test milestone_"+randomString+"_EDITED" {
		t.Errorf("got %s, want %s", msEdited.Name, "A test milestone_"+randomString+"_EDITED")
	}

	// Delete Milestone
	_, err = Client.Milestone.Delete(msEdited.ID)
	if err != nil {
		t.Error(err)
	}
}
