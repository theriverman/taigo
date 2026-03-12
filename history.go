package taigo

import (
	"fmt"
	"strconv"
)

// HistoryEntry is a raw DTO for /history endpoints.
type HistoryEntry = RawResource

// HistoryService is a handle to actions related to object history.
type HistoryService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type historyDeleteCommentQueryParams struct {
	ID string `url:"id"`
}

// UserStory -> https://docs.taiga.io/api.html#history-user-story-task-issue-or-wiki
func (s *HistoryService) UserStory(userStoryID int) ([]HistoryEntry, error) {
	if err := requirePositiveID("userStoryID", userStoryID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPath(s.client, s.Endpoint, "userstory", strconv.Itoa(userStoryID))
}

// Task -> https://docs.taiga.io/api.html#history-user-story-task-issue-or-wiki
func (s *HistoryService) Task(taskID int) ([]HistoryEntry, error) {
	if err := requirePositiveID("taskID", taskID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPath(s.client, s.Endpoint, "task", strconv.Itoa(taskID))
}

// Issue -> https://docs.taiga.io/api.html#history-user-story-task-issue-or-wiki
func (s *HistoryService) Issue(issueID int) ([]HistoryEntry, error) {
	if err := requirePositiveID("issueID", issueID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPath(s.client, s.Endpoint, "issue", strconv.Itoa(issueID))
}

// Wiki -> https://docs.taiga.io/api.html#history-user-story-task-issue-or-wiki
func (s *HistoryService) Wiki(wikiID int) ([]HistoryEntry, error) {
	if err := requirePositiveID("wikiID", wikiID); err != nil {
		return nil, err
	}
	return getRawResourceListAtPath(s.client, s.Endpoint, "wiki", strconv.Itoa(wikiID))
}

// DeleteUserStoryComment -> https://docs.taiga.io/api.html#history-delete-comment
func (s *HistoryService) DeleteUserStoryComment(userStoryID int, commentID any) (*RawResource, error) {
	return s.deleteComment("userstory", userStoryID, commentID)
}

// DeleteTaskComment -> https://docs.taiga.io/api.html#history-delete-comment
func (s *HistoryService) DeleteTaskComment(taskID int, commentID any) (*RawResource, error) {
	return s.deleteComment("task", taskID, commentID)
}

// DeleteIssueComment -> https://docs.taiga.io/api.html#history-delete-comment
func (s *HistoryService) DeleteIssueComment(issueID int, commentID any) (*RawResource, error) {
	return s.deleteComment("issue", issueID, commentID)
}

// DeleteWikiComment -> https://docs.taiga.io/api.html#history-delete-comment
func (s *HistoryService) DeleteWikiComment(wikiID int, commentID any) (*RawResource, error) {
	return s.deleteComment("wiki", wikiID, commentID)
}

func (s *HistoryService) deleteComment(resourceType string, resourceID int, commentID any) (*RawResource, error) {
	if err := requirePositiveID("resourceID", resourceID); err != nil {
		return nil, err
	}
	if commentID == nil {
		return nil, fmt.Errorf("commentID must not be nil")
	}
	commentIDValue := fmt.Sprint(commentID)
	if commentIDValue == "" {
		return nil, fmt.Errorf("commentID is required")
	}
	params := &historyDeleteCommentQueryParams{ID: commentIDValue}
	return postRawResourceAtPathWithQuery(
		s.client,
		nil,
		params,
		s.Endpoint,
		resourceType,
		strconv.Itoa(resourceID),
		"delete_comment",
	)
}
