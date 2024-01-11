package auth

import (
	"context"
	"go-sso/intenal/domain/models"
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
