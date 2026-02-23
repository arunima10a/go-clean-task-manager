package repo

import (
	"context"
	"fmt"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/pkg/postgres"
)

type UserRepo struct {
	pg *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}

}

func (r *UserRepo) Create(ctx context.Context, u entity.User) (int, error) {
	sql := `INSERT INTO users (email, password) VALUES ($1,$2) RETURNING id`

	var id int

	err := r.pg.GetQueryer(ctx).QueryRow(ctx, sql, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserRepo - Create - Scan ID: %w", err)
	}
	return id, nil

}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sql := `SELECT id, email, password FROM users WHERE email = $1`
	var u entity.User
	err := r.pg.Pool.QueryRow(ctx, sql, email).Scan(&u.ID, &u.Email, &u.Password)
	return u, err
}
