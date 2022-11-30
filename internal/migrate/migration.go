package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(migrationDir, dsn string) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", migrationDir),
		"postgres", driver)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.Up()) // or m.Step(2) if you want to explicitly set the number of migrations to run
}
