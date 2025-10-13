package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"sync"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

var (
	ErrUserExists      = fmt.Errorf("loginName is already in use")
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

type Authenticator interface {
	AuthenticateUser(username, password string) (jwt string, err error)
	AuthenticateToken(jwt string) (err error)

	CreateUser(username, password string) (id int, err error)
}

type authy struct {
	DB     *database.Database
	userMu sync.Mutex
}

func NewAuthenticator(db *database.Database) *authy {
	a := new(authy)
	a.DB = db
	return a
}

func (a *authy) CreateUser(username, password string) (id int, err error) {
	a.userMu.Lock()
	defer a.userMu.Unlock()
	user, err := a.DB.Queries.GetUserByLoginName(username)
	if err == nil {
		return user.ID, ErrUserExists
	}

	hash, err := GeneratePasswordHash(password, &DefaultPasswordParams)
	if err != nil {
		return -1, err
	}

	return a.DB.Commands.CreateUser(model.CreateUserModel{
		LoginName:    username,
		PasswordHash: hash,
	})
}

func (a *authy) AuthenticateUser(username, password string) (token string, err error) {
	log := zap.L().Named("auth")
	user, err := a.DB.Queries.GetUserByLoginName(username)
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

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "timekeeper",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%v", time.Now().Unix()),
	})
	log.Info("user authenticated. created jwt")
	return t.SignedString(config.HmacSecret())
}

func (a *authy) AuthenticateToken(token string) (err error) {
	log := zap.L().Named("auth")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.HmacSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Error("failed to parse jwt", zap.Error(err))
		return err
	}

	expiration, err := parsedToken.Claims.GetExpirationTime()
	if err != nil {
		log.Error("failed to get expiration", zap.Error(err))
		return err
	}

	if expiration.Before(time.Now()) {
		err := fmt.Errorf("expired")
		log.Error("token is expired", zap.Error(err))
		return err
	}

	log.Debug("authenticated token")
	return nil
}
