package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
	"context"
	"fmt"

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

func (r *UserRepo) UserSetIsActiveByID(ctx context.Context, id string, active bool) (*entity.User, error) {
	query := sq.Update("users").Set("is_active", active).Where(sq.Eq{"user_id": id}).Suffix("RETURNING user_id, is_active, team_name, username").
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
	// First check if team exists
	teamQuery := sq.Select("id").From("teams").Where(sq.Eq{"name": teamname}).
		PlaceholderFormat(sq.Dollar)

	teamSql, teamArgs, err := teamQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var teamId int
	err = r.db.Pool.QueryRow(ctx, teamSql, teamArgs...).Scan(&teamId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("team not found: %s", teamname)
		}
		return nil, err
	}

	// Then get team members
	query := sq.Select("user_id, username, is_active, team_name").From("users").Where(sq.Eq{"team_name": teamname}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.User
	for rows.Next() {
		var member entity.User
		err := rows.Scan(
			&member.Id,
			&member.Username,
			&member.IsActive,
			&member.TeamName,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *UserRepo) GetActiveMembersByTeamName(ctx context.Context, teamname string) ([]entity.User, error) {
	query := sq.Select("user_id, username, is_active, team_name").From("users").Where(sq.Eq{"team_name": teamname, "is_active": true}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.User
	for rows.Next() {
		var member entity.User
		err := rows.Scan(
			&member.Id,
			&member.Username,
			&member.IsActive,
			&member.TeamName,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	query := sq.Select("user_id, is_active, team_name, username").From("users").Where(sq.Eq{"user_id": id}).
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
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	return &user, nil
}
