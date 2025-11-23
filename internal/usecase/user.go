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
}

func (uc UserUsecase) UserSetIsActive(query *queries.UserQuery) dto.UserDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(query.UserId)
	u, err := uc.userRepo.UserSetIsActiveByID(ctx, userId, query.IsActive)
	if err != nil {
		log.Fatal("UserIsActive service failed")
	}
	return dto.UserDTO{
		UserId:   strconv.Itoa(u.Id),
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
