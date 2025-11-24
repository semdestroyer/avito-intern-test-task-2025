package dto

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequestDTO struct {
	Id                string
	PullRequestId     string `json: ""`
	PullRequestName   string `json: ""`
	AuthorId          string `json: ""`
	Status            Status `json: ""`
	AssignedReviewers []UserDTO
}

type PullRequestReassignDTO struct {
	OldUserId int
	NewUserId int
	PrId      int
}
