package dto

type UserDTO struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserPrsDTO struct {
	UserId       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

type UserIsActiveDTO struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
