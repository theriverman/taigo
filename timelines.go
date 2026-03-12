package taigo

import (
	"errors"
	"strconv"
)

// TimelineEntry is a raw DTO for /timeline endpoints.
type TimelineEntry = RawResource

// TimelineQueryParams holds query parameters for timeline endpoints.
type TimelineQueryParams struct {
	Page int `url:"page,omitempty"`
}

// TimelineService is a handle to timeline actions.
type TimelineService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// User -> https://docs.taiga.io/api.html#timelines-user-line
func (s *TimelineService) User(userID int, queryParams *TimelineQueryParams) ([]TimelineEntry, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPathWithQuery(s.client, queryParams, s.Endpoint, "user", strconv.Itoa(userID))
}

// Project -> https://docs.taiga.io/api.html#timelines-project-line
func (s *TimelineService) Project(projectID int, queryParams *TimelineQueryParams) ([]TimelineEntry, error) {
	switch {
	case projectID != 0:
	case s.defaultProjectID != 0:
		projectID = s.defaultProjectID
	default:
		return nil, errors.New("projectID is required")
	}
	return getRawResourceListAtPathWithQuery(s.client, queryParams, s.Endpoint, "project", strconv.Itoa(projectID))
}

// Profile -> https://docs.taiga.io/api.html#timelines-user-profile-line
func (s *TimelineService) Profile(userID int, queryParams *TimelineQueryParams) ([]TimelineEntry, error) {
	if err := requirePositiveID("userID", userID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPathWithQuery(s.client, queryParams, s.Endpoint, "profile", strconv.Itoa(userID))
}
