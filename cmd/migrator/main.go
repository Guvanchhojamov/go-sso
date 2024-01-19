package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
)

const (
	emptyValInt = 0
)

func main() {
	var databasePath, migrationsPath, migrationsTable, command string
	flag.StringVar(&databasePath, "database-path", "./database/sso.db", "path to database folder")
	flag.StringVar(&migrationsPath, "migrations-path", "./migrations", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "schema_migrations", "name of migrations table")
	flag.StringVar(&command, "command", "run", "migration command up|down")

	flag.Parse()

	if len(databasePath) == emptyValInt {
		log.Fatal("database-path required")
	}
	if len(migrationsPath) == emptyValInt {
		log.Fatal("migrations-path required")
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", databasePath, migrationsTable),
	)
	if err != nil {
		log.Fatal(err)
	}
	// command for migrate up|down:
	// go run ./cmd/migrator --command=up
	// go run ./cmd/migrator --command=down

	switch command {
	case "up":
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("migrations up to date")
			}
			panic(err)
		}
		fmt.Println("migrations up success")
	case "down":
		if err = m.Down(); err != nil {
			fmt.Println("migrations down error")
			panic(err)
		}
		fmt.Println("migrations down success")
	default:
		fmt.Println("please enter flag --command=up|down")
	}

}
