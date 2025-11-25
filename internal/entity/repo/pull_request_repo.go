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
	db       *db.DB
	userRepo *UserRepo
}

func NewPrRepo(db *db.DB, userRepo *UserRepo) *PrRepo {
	return &PrRepo{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *PrRepo) GetPRsByUser(ctx context.Context, userId string) ([]entity.PullRequest, error) {
	// First verify user exists
	_, err := r.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", userId)
	}

	// Get PRs where user is assigned as reviewer (not author)
	query := sq.Select("pr.id, pr.author_id, pr.pull_request_name, s.status").
		From("pull_requests pr").
		Join("statuses s ON pr.status_id = s.id").
		Join("assigned_reviewers ar ON pr.id = ar.pull_request_id").
		Where(sq.Eq{"ar.reviewer_id": userId}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var prs []entity.PullRequest
	for rows.Next() {
		var pr entity.PullRequest
		var statusStr string
		var prId string
		var authorUserId string
		err = rows.Scan(
			&prId,
			&authorUserId,
			&pr.PullRequestName,
			&statusStr,
		)
		if err != nil {
			return nil, err
		}

		pr.Id = prId

		// Load author user
		author, err := r.userRepo.GetUserById(ctx, authorUserId)
		if err != nil {
			return nil, fmt.Errorf("failed to load author: %w", err)
		}
		pr.AuthorId = *author

		// Convert status string to enum
		if statusStr == "OPENED" {
			pr.Status = entity.OPEN
		} else if statusStr == "MERGED" {
			pr.Status = entity.MERGED
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PrRepo) MarkAsMerged(ctx context.Context, id string) (*entity.PullRequest, error) {
	query := sq.Update("pull_requests").Set("status_id", 2).Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, pull_request_name, author_id, status_id").
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	var statusId int
	var prId string
	var authorUserId string
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&prId,
		&pr.PullRequestName,
		&authorUserId,
		&statusId,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("PR not found: %s", id)
		}
		return nil, err
	}

	pr.Id = prId

	// Load author user
	author, err := r.userRepo.GetUserById(ctx, authorUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}
	pr.AuthorId = *author

	// Set status based on ID
	if statusId == 2 {
		pr.Status = entity.MERGED
	} else {
		pr.Status = entity.OPEN
	}

	// Load assigned reviewers
	reviewers, err := r.GetAssignedReviewers(ctx, pr.Id)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PrRepo) Reassign(ctx context.Context, oldUserId string, newUserId string, prId string) (*entity.PullRequest, error) {
	// First update the assigned reviewer
	query := sq.Update("assigned_reviewers").Set("reviewer_id", newUserId).Where(sq.Eq{"pull_request_id": prId, "reviewer_id": oldUserId}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	// Then get the updated PR
	return r.GetByID(ctx, prId)
}

func (r *PrRepo) GetByID(ctx context.Context, id string) (*entity.PullRequest, error) {
	query := sq.Select("pr.id, pr.pull_request_name, pr.author_id, s.status").
		From("pull_requests pr").Join("statuses s ON pr.status_id = s.id").
		Where(sq.Eq{"pr.id": id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	var statusStr string
	var prId string
	var authorUserId string
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&prId,
		&pr.PullRequestName,
		&authorUserId,
		&statusStr,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("PR not found: %s", id)
		}
		return nil, err
	}

	pr.Id = prId

	// Load author user
	author, err := r.userRepo.GetUserById(ctx, authorUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load author: %w", err)
	}
	pr.AuthorId = *author

	// Convert status string to enum
	if statusStr == "OPENED" {
		pr.Status = entity.OPEN
	} else if statusStr == "MERGED" {
		pr.Status = entity.MERGED
	}

	// Load assigned reviewers
	reviewers, err := r.GetAssignedReviewers(ctx, pr.Id)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PrRepo) GetAssignedReviewers(ctx context.Context, prId string) ([]entity.User, error) {
	query := sq.Select("u.user_id, u.username, u.is_active, u.team_name").
		From("assigned_reviewers ar").Join("users u ON ar.reviewer_id = u.user_id").
		Where(sq.Eq{"ar.pull_request_id": prId}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []entity.User
	for rows.Next() {
		var reviewer entity.User
		err = rows.Scan(
			&reviewer.Id,
			&reviewer.Username,
			&reviewer.IsActive,
			&reviewer.TeamName,
		)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func (r *PrRepo) Update(ctx context.Context, pr *entity.PullRequest) (*entity.PullRequest, error) {
	// Update assigned reviewers
	// First delete existing assignments
	deleteQuery := sq.Delete("assigned_reviewers").Where(sq.Eq{"pull_request_id": pr.Id}).PlaceholderFormat(sq.Dollar)
	deleteSql, deleteArgs, err := deleteQuery.ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.Pool.Exec(ctx, deleteSql, deleteArgs...)
	if err != nil {
		return nil, err
	}

	// Then insert new assignments
	for _, reviewer := range pr.AssignedReviewers {
		insertQuery := sq.Insert("assigned_reviewers").Columns("pull_request_id", "reviewer_id").
			Values(pr.Id, reviewer.Id).
			PlaceholderFormat(sq.Dollar)
		insertSql, insertArgs, err := insertQuery.ToSql()
		if err != nil {
			return nil, err
		}
		_, err = r.db.Pool.Exec(ctx, insertSql, insertArgs...)
		if err != nil {
			return nil, err
		}
	}

	return pr, nil
}

func (r *PrRepo) Create(ctx context.Context, pr *entity.PullRequest) (*entity.PullRequest, error) {
	query := sq.Insert("pull_requests").Columns("id", "pull_request_name", "author_id", "status_id").
		Values(pr.Id, pr.PullRequestName, pr.AuthorId.Id, int(pr.Status)+1).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var id string
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("failed to create PR")
		}
		return nil, err
	}

	createdPr := &entity.PullRequest{
		Id:                id,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: []entity.User{},
	}

	return createdPr, nil
}
