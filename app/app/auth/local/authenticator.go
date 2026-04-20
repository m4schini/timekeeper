package local

import (
	"context"
	"fmt"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/config"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrUserExists      = fmt.Errorf("loginName is already in use")
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

type Authenticator interface {
	AuthenticateUser(username, password string) (jwt string, err error)

	CreateUser(username, password string) (id int, err error)
}

type authy struct {
	getUserByLoginName query.GetUserByLoginName
	createUser         command.CreateUser
	updateLastLogin    command.UpdateLastLogin
	userMu             sync.Mutex
}

func NewAuthenticator(db *database.Database) *authy {
	a := new(authy)
	a.getUserByLoginName = db.Queries.UserByLoginName
	a.createUser = db.Commands.CreateUser
	a.updateLastLogin = db.Commands.UpdateLastLogin
	return a
}

func (a *authy) CreateUser(username, password string) (id int, err error) {
	a.userMu.Lock()
	defer a.userMu.Unlock()
	user, err := a.getUserByLoginName.Query(context.TODO(), query.GetUserByLoginNameRequest{LoginName: username})
	if err == nil {
		return user.ID, ErrUserExists
	}

	hash, err := GeneratePasswordHash(password, &DefaultPasswordParams)
	if err != nil {
		return -1, err
	}

	return a.createUser.Execute(context.TODO(), command.CreateUserRequest{
		LoginName:    username,
		PasswordHash: hash,
	})
}

func (a *authy) AuthenticateUser(username, password string) (token string, err error) {
	log := zap.L().Named("auth")
	user, err := a.getUserByLoginName.Query(context.TODO(), query.GetUserByLoginNameRequest{LoginName: username})
	if err != nil {
		return "", err
	}

	matches, err := ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !matches {
		return "", ErrInvalidPassword
	}

	//go a.updateLastLogin.Execute(context.TODO(), command.UpdateLastLoginRequest{
	//	ID:        user.ID,
	//	Timestamp: time.Now(),
	//})

	expiresAt := time.Now().Add(72 * time.Hour)
	jwtId := time.Now().Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "raumzeitalpaka",
		Subject:   fmt.Sprintf("%v", user.ID),
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%v", jwtId),
	})
	log.Info("user authenticated - jwt created", zap.Int("user", user.ID), zap.String("username", user.LoginName), zap.Time("expires_at", expiresAt), zap.Int64("jwt_id", jwtId))
	return t.SignedString(config.HmacSecret())
}
