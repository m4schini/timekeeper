package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
	"timekeeper/config"
)

type Authenticator interface {
	AuthenticateUser(username, password string) (jwt string, err error)
	AuthenticateToken(jwt string) (err error)
}

type authy struct {
}

func NewAuthenticator() *authy {
	return new(authy)
}

func (a *authy) AuthenticateUser(username, password string) (token string, err error) {
	log := zap.L().Named("auth")
	if username != "admin" || password != config.AdminPassword() {
		err := fmt.Errorf("invalid credentials")
		log.Error("login failed", zap.Error(err))
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "timekeeper",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(18 * time.Hour)),
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
