package entity

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequest struct {
	PullRequestId       string
	PullRequestquesName string
	AuthorId            string
	Status              Status
	AssignedReviewers   []User
}
