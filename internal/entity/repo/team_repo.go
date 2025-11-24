package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type TeamRepo struct {
	db *db.DB
}

func NewTeamRepo(db *db.DB) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (tr TeamRepo) CreateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	query := sq.Insert("teams(Name)").Values(team.Name)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	err = tr.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&pr.Id,
		&pr.PullRequestName,
		&pr.AuthorId,
		&pr.Status,
		&pr.AssignedReviewers,
	)

	members := make([]entity.User, 0)
	for _, member := range team.Members {
		m := entity.User{}
		queryUsers := sq.Update("users").Set("team_name", team.Name).Where(sq.Eq{"id": member.Id}).
			Suffix("RETURNING id, is_active, username, team_name").
			PlaceholderFormat(sq.Dollar)
		sql, args, err := queryUsers.ToSql()
		err = tr.db.Pool.QueryRow(ctx, sql, args...).Scan(
			&m.Id,
			&m.IsActive,
			&m.Username,
			&m.TeamName,
		)
		if err != nil {
			return nil, err
		}

		members = append(members)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &entity.Team{
		Name:    team.Name,
		Members: members,
	}, nil
}
