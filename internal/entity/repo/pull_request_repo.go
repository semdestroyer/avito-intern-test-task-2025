package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type PrRepo struct {
	db *db.DB
}

func NewPrRepo(db *db.DB) *PrRepo {
	return &PrRepo{
		db: db,
	}
}

func (r *PrRepo) GetPRsByUser(ctx context.Context, id int) (*entity.PullRequest, error) {

	query := sq.Select("pr.id, pr.author_id, pr.pull_request_name, s.status").
		From("pull_reqests pr").Join("statuses s ON pr.status_id = s.id").
		Where(sq.Eq{"author_id": id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql, args)
	var pr entity.PullRequest
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&pr.Id,
		&pr.AuthorId,
		&pr.PullRequestName,
		&pr.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &pr, nil
}
