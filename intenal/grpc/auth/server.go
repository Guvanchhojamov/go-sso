package auth

import (
	"context"
	"errors"
	ssov1 "github.com/GolangLessons/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emtyValueString = ""
	emptyValueInt   = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email, password string, appId int32) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userId int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}
func (s *serverAPI) Login(ctx context.Context, input *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// validate login creditianals
	if correctInput, err := s.validateLoginInput(input); !correctInput {
		return nil, err
	}
	token, err := s.auth.Login(ctx, input.GetEmail(), input.GetPassword(), input.GetAppId())
	if err != nil {
		// Ошибку auth.ErrInvalidCredentials мы создадим ниже
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "login failed")
	}
	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, input *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	// TODO
	if correctInput, err := s.validateRegisterInput(input); !correctInput {
		return nil, err
	}

	uid, err := s.auth.RegisterNewUser(ctx, input.GetEmail(), input.GetPassword())
	if err != nil {
		if errors.Is(err, db.ErrUserExits) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "register new user failed")
	}
	return &ssov1.RegisterResponse{UserId: uid}, nil
}

// validate inputs
func (s *serverAPI) validateLoginInput(input *ssov1.LoginRequest) (correct bool, err error) {
	if input.Email == emtyValueString || len(input.Email) == emptyValueInt {
		return false, status.Error(codes.InvalidArgument, "email is required")
	}
	if input.Password == emtyValueString || len(input.Password) == emptyValueInt {
		return false, status.Error(codes.InvalidArgument, "Password is required")
	}
	if input.GetAppId() == emptyValueInt {
		return false, status.Error(codes.InvalidArgument, "app_id is required")
	}
	return true, nil
}
func (s *serverAPI) validateRegisterInput(input *ssov1.RegisterRequest) (correct bool, err error) {
	if input.Email == emtyValueString || len(input.Email) == emptyValueInt {
		return false, status.Error(codes.InvalidArgument, "email is required")
	}
	if input.Password == emtyValueString || len(input.Password) == emptyValueInt {
		return false, status.Error(codes.InvalidArgument, "Password is required")
	}
	return true, nil
}
