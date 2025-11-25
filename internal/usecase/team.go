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
	teamRepo *repo.TeamRepo
	userRepo *repo.UserRepo
}

func NewTeamUsecase(tr *repo.TeamRepo, ur *repo.UserRepo) *TeamUsecase {
	return &TeamUsecase{
		teamRepo: tr,
		userRepo: ur,
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
