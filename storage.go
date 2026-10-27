package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errNotFound = errors.New("not found")

type Task struct {
	id          int64
	title       string
	description string
	status      string
	assigneeID  *int64
	createdAt   time.Time
	updatedAt   time.Time
}

type TaskStorage interface {
	CreateTask(ctx context.Context, title, description, status string, assigneeID *int64) (Task, error)
	GetTask(ctx context.Context, id int64) (Task, error)
	UpdateTask(ctx context.Context, id int64, title, description, status string, assigneeID *int64) (Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

type PostgresTaskStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresTaskStorage(pool *pgxpool.Pool) *PostgresTaskStorage {
	return &PostgresTaskStorage{pool: pool}
}

func (s *PostgresTaskStorage) CreateTask(ctx context.Context, title, description, status string, assigneeID *int64) (Task, error) {
	var t Task
	err := s.pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, assignee_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, title, description, status, assignee_id, created_at, updated_at`,
		title, description, status, assigneeID,
	).Scan(&t.id, &t.title, &t.description, &t.status, &t.assigneeID, &t.createdAt, &t.updatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("insert task: %w", err)
	}
	return t, nil
}

func (s *PostgresTaskStorage) GetTask(ctx context.Context, id int64) (Task, error) {
	var t Task
	err := s.pool.QueryRow(ctx,
		`SELECT id, title, description, status, assignee_id, created_at, updated_at
		 FROM tasks
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&t.id, &t.title, &t.description, &t.status, &t.assigneeID, &t.createdAt, &t.updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, errNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("select task: %w", err)
	}
	return t, nil
}

func (s *PostgresTaskStorage) UpdateTask(ctx context.Context, id int64, title, description, status string, assigneeID *int64) (Task, error) {
	var t Task
	err := s.pool.QueryRow(ctx,
		`UPDATE tasks
		 SET title = $2,
		     description = $3,
		     status = $4,
		     assignee_id = $5,
		     updated_at = now()
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, title, description, status, assignee_id, created_at, updated_at`,
		id, title, description, status, assigneeID,
	).Scan(&t.id, &t.title, &t.description, &t.status, &t.assigneeID, &t.createdAt, &t.updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, errNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	return t, nil
}

func (s *PostgresTaskStorage) DeleteTask(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE tasks
		 SET deleted_at = now(),
		     updated_at = now()
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errNotFound
	}
	return nil
}
