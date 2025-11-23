package queries

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequestQuery struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            Status
	AssignedReviewers []User
}
