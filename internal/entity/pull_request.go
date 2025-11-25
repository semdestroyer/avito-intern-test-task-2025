package entity

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequest struct {
	Id                string
	PullRequestName   string
	AuthorId          User
	Status            Status
	AssignedReviewers []User
}
