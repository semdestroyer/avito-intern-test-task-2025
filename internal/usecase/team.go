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

type TeamUsecase struct {
	teamRepo *repo.TeamRepo
	userRepo *repo.UserRepo
}

func NewTeamUsecase(tr *repo.TeamRepo, ur *repo.UserRepo) *TeamUsecase {
	return &TeamUsecase{
		teamRepo: tr,
		userRepo: ur,
	}
}

func (tc TeamUsecase) GetTeamMembersByName(query *queries.TeamNameQuery) dto.TeamDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := tc.userRepo.GetMembersByTeamName(ctx, query.TeamName)

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

func (tc TeamUsecase) AddTeam() dto.TeamDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(query.TeamName)
	u, err := tc.userRepo.GetMembersByTeamName()

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
