package main

import (
	"testing"
)

func TestStats(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	discoverStats, err := Client.Stats.GetDiscoverStats()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Projects.Total: %d", discoverStats.Projects.Total)
	}

	// This endpoint is disabled in the default Taiga deployment
	// systemStats, err := Client.Stats.GetSystemStats()
	// if err != nil {
	// 	t.Error(err)
	// } else {
	// 	t.Logf("Projects.AverageLastFiveWorkingDays: %f", systemStats.Projects.AverageLastFiveWorkingDays)
	// 	t.Logf("Projects.AverageLastSevenDays: %f", systemStats.Projects.AverageLastSevenDays)
	// 	t.Logf("Projects.Today: %d", systemStats.Projects.Today)
	// }
}
