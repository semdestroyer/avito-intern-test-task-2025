package usecase

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"context"
	"log"
	"strconv"
	"time"
)

type PullRequestUsecase struct {
	prRepo   *repo.PrRepo
	userRepo *repo.UserRepo
}

func NewPullRequestUsecase(up *repo.UserRepo, pr *repo.PrRepo) *PullRequestUsecase {
	return &PullRequestUsecase{
		prRepo:   pr,
		userRepo: up,
	}
}

func (pc PullRequestUsecase) Create(prDTO dto.PullRequestDTO) dto.PullRequestDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := strconv.Atoi(prDTO.AuthorId)
	if err != nil {
		log.Fatal(err.Error())
	}
	user, err := pc.userRepo.GetUserById(ctx, id)
	pr := entity.PullRequest{
		PullRequestName: prDTO.PullRequestName,
		AuthorId: entity.User{
			Id:       user.Id,
			Username: user.Username,
			IsActive: user.IsActive,
			TeamName: user.TeamName,
		},
		Status: entity.OPEN, //TODO: уточнить этот момент по спеке типа может ли создаваться со статусом merged
		//	AssignedReviewers: prDTO.AssignedReviewers, //TODO: разобраться с тем как назначаются ревьюеры
	}

	prc, err := pc.prRepo.Create(ctx, &pr)

	if err != nil {
		log.Fatal("user service failed: ", err)
	}

	asr := make([]dto.UserDTO, 0)

	for _, reviewer := range prc.AssignedReviewers {
		revDto := dto.UserDTO{
			UserId:   string(rune(reviewer.Id)),
			Username: reviewer.Username,
			TeamName: reviewer.TeamName,
			IsActive: reviewer.IsActive,
		}
		asr = append(asr, revDto)
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     string(rune(prc.Id)),
		PullRequestName:   prc.PullRequestName,
		AuthorId:          string(rune(prc.AuthorId.Id)), //TODO: посмотреть переименовать authorid т.к. это объект а не айди
		Status:            dto.Status(prc.Status),
		AssignedReviewers: asr,
	}

	return resDto
}

func (pc PullRequestUsecase) Reassign(prDTO dto.PullRequestReassignDTO) dto.PullRequestDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reassign, err := pc.prRepo.Reassign(ctx, prDTO.OldUserId, prDTO.NewUserId, prDTO.PrId)
	if err != nil {
		return dto.PullRequestDTO{}
	}

	if err != nil {
		log.Fatal("user service failed: ", err)
	}

	asr := make([]dto.UserDTO, 0)

	for _, reviewer := range reassign.AssignedReviewers {
		revDto := dto.UserDTO{
			UserId:   string(rune(reviewer.Id)),
			Username: reviewer.Username,
			TeamName: reviewer.TeamName,
			IsActive: reviewer.IsActive,
		}
		asr = append(asr, revDto)
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     string(rune(reassign.Id)),
		PullRequestName:   reassign.PullRequestName,
		AuthorId:          string(rune(reassign.AuthorId.Id)), //TODO: посмотреть переименовать authorid т.к. это объект а не айди
		Status:            dto.Status(reassign.Status),
		AssignedReviewers: asr,
	}

	return resDto
}

func (pc PullRequestUsecase) Merge(prDTO dto.PullRequestDTO) dto.PullRequestDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(prDTO.Id)

	merged, err := pc.prRepo.MarkAsMerged(ctx, id)

	asr := make([]dto.UserDTO, 0)

	for _, reviewer := range merged.AssignedReviewers {
		revDto := dto.UserDTO{
			UserId:   string(rune(reviewer.Id)),
			Username: reviewer.Username,
			TeamName: reviewer.TeamName,
			IsActive: reviewer.IsActive,
		}
		asr = append(asr, revDto)
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     string(rune(merged.Id)),
		PullRequestName:   merged.PullRequestName,
		AuthorId:          string(rune(merged.AuthorId.Id)), //TODO: посмотреть переименовать authorid т.к. это объект а не айди
		Status:            dto.Status(merged.Status),
		AssignedReviewers: asr,
	}

	if err != nil {
		return dto.PullRequestDTO{}
	}

	return resDto

}
