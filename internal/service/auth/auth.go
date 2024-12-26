package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/paveldroo/sso-service/internal/domain/models"
	"github.com/paveldroo/sso-service/internal/lib/jwt"
	"github.com/paveldroo/sso-service/internal/lib/logger/sl"
	"github.com/paveldroo/sso-service/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid user credentials")

type UserSaver interface {
	AddUser(ctx context.Context, email string, pass_hash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Auth struct {
	log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	AppProvider  AppProvider
	tokenTTL     time.Duration
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) Auth {
	return Auth{
		log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		AppProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	slog.Info("attempting login user")

	user, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			slog.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("get user from storage: %w", ErrInvalidCredentials)
		}
		slog.Warn("failed to get user")
		return "", fmt.Errorf("get user from storage: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		slog.Warn("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("compare user password and hash: %w", ErrInvalidCredentials)
	}

	app, err := a.AppProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("get app from storage: %w", err)
	}

	slog.Info("user logged successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		slog.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("get new token: %w", err)
	}

	return token, nil
}

func (a Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	slog.Info("registering user")
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate hash from password: %w", err)
	}

	userID, err := a.UserSaver.AddUser(ctx, email, hashedPass)
	if err != nil {
		return 0, fmt.Errorf("add user to storage: %w", err)
	}

	return userID, nil
}

func (a Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	slog.Info("checking if user is admin")
	isAdmin, err := a.UserProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("get user is_admin property from storage: %w", err)
	}

	slog.Info("checked is user admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
