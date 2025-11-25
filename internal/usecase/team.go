package usecase

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/queries"
	"context"
	"fmt"
	"time"
)

type TeamUsecase struct {
	teamRepo  *repo.TeamRepo
	userRepo  *repo.UserRepo
	prRepo    *repo.PrRepo
	prUsecase *PullRequestUsecase
}

func NewTeamUsecase(tr *repo.TeamRepo, ur *repo.UserRepo, pr *repo.PrRepo, pc *PullRequestUsecase) *TeamUsecase {
	return &TeamUsecase{
		teamRepo:  tr,
		userRepo:  ur,
		prRepo:    pr,
		prUsecase: pc,
	}
}

func (tc TeamUsecase) GetTeamMembersByName(query *queries.TeamNameQuery) (dto.TeamDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tc.teamRepo.GetTeamByName(ctx, query.TeamName)
	if err != nil {
		return dto.TeamDTO{}, fmt.Errorf("team not found")
	}

	users, err := tc.userRepo.GetMembersByTeamName(ctx, query.TeamName)
	if err != nil {
		return dto.TeamDTO{}, err
	}

	udto := make([]dto.TeamMemberDTO, 0)
	for _, u := range users {
		member := dto.TeamMemberDTO{
			UserId:   u.Id,
			Username: u.Username,
			IsActive: u.IsActive,
		}
		udto = append(udto, member)
	}

	return dto.TeamDTO{
		TeamName: query.TeamName,
		Members:  udto,
	}, nil
}

func (tc TeamUsecase) AddTeam(teamDTO dto.TeamDTO) (dto.TeamDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	members := make([]entity.User, 0)

	for _, member := range teamDTO.Members {
		user := entity.User{
			Id:       member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: teamDTO.TeamName,
		}

		members = append(members, user)
	}

	team := entity.Team{
		Name:    teamDTO.TeamName,
		Members: members,
	}
	t, err := tc.teamRepo.CreateTeam(ctx, &team)
	if err != nil {
		return dto.TeamDTO{}, err
	}

	membersDto := make([]dto.TeamMemberDTO, 0)

	for _, member := range t.Members {
		user := dto.TeamMemberDTO{
			UserId:   member.Id,
			Username: member.Username,
			IsActive: member.IsActive,
		}

		membersDto = append(membersDto, user)
	}

	return dto.TeamDTO{
		TeamName: t.Name,
		Members:  membersDto,
	}, nil
}

func (tc TeamUsecase) BulkDeactivate(teamDTO dto.TeamBulkDeactivateDTO) (dto.BulkDeactivateResultDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if teamDTO.TeamName == "" {
		return dto.BulkDeactivateResultDTO{}, fmt.Errorf("team_name is required")
	}

	if len(teamDTO.UserIds) == 0 {
		return dto.BulkDeactivateResultDTO{}, fmt.Errorf("user_ids is required")
	}

	_, err := tc.teamRepo.GetTeamByName(ctx, teamDTO.TeamName)
	if err != nil {
		return dto.BulkDeactivateResultDTO{}, fmt.Errorf("team not found")
	}

	_, err = tc.userRepo.BulkSetInactiveByTeam(ctx, teamDTO.TeamName, teamDTO.UserIds)
	if err != nil {
		return dto.BulkDeactivateResultDTO{}, err
	}

	reassignments := make([]dto.BulkReassignmentDTO, 0)
	for _, userId := range teamDTO.UserIds {
		prs, err := tc.prRepo.GetPRsByUser(ctx, userId)
		if err != nil {
			continue
		}
		for _, pr := range prs {
			if pr.Status == entity.MERGED {
				continue
			}
			assignDTO := dto.PullRequestReassignDTO{
				PullRequestId: pr.Id,
				OldUserId:     userId,
			}
			_, replacedBy, err := tc.prUsecase.reassignWithContext(ctx, assignDTO)
			if err != nil {
				continue
			}
			reassignments = append(reassignments, dto.BulkReassignmentDTO{
				PullRequestId: pr.Id,
				ReplacedUser:  userId,
				NewReviewer:   replacedBy,
			})
		}
	}

	return dto.BulkDeactivateResultDTO{
		TeamName:      teamDTO.TeamName,
		Deactivated:   teamDTO.UserIds,
		Reassignments: reassignments,
	}, nil
}
