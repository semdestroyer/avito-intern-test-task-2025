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

	// If team exists and no members provided, return error
	if teamExists && len(team.Members) == 0 {
		return nil, fmt.Errorf("team already exists")
	}

	// If team exists, check if all members are already in this team
	if teamExists && len(team.Members) > 0 {
		userIds := getUserIds(team.Members)
		// Build query with IN clause for array of user IDs
		query := sq.Select("user_id").
			From("users").
			Where(sq.Eq{"team_name": team.Name}).
			Where(sq.Eq{"user_id": userIds}).
			PlaceholderFormat(sq.Dollar)
		
		sql, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}

		existingMembers, err := tr.db.Pool.Query(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		defer existingMembers.Close()

		existingUserIds := make(map[string]bool)
		for existingMembers.Next() {
			var userId string
			if err := existingMembers.Scan(&userId); err != nil {
				return nil, err
			}
			existingUserIds[userId] = true
		}

		// Check if all requested members are already in the team
		allMembersInTeam := true
		for _, member := range team.Members {
			if !existingUserIds[member.Id] {
				allMembersInTeam = false
				break
			}
		}

		if allMembersInTeam {
			return nil, fmt.Errorf("team already exists")
		}
	}

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

// getUserIds extracts user IDs from a slice of users
func getUserIds(users []entity.User) []string {
	ids := make([]string, len(users))
	for i, user := range users {
		ids[i] = user.Id
	}
	return ids
}
