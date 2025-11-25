package dto

import "avito-intern-test-task-2025/internal/entity"

func ConvertStatusToString(status entity.Status) string {
	switch status {
	case entity.OPEN:
		return "OPEN"
	case entity.MERGED:
		return "MERGED"
	default:
		return "UNKNOWN"
	}
}

type PullRequestDTO struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	MergedAt          string   `json:"mergedAt,omitempty"`
}

type PullRequestReassignDTO struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}
