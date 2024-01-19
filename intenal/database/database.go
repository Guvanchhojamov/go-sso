package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"go-sso/database"
	"go-sso/intenal/domain/models"
	"log"
)

type DBSqlite struct {
	db *sql.DB
}

func NewDBSqlite(databasePath string) (*DBSqlite, error) {
	const op = "database.NewDBSqlite"
	// path to file
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("%s: %w", op, err)
	}
	return &DBSqlite{db: db}, nil
}
func (s *DBSqlite) SaveUSer(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "database.database.SaveUser"
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		// Небольшое кунг-фу для выявления ошибки ErrConstraintUnique
		// (см. подробности ниже)
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, database.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *DBSqlite) User(ctx context.Context, email string) (models.User, error) {
	const op = "database.database.User"
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, database.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *DBSqlite) App(ctx context.Context, appId int) (models.AppDetail, error) {
	const op = "database.database.App"
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.AppDetail{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appId)

	var app models.AppDetail
	err = row.Scan(&app.ID, &app.Name, &app.SecretKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.AppDetail{}, fmt.Errorf("%s: %w", op, database.ErrAppNotFound)
		}
		return models.AppDetail{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
