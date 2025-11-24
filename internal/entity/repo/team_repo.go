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

	for _, member := range team.Members {
		query := sq.Update("users").Set("team_name", team.Name).Where(sq.Eq{"id": member.Id}).
			Suffix("RETURNING id, pull_request_name, author_id, status, assigned_reviewers").
			PlaceholderFormat(sq.Dollar)
	}

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

	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &pr, nil
}
