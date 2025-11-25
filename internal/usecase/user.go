package usecase

import (
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/queries"
	"context"
	"time"
)

type UserUsecase struct {
	userRepo *repo.UserRepo
	prRepo   *repo.PrRepo
}

func NewUserUsecase(rp *repo.UserRepo, pr *repo.PrRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: rp,
		prRepo:   pr,
	}
}

func (uc UserUsecase) UserSetIsActive(query *queries.UserIsActiveQuery) (dto.UserDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := uc.userRepo.UserSetIsActiveByID(ctx, query.UserId, query.IsActive)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		UserId:   u.Id,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}, nil
}

func (uc UserUsecase) UserGetReviews(query *queries.UserIdQuery) (dto.UserPrsDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prs, err := uc.prRepo.GetPRsByUser(ctx, query.UserId)
	if err != nil {
		return dto.UserPrsDTO{}, err
	}
	prDtos := make([]dto.PullRequestShortDTO, 0)
	for _, pr := range prs {
		prDtos = append(prDtos, dto.PullRequestShortDTO{
			PullRequestId:   pr.Id,
			AuthorId:        pr.AuthorId.Id,
			Status:          dto.ConvertStatusToString(pr.Status),
			PullRequestName: pr.PullRequestName,
		})
	}
	return dto.UserPrsDTO{
		UserId:       query.UserId,
		PullRequests: prDtos,
	}, nil
}
