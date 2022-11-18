package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

// func NewUserRepository()
var _ auth.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// _, err := r.GetUser(ctx, user.ID)
	// if err != nil {
	// 	return err
	// }
	tx, err := r.db.Begin(ctx)
	if err != nil {

		return err
	}
	defer func() {
		if err != nil {
			rerr := tx.Rollback(ctx)
			if rerr != nil {
				err = rerr
			}
		}
	}()

	insrtStmt, err := tx.Prepare(ctx, "insert", `INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insrtStmt.Name, user.Login, user.Password, time.Now())
	if err != nil {
		return err
	}
	return err

}

func (r *UserRepository) GetUser(ctx context.Context, login, password string) (*models.User, error) {

	var id string
	err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE login = $1 AND password = $2;`, login, password).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       id,
		Login:    login,
		Password: password,
	}, nil

}
