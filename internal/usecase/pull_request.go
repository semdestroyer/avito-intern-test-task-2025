package usecase

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
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

func (pc PullRequestUsecase) Create(prDTO dto.PullRequestDTO) (dto.PullRequestDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if prDTO.PullRequestId != "" {
		existingPr, _ := pc.prRepo.GetByID(ctx, prDTO.PullRequestId)
		if existingPr != nil {
			return dto.PullRequestDTO{}, errors.New("PR_EXISTS")
		}
	}

	user, err := pc.userRepo.GetUserById(ctx, prDTO.AuthorId)
	if err != nil {
		return dto.PullRequestDTO{}, errors.New("AUTHOR_NOT_FOUND")
	}

	pr := entity.PullRequest{
		Id:              prDTO.PullRequestId,
		PullRequestName: prDTO.PullRequestName,
		AuthorId: entity.User{
			Id:       user.Id,
			Username: user.Username,
			IsActive: user.IsActive,
			TeamName: user.TeamName,
		},
		Status: entity.OPEN,
	}

	prc, err := pc.prRepo.Create(ctx, &pr)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	reviewers, err := pc.assignReviewers(ctx, user.Id, user.TeamName)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	prc.AssignedReviewers = reviewers
	prc, err = pc.prRepo.Update(ctx, prc)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	reviewerIds := make([]string, len(prc.AssignedReviewers))
	for i, reviewer := range prc.AssignedReviewers {
		reviewerIds[i] = reviewer.Id
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     prc.Id,
		PullRequestName:   prc.PullRequestName,
		AuthorId:          prc.AuthorId.Id,
		Status:            dto.ConvertStatusToString(prc.Status),
		AssignedReviewers: reviewerIds,
	}

	return resDto, nil
}

func (pc PullRequestUsecase) assignReviewers(ctx context.Context, authorId string, teamName string) ([]entity.User, error) {
	candidates, err := pc.userRepo.GetActiveMembersByTeamName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	filteredCandidates := make([]entity.User, 0)
	for _, candidate := range candidates {
		if candidate.Id != authorId {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	reviewers := make([]entity.User, 0)
	maxReviewers := 2
	if len(filteredCandidates) < maxReviewers {
		maxReviewers = len(filteredCandidates)
	}

	// Shuffle candidates and select first maxReviewers
	for i := range filteredCandidates {
		j := rand.Intn(i + 1)
		filteredCandidates[i], filteredCandidates[j] = filteredCandidates[j], filteredCandidates[i]
	}

	for i := 0; i < maxReviewers; i++ {
		reviewers = append(reviewers, filteredCandidates[i])
	}

	return reviewers, nil
}

func (pc PullRequestUsecase) Reassign(prDTO dto.PullRequestReassignDTO) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resDto, replacedBy, err := pc.reassignWithContext(ctx, prDTO)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"pr":          resDto,
		"replaced_by": replacedBy,
	}, nil
}

func (pc PullRequestUsecase) reassignWithContext(ctx context.Context, prDTO dto.PullRequestReassignDTO) (dto.PullRequestDTO, string, error) {
	pr, err := pc.prRepo.GetByID(ctx, prDTO.PullRequestId)
	if err != nil {
		return dto.PullRequestDTO{}, "", errors.New("PR_NOT_FOUND")
	}

	if pr.Status == entity.MERGED {
		return dto.PullRequestDTO{}, "", errors.New("PR_MERGED")
	}

	oldUserId := prDTO.OldUserId
	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer.Id == oldUserId {
			found = true
			break
		}
	}
	if !found {
		return dto.PullRequestDTO{}, "", errors.New("REVIEWER_NOT_ASSIGNED")
	}

	oldReviewer, err := pc.userRepo.GetUserById(ctx, oldUserId)
	if err != nil {
		return dto.PullRequestDTO{}, "", errors.New("REVIEWER_NOT_ASSIGNED")
	}

	candidates, err := pc.userRepo.GetActiveMembersByTeamName(ctx, oldReviewer.TeamName)
	if err != nil {
		return dto.PullRequestDTO{}, "", errors.New("NO_CANDIDATE")
	}

	filteredCandidates := make([]entity.User, 0)
	for _, candidate := range candidates {
		if candidate.Id != pr.AuthorId.Id && candidate.Id != oldUserId {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	if len(filteredCandidates) == 0 {
		return dto.PullRequestDTO{}, "", errors.New("NO_CANDIDATE")
	}

	newReviewer := filteredCandidates[rand.Intn(len(filteredCandidates))]

	updatedReviewers := make([]entity.User, len(pr.AssignedReviewers))
	copy(updatedReviewers, pr.AssignedReviewers)
	replaced := false
	for i, reviewer := range updatedReviewers {
		if reviewer.Id == oldUserId && !replaced {
			updatedReviewers[i] = newReviewer
			replaced = true
			break
		}
	}

	pr.AssignedReviewers = updatedReviewers
	updatedPr, err := pc.prRepo.Update(ctx, pr)
	if err != nil {
		return dto.PullRequestDTO{}, "", err
	}

	reviewerIds := make([]string, len(updatedPr.AssignedReviewers))
	for i, reviewer := range updatedPr.AssignedReviewers {
		reviewerIds[i] = reviewer.Id
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     updatedPr.Id,
		PullRequestName:   updatedPr.PullRequestName,
		AuthorId:          updatedPr.AuthorId.Id,
		Status:            dto.ConvertStatusToString(updatedPr.Status),
		AssignedReviewers: reviewerIds,
	}

	return resDto, newReviewer.Id, nil
}

func (pc PullRequestUsecase) Merge(prDTO dto.PullRequestDTO) (dto.PullRequestDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pr, err := pc.prRepo.GetByID(ctx, prDTO.PullRequestId)
	if err != nil {
		return dto.PullRequestDTO{}, errors.New("PR_NOT_FOUND")
	}

	if pr.Status == entity.MERGED {
		reviewerIds := make([]string, len(pr.AssignedReviewers))
		for i, reviewer := range pr.AssignedReviewers {
			reviewerIds[i] = reviewer.Id
		}

		return dto.PullRequestDTO{
			PullRequestId:     pr.Id,
			PullRequestName:   pr.PullRequestName,
			AuthorId:          pr.AuthorId.Id,
			Status:            dto.ConvertStatusToString(pr.Status),
			AssignedReviewers: reviewerIds,
			MergedAt:          time.Now().Format(time.RFC3339),
		}, nil
	}

	merged, err := pc.prRepo.MarkAsMerged(ctx, pr.Id)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	reviewerIds := make([]string, len(merged.AssignedReviewers))
	for i, reviewer := range merged.AssignedReviewers {
		reviewerIds[i] = reviewer.Id
	}

	resDto := dto.PullRequestDTO{
		PullRequestId:     merged.Id,
		PullRequestName:   merged.PullRequestName,
		AuthorId:          merged.AuthorId.Id,
		Status:            dto.ConvertStatusToString(merged.Status),
		AssignedReviewers: reviewerIds,
		MergedAt:          time.Now().Format(time.RFC3339),
	}

	return resDto, nil
}
