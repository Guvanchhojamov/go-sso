package auth

import (
	"context"
	"fmt"
	"go-sso/intenal/domain/models"
	"go-sso/intenal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	tokenTTL = time.Second * 10
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
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appDetails   AppDetailsProvider
	tokenTTL     time.Duration
}

func (a *Auth) NewAuth(log *slog.Logger, userSaver UserSaver, userProider UserProvider, appDetails AppDetailsProvider) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProider,
		appDetails:   appDetails,
		tokenTTL:     tokenTTL,
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
