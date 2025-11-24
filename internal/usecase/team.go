package usecase

import (
	"avito-intern-test-task-2025/internal/entity"
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

	users, err := tc.userRepo.GetMembersByTeamName(ctx, query.TeamName)

	if err != nil {
		log.Fatal("user service failed: ", err)
	}

	udto := make([]dto.TeamMemberDTO, 0)
	for _, u := range users {
		member := dto.TeamMemberDTO{
			UserId:   string(rune(u.Id)),
			Username: u.Username,
			IsActive: u.IsActive,
		}
		udto = append(udto, member)
	}

	return dto.TeamDTO{
		Name:    query.TeamName,
		Members: udto,
	}
}

func (tc TeamUsecase) AddTeam(teamDTO dto.TeamDTO) dto.TeamDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	members := make([]entity.User, 0)

	for _, member := range teamDTO.Members {

		id, err := strconv.Atoi("member.UserId")
		if err != nil {
			log.Fatal("error during conv")
		}

		user := entity.User{
			Id:       id,
			Username: member.Username,
		}

		members = append(members, user)
	}

	team := entity.Team{
		Name:    teamDTO.Name,
		Members: members,
	}
	t, err := tc.teamRepo.CreateTeam(ctx, &team)

	if err != nil {
		log.Fatal("user service failed: ", err)
	}

	membersDto := make([]dto.TeamMemberDTO, 0)

	for _, member := range t.Members {
		user := dto.TeamMemberDTO{
			UserId:   string(rune(member.Id)),
			Username: member.Username,
			IsActive: member.IsActive,
		}

		membersDto = append(membersDto, user)
	}

	return dto.TeamDTO{
		Name:    t.Name,
		Members: membersDto,
	}
}
