package usecase

import (
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/queries"
	"context"
	"log"
	"strconv"
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

func (uc UserUsecase) UserSetIsActive(query *queries.UserIsActiveQuery) dto.UserDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(query.UserId)
	u, err := uc.userRepo.UserSetIsActiveByID(ctx, userId, query.IsActive)

	if err != nil {
		log.Fatal("user service failed: ", err)
	}

	return dto.UserDTO{
		UserId:   strconv.Itoa(u.Id),
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func (uc UserUsecase) UserGetReviews(query *queries.UserIdQuery) dto.UserPrsDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(query.UserId)
	prs, err := uc.prRepo.GetPRsByUser(ctx, userId)

	if err != nil {
		log.Fatal("user service failed: ", err)
	}
	prDtos := make([]dto.PullRequestShortDTO, 0)
	for _, pr := range prs {
		prDtos = append(prDtos, dto.PullRequestShortDTO{
			PullRequestId:   string(rune(pr.Id)),
			AuthorId:        string(rune(pr.AuthorId.Id)),
			Status:          dto.Status(pr.Status),
			PullRequestName: pr.PullRequestName,
		})
	}
	return dto.UserPrsDTO{
		UserId:       query.UserId,
		PullRequests: prDtos,
	}
}
