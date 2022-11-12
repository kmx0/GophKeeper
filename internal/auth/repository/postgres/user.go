package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID       string
	Login    string
	Password string
}
type UserRepository struct {
	db  *pgxpool.Pool
}

func NewUserRepository()