package repo

import (
	"context"
	"fmt"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/pkg/postgres"
)

type CategoryRepo struct {
	pg *postgres.Postgres
}

func NewcategoryRepo(pg *postgres.Postgres) *CategoryRepo {
	return &CategoryRepo{pg}

}

func (r *CategoryRepo) Create(ctx context.Context, name string, userID int) (int, error) {
	sql := `INSERT INTO categories (name, user_id) VALUES ($1,$2) RETURNING id`

	var id int

	err := r.pg.GetQueryer(ctx).QueryRow(ctx, sql, name, userID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CategoryRepo - Create - Scan ID: %w", err)
	}
	return id, nil

}

func (r *CategoryRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sql := `SELECT id, email, password FROM users WHERE email = $1`
	var u entity.User
	err := r.pg.Pool.QueryRow(ctx, sql, email).Scan(&u.ID, &u.Email, &u.Password)
	return u, err
}

func (r *CategoryRepo) GetAll(ctx context.Context, userID int) ([]entity.Category, error) {
	sql := `SELECT id, name, user_id FROM categories WHERE user_id = $1`

	rows, err := r.pg.GetQueryer(ctx).Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	categories := make([]entity.Category, 0)
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.UserID); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
