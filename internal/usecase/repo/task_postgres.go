package repo

import (
	"context"
	"fmt"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/pkg/logger"
	"github.com/arunima10a/task-manager/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type TaskRepo struct {
	pg *postgres.Postgres
	l  *logger.Logger
}

func New(pg *postgres.Postgres, l *logger.Logger) *TaskRepo {
	return &TaskRepo{
		pg: pg,
		l:  l,
	}
}

func (r *TaskRepo) Store(ctx context.Context, t entity.Task, userID int) error {
	sql := `INSERT INTO tasks (title, description, status, user_id) VALUES ($1, $2, $3, $4)`

	_, err := r.pg.Pool.Exec(ctx, sql, t.Title, t.Description, t.Status, userID)
	if err != nil {
		r.l.Error(err, "TaskRepo - Store - r.pg.Pool.Exec")
		return fmt.Errorf("database error: %w", err)
	}
	fmt.Println("Repo: Insert successful")
	return nil
}

func (r *TaskRepo) GetAll(ctx context.Context, userID int, status string, limit, offset int) ([]entity.Task, error) {

	var rows pgx.Rows
	var err error

	if status != "" {
		sql := `
			SELECT 
                t.id, t.title, t.description, t.status, 
                COALESCE(t.category_id, 0), -- Fix: If NULL, return 0
                COALESCE(c.name, 'Uncategorized')
			FROM tasks t
			LEFT JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1 AND t.status = $2
			LIMIT $3 OFFSET $4`

		rows, err = r.pg.GetQueryer(ctx).Query(ctx, sql, userID, status, limit, offset)
	} else {
		sql := `
			SELECT 
                t.id, t.title, t.description, t.status, 
                COALESCE(t.category_id, 0), -- Fix: If NULL, return 0
                COALESCE(c.name, 'Uncategorized')
			FROM tasks t
			LEFT JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1
			LIMIT $2 OFFSET $3`

		rows, err = r.pg.GetQueryer(ctx).Query(ctx, sql, userID, limit, offset)
	}
	if err != nil {
		r.l.Error(err, "TaskRepo - GetAll - Query")
		return nil, fmt.Errorf("TaskRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)

	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CategoryID, &t.CategoryName)
		if err != nil {
			return nil, fmt.Errorf("TaskRepo - GetAll - Scan: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}
func (r *TaskRepo) Update(ctx context.Context, t entity.Task, userID int) error {
	var sql string
	var err error

	if userID == 0 {
		sql = `UPDATE tasks SET description = $1 WHERE id = $2`
		_, err = r.pg.GetQueryer(ctx).Exec(ctx, sql, t.Description, t.ID)
	} else {
		sql = `UPDATE tasks SET status = $1 WHERE id = $2 AND user_id = $3`
		_, err = r.pg.GetQueryer(ctx).Exec(ctx, sql, t.Status, t.ID, userID)
	}

	if err != nil {
		r.l.Error(err, "TaskRepo - Update - Exec")
		return fmt.Errorf("TaskRepo - Update - Exec: %w", err)
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM tasks WHERE id = $1`

	res, err := r.pg.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("TaskRepo - Delete: task not found")
	}

	return nil
}

func (r *TaskRepo) GetFiltered(ctx context.Context, userID int, status string) ([]entity.Task, error) {
	sql := `SELECT id, title, description, status FROM tasks WHERE user_id = $1`

	var rows pgx.Rows
	var err error

	if status != "" {
		sql += " AND status = $2"
		rows, err = r.pg.Pool.Query(ctx, sql, userID, status)
	} else {
		rows, err = r.pg.Pool.Query(ctx, sql, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("TaskRepo - GetFiltered - Query: %w", err)
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)
	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("TaskRepo - GetFiltered - Scan: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
