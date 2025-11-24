package dto

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequestDTO struct {
	PullRequestId     string `json: ""`
	PullRequestName   string `json: ""`
	AuthorId          string `json: ""`
	Status            Status `json: ""`
	AssignedReviewers []UserDTO
}
