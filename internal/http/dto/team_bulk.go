package dto

type TeamBulkDeactivateDTO struct {
	TeamName string   `json:"team_name"`
	UserIds  []string `json:"user_ids"`
}

type BulkReassignmentDTO struct {
	PullRequestId string `json:"pull_request_id"`
	ReplacedUser  string `json:"replaced_user_id"`
	NewReviewer   string `json:"new_reviewer_id"`
}

type BulkDeactivateResultDTO struct {
	TeamName      string                `json:"team_name"`
	Deactivated   []string              `json:"deactivated_user_ids"`
	Reassignments []BulkReassignmentDTO `json:"reassignments"`
}
