package entity

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequest struct {
	Id                int
	PullRequestName   string
	AuthorId          User
	Status            Status
	AssignedReviewers []User
}
