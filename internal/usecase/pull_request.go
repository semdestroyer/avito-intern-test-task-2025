package usecase

import "avito-intern-test-task-2025/internal/entity/repo"

type PullRequestUsecase struct {
	prRepo *repo.PrRepo
}

func NewPullRequestUsecase(rp *repo.UserRepo, pr *repo.PrRepo) *PullRequestUsecase {
	return &PullRequestUsecase{
		prRepo: pr,
	}
}
