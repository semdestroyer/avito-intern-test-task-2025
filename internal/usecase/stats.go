package usecase

import (
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"context"
	"time"
)

type StatsUsecase struct {
	prRepo *repo.PrRepo
}

func NewStatsUsecase(pr *repo.PrRepo) *StatsUsecase {
	return &StatsUsecase{
		prRepo: pr,
	}
}

func (su StatsUsecase) GetAssignmentStats() (dto.AssignmentStatsDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCounts, err := su.prRepo.GetUserAssignmentCounts(ctx)
	if err != nil {
		return dto.AssignmentStatsDTO{}, err
	}

	prCounts, err := su.prRepo.GetPullRequestAssignmentCounts(ctx)
	if err != nil {
		return dto.AssignmentStatsDTO{}, err
	}

	userAssignments := make([]dto.UserAssignmentDTO, 0, len(userCounts))
	for _, c := range userCounts {
		userAssignments = append(userAssignments, dto.UserAssignmentDTO{
			UserId:      c.ID,
			Assignments: c.Count,
		})
	}

	prAssignments := make([]dto.PullRequestAssignmentDTO, 0, len(prCounts))
	for _, c := range prCounts {
		prAssignments = append(prAssignments, dto.PullRequestAssignmentDTO{
			PullRequestId: c.ID,
			Reviewers:     c.Count,
		})
	}

	return dto.AssignmentStatsDTO{
		UserAssignments:        userAssignments,
		PullRequestAssignments: prAssignments,
	}, nil
}
