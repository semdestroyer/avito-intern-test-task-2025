package queries

type PullRequestShortQuery struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          Status
}
