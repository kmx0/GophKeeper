package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

var _ auth.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {

		return fmt.Errorf("error on CreateUser: %w", err)
	}

	defer func() {
		if err != nil {
			rerr := tx.Rollback(ctx)
			if rerr != nil {
				fmt.Printf("error on CreateUser Rollback: %v", rerr)
				err = rerr
			}
		}
	}()

	var id int
	err = r.db.QueryRow(ctx, `SELECT id FROM users WHERE login = $1;`, user.Login).Scan(&id)
	if err == nil {
		return fmt.Errorf("error on CreateUser: %w", auth.ErrLoginBusy)
	}

	insrtStmt, err := tx.Prepare(ctx, "user.insert", `INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3);`)
	if err != nil {
		return fmt.Errorf("error on CreateUser: %w", err)
	}
	_, err = tx.Exec(ctx, insrtStmt.Name, user.Login, user.Password, time.Now())
	if err != nil {

		return fmt.Errorf("error on CreateUser: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {

		return fmt.Errorf("error on CreateUser: %w", err)
	}
	return nil

}

func (r *UserRepository) GetUser(ctx context.Context, login, password string) (*models.User, error) {

	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE login = $1 AND password = $2;`, login, password).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error on GetUser: %w", auth.ErrUserNotFound)
	}
	return &models.User{
		ID:       id,
		Login:    login,
		Password: password,
	}, nil

}
