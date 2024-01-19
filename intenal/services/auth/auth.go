package auth

import (
	"context"
	"errors"
	"fmt"
	"go-sso/database"
	"go-sso/intenal/domain/models"
	"go-sso/intenal/lib/jwt"
	"go-sso/intenal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserDatabase interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
}
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}
type AppDetailsProvider interface {
	AppDetail(ctx context.Context, appId int) (models.AppDetail, error)
}

type Auth struct {
	log                *slog.Logger
	userSaver          UserSaver
	userProvider       UserProvider
	AppDetailsProvider AppDetailsProvider
	tokenTTL           time.Duration
}

func (a *Auth) NewAuth(log *slog.Logger, userSaver UserSaver, userProider UserProvider, AppDetailsProvider AppDetailsProvider) *Auth {
	return &Auth{
		log:                log,
		userSaver:          userSaver,
		userProvider:       userProider,
		AppDetailsProvider: AppDetailsProvider,
		tokenTTL:           tokenTTL,
	}
}
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "Auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registering user")

	//generate password hash
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// Save user to database
	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("username", email))

	log.Info("attemping to Login user")

	//get user data from DB
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	// check password
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid user credentials")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	//get app details
	app, err := a.AppDetailsProvider.AppDetail(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	// if ok generate token
	token, err := jwt.NewToken(user, app, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
