package main

import (
	"testing"
)

func TestClient(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	var makeurltests = []struct {
		in  []string
		out string
	}{
		{[]string{"epics"}, "http://localhost:9000/api/v1/epics"},
		{[]string{"epics", "5"}, "http://localhost:9000/api/v1/epics/5"},
		{[]string{"epics", "bulk_create"}, "http://localhost:9000/api/v1/epics/bulk_create"},
		{[]string{"epics", "attachments", "5"}, "http://localhost:9000/api/v1/epics/attachments/5"},
	}

	for _, tt := range makeurltests {
		s := Client.MakeURL(tt.in...)
		if s != tt.out {
			t.Errorf("got %q, want %q", s, tt.out)
		}
	}

}
