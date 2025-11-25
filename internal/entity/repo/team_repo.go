package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
	"context"
	"fmt"
	"strings"

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
	// Check if team already exists
	_, err := tr.GetTeamByName(ctx, team.Name)
	teamExists := err == nil

	// If team doesn't exist, create it
	if !teamExists {
		query := sq.Insert("teams").Columns("name").Values(team.Name).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}

		var teamId int
		err = tr.db.Pool.QueryRow(ctx, sql, args...).Scan(&teamId)
		if err != nil {
			// Check if it's a duplicate key error (race condition)
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				// Team was created by another request, continue
				teamExists = true
			} else {
				return nil, err
			}
		}
	}

	// Create or update users (upsert)
	members := make([]entity.User, 0)
	for _, member := range team.Members {
		// Try to update existing user first
		updateQuery := sq.Update("users").
			Set("team_name", team.Name).
			Set("username", member.Username).
			Set("is_active", member.IsActive).
			Where(sq.Eq{"user_id": member.Id}).
			Suffix("RETURNING user_id, is_active, username, team_name").
			PlaceholderFormat(sq.Dollar)
		
		updateSql, updateArgs, err := updateQuery.ToSql()
		if err != nil {
			return nil, err
		}

		var updatedUser entity.User
		err = tr.db.Pool.QueryRow(ctx, updateSql, updateArgs...).Scan(
			&updatedUser.Id,
			&updatedUser.IsActive,
			&updatedUser.Username,
			&updatedUser.TeamName,
		)

		// If user doesn't exist (pgx.ErrNoRows), create it
		if err != nil {
			if err == pgx.ErrNoRows {
				// User doesn't exist, create it
				insertQuery := sq.Insert("users").
					Columns("user_id", "username", "is_active", "team_name").
					Values(member.Id, member.Username, member.IsActive, team.Name).
					Suffix("RETURNING user_id, is_active, username, team_name").
					PlaceholderFormat(sq.Dollar)
				
				insertSql, insertArgs, err := insertQuery.ToSql()
				if err != nil {
					return nil, err
				}

				err = tr.db.Pool.QueryRow(ctx, insertSql, insertArgs...).Scan(
					&updatedUser.Id,
					&updatedUser.IsActive,
					&updatedUser.Username,
					&updatedUser.TeamName,
				)
				if err != nil {
					return nil, err
				}
			} else {
				// Other error occurred
				return nil, err
			}
		}

		members = append(members, updatedUser)
	}

	return &entity.Team{
		Name:    team.Name,
		Members: members,
	}, nil
}

func (tr TeamRepo) GetTeamByName(ctx context.Context, teamName string) (*entity.Team, error) {
	// First check if team exists
	teamQuery := sq.Select("id").From("teams").Where(sq.Eq{"name": teamName}).
		PlaceholderFormat(sq.Dollar)

	teamSql, teamArgs, err := teamQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var teamId int
	err = tr.db.Pool.QueryRow(ctx, teamSql, teamArgs...).Scan(&teamId)
	if err != nil {
		return nil, fmt.Errorf("team not found: %s", teamName)
	}

	return &entity.Team{
		Name: teamName,
	}, nil
}
