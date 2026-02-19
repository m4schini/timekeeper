package auth

import (
	"fmt"
	"raumzeitalpaka/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	UserId any    `json:"userId"`
	Issuer string `json:"issuer"`
}

func NewJWT(claims Claims) (string, error) {
	expiresAt := time.Now().Add(24 * 7 * time.Hour)
	jwtId := time.Now().Unix()

	if claims.Issuer == "" {
		claims.Issuer = "local"
	}

	// just completely beating registered claims to my usecase because im too lazy to look up how to do this correctly
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("raumzeitalpaka/%v", claims.Issuer),
		Subject:   fmt.Sprintf("%v", claims.UserId),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%v", jwtId),
	})
	return t.SignedString(config.HmacSecret())
}

func AuthenticateJWT(token string) (userId int, err error) {
	log := zap.L().Named("auth")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.HmacSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Error("failed to parse jwt", zap.Error(err))
		return -1, err
	}

	expiration, err := parsedToken.Claims.GetExpirationTime()
	if err != nil {
		log.Error("failed to get expiration", zap.Error(err))
		return -1, err
	}

	if expiration.Before(time.Now()) {
		err = fmt.Errorf("expired")
		log.Error("token is expired", zap.Error(err))
		return -1, err
	}

	userIdStr, err := parsedToken.Claims.GetSubject()
	if err != nil {
		err = fmt.Errorf("invalid token")
		log.Error("token misses user", zap.Error(err))
		return -1, err
	}

	_userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid userId")
		log.Error("cannot parse subjet", zap.Error(err), zap.String("userId", userIdStr))
		return -1, err
	}

	return int(_userId), err
}
