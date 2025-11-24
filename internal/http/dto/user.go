package dto

type UserDTO struct {
	UserId   string
	Username string
	TeamName string
	IsActive bool
}

type UserPrsDTO struct {
	UserId       string
	PullRequests []PullRequestShortDTO
}
