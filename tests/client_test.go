package main

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	var makeurltests = []struct {
		in  []string
		out string
	}{
		{[]string{"epics"}, fmt.Sprintf("%s/api/v1/epics", testHostURL)},
		{[]string{"epics", "5"}, fmt.Sprintf("%s/api/v1/epics/5", testHostURL)},
		{[]string{"epics", "bulk_create"}, fmt.Sprintf("%s/api/v1/epics/bulk_create", testHostURL)},
		{[]string{"epics", "attachments", "5"}, fmt.Sprintf("%s/api/v1/epics/attachments/5", testHostURL)},
	}

	for _, tt := range makeurltests {
		s := Client.MakeURL(tt.in...)
		if s != tt.out {
			t.Errorf("got %s, want %s", s, tt.out)
		}
	}

}
