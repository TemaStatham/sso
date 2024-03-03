package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/TemaStatham/sso/internal/domain/models"
	"github.com/TemaStatham/sso/internal/lib/jwt"
	"github.com/TemaStatham/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorInvalidCredentials = errors.New("invalid credentional")
	ErrorAppID              = errors.New("app id is unvalue")
	ErrorUserExist          = errors.New("user already exist")
)

// Auth сервис аунтетификации
type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

// UserSaver интерфейс сохранения пользователя
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (ui int64, err error)
}

// UserProvider интерфейс логинации
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// AppProvider интерфейс приложения
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New возвращает новый Auth сервис
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("logining user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		a.log.Warn("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credential", err)
		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is logined successfuly")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Warn("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// Register
func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			a.log.Warn("app id not found", err)
			return 0, fmt.Errorf("%s: %w", op, ErrorUserExist)
		}
		log.Error("failed to get password hash: ", err)
		return 0, fmt.Errorf("%s %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user: ", err)
	}

	log.Info("user registrated")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)),
	)

	log.Info("is admin checked")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app id not found", err)
			return false, fmt.Errorf("%s: %w", op, ErrorAppID)
		}

		log.Error("is admin failed: ", err)
		return false, fmt.Errorf("%s %w", op, err)
	}

	log.Info("is admin user successfully checked")

	return isAdmin, nil
}
