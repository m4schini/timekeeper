package auth

import (
	"fmt"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/config"
	"strconv"
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
	AuthenticateToken(jwt string) (userId int, role model.Role, err error)

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

func (a *authy) AuthenticateToken(token string) (userId int, role model.Role, err error) {
	log := zap.L().Named("auth")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.HmacSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Error("failed to parse jwt", zap.Error(err))
		return -1, model.RoleParticipant, err
	}

	expiration, err := parsedToken.Claims.GetExpirationTime()
	if err != nil {
		log.Error("failed to get expiration", zap.Error(err))
		return -1, model.RoleParticipant, err
	}

	if expiration.Before(time.Now()) {
		err := fmt.Errorf("expired")
		log.Error("token is expired", zap.Error(err))
		return -1, model.RoleParticipant, err
	}

	userIdStr, err := parsedToken.Claims.GetSubject()
	if err != nil {
		err := fmt.Errorf("invalid token")
		log.Error("token misses user", zap.Error(err))
		return -1, model.RoleParticipant, err
	}

	_userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid userId")
		log.Error("cannot parse subjet", zap.Error(err), zap.String("userId", userIdStr))
		return -1, model.RoleParticipant, err
	}

	return int(_userId), model.RoleOrganizer, nil
}
