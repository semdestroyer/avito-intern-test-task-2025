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
	query := sq.Update("users").Set("is_active", active).Where(sq.Eq{"id": id}).Suffix("RETURNING id, is_active, team_name, username").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.IsActive,
		&user.TeamName,
		&user.Username,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetMembersByTeamName(ctx context.Context, teamname string) ([]entity.User, error) {
	query := sq.Select("id, username, is_active").Where(sq.Eq{"teamname": teamname}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var members []entity.User
	rows, err := r.db.Pool.Query(ctx, sql, args...)

	for rows.Next() {
		var member entity.User
		err := rows.Scan(
			&member.Id,
			&member.Username,
			&member.IsActive,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return members, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	query := sq.Select("id, is_active, team_name, username").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()

	var user entity.User
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.IsActive,
		&user.TeamName,
		&user.Username,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
