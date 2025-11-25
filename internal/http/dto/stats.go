package dto

type AssignmentStatsDTO struct {
	UserAssignments        []UserAssignmentDTO        `json:"user_assignments"`
	PullRequestAssignments []PullRequestAssignmentDTO `json:"pull_request_assignments"`
}

type UserAssignmentDTO struct {
	UserId      string `json:"user_id"`
	Assignments int    `json:"assignments"`
}

type PullRequestAssignmentDTO struct {
	PullRequestId string `json:"pull_request_id"`
	Reviewers     int    `json:"reviewers"`
}
