package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	db *db.DB
}

func NewUserRepo(db *db.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) UserSetIsActiveByID(ctx context.Context, id int, active bool) (*entity.User, error) {
	query := sq.Update("users").Set("is_active", active).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
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
