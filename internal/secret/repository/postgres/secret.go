package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

type Secret struct {
	//int
	ID        int
	UserID    int
	Key       string
	Value     string
	CreatedAt string
}
type SecretRepository struct {
	db *pgxpool.Pool
}

// func NewUserRepository()
var _ secret.Repository = (*SecretRepository)(nil)

func NewSecretRepository(db *pgxpool.Pool) *SecretRepository {
	return &SecretRepository{
		db: db,
	}
}
func (r *SecretRepository) CreateSecret(ctx context.Context, user *models.User, secret *models.Secret) error {
	// _, err := r.GetSecret(ctx, user, sc.ID)
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

	insrtStmt, err := tx.Prepare(ctx, "secret.insert", `INSERT INTO secrets (users_id, key, value, type, created_at) VALUES ($1, $2, $3, $4,$5);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insrtStmt.Name, user.ID, secret.Key, secret.Value, secret.Type, time.Now())
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)

	return err

}

func (r *SecretRepository) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	var id int
	var value string
	var secretType string
	err := r.db.QueryRow(ctx, `SELECT id, type, value FROM secrets WHERE users_id = $1 AND key = $2;`, user.ID, key).Scan(&id, &secretType, &value)
	if err != nil {
		return nil, err
	}
	return &models.Secret{
		ID:     id,
		UserID: user.ID,
		Type:   secretType,
		Key:    key,
		Value:  value,
	}, nil
}

func (r *SecretRepository) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	result := make([]*models.Secret, 0)
	rows, err := r.db.Query(ctx, `SELECT id, type, key, value FROM secrets WHERE users_id = $1;`, user.ID)
	if err != nil {
		err = fmt.Errorf("queryRow failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var secretType string
		var key string
		var value string
		err := rows.Scan(&id, &secretType, &key, &value)
		if err != nil {
			return nil, err
		}
		result = append(result, &models.Secret{ID: id, UserID: user.ID, Type: secretType, Key: key, Value: value})
	}
	if rows.Err() != nil {
		return nil, err
	}
	return result, nil
}
func (r *SecretRepository) DeleteSecret(ctx context.Context, user *models.User, key string) error {
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

	deleteStmt, err := tx.Prepare(ctx, "secret.delete", `DELETE FROM secrets WHERE users_id = $1 AND key = $2;`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, deleteStmt.Name, user.ID, key)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return err
}
