package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func (r *UserRepo) UserSetIsActiveByID(ctx context.Context, id int, active bool) (*entity.User, error) {
	query := sq.Update("users").Set("is_active", active).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
