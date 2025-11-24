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

func (r *PrRepo) GetPRsByUser(ctx context.Context, id int) ([]entity.PullRequest, error) {

	query := sq.Select("pr.id, pr.author_id, pr.pull_request_name, s.status").
		From("pull_reqests pr").Join("statuses s ON pr.status_id = s.id").
		Where(sq.Eq{"author_id": id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql, args)
	var prs []entity.PullRequest

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	for rows.Next() {
		var pr entity.PullRequest
		err = rows.Scan(
			&pr.Id,
			&pr.AuthorId,
			&pr.PullRequestName,
			&pr.Status,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	return prs, nil
}
func (r *PrRepo) MarkAsMerged(ctx context.Context, id int) (*entity.PullRequest, error) {
	query := sq.Update("pull_requests").Set("status", entity.MERGED).Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, pull_request_name, author_id, status, assigned_reviewers").
		PlaceholderFormat(sq.Dollar)
	//TODO: подумать над timestamp merged_at
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(
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

func (r *PrRepo) Reassign(ctx context.Context, oldUserId int, newUserId, prId int) (*entity.PullRequest, error) {
	query := sq.Update("assigned_reviewers").Set("reviewer_id", newUserId).Where(sq.Eq{"pull_request_id": prId, "reviewer_id": oldUserId}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	query = sq.Select("id, pull_request_name, author_id, status_id").From("pull_requests").Where(sq.Eq{"id": prId}).PlaceholderFormat(sq.Dollar)
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	query = sq.Select("id, pull_request_name, author_id, status_id").From("pull_requests").Where(sq.Eq{"id": prId}).PlaceholderFormat(sq.Dollar)
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan()

	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &pr, nil
}

func (r *PrRepo) Create(ctx context.Context, pr *entity.PullRequest) (*entity.PullRequest, error) {
	query := sq.Insert().Into("pull_requests")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var pr entity.PullRequest
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan()

	if err != nil {
		if err == pgx.ErrNoRows {
			//return nil //,ErrUserNotFound
		}
		return nil, err
	}

	return &pr, nil
}
