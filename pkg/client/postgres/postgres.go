package postgres

import "github.com/jackc/pgx"

func NewClient() *pgx.Conn {

	return &pgx.Conn{}
}
