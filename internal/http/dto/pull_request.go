package dto

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequestDTO struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            Status
	AssignedReviewers []User
}
