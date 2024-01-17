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
	var databasePath, migrationsPath, migrationsTable, command, down string
	flag.StringVar(&databasePath, "database-path", "./database", "path to database folder")
	flag.StringVar(&migrationsPath, "migrations-path", "./database/migration_up.sql", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&command, "command", "", "up all migrations")

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
	if down == "up" {
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("migrations up to date")
			}
			panic(err)
		}
	} else {

	}

}
