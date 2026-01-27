package auth

import (
	"fmt"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	UserId   int        `json:"userId"`
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
}

func NewJWT(claims Claims) (string, error) {
	expiresAt := time.Now().Add(72 * time.Hour)
	jwtId := time.Now().Unix()

	// just completely beating registered claims to my usecase because im too lazy to look up how to do this correctly
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    claims.Username,
		Subject:   fmt.Sprintf("%v", claims.UserId),
		Audience:  jwt.ClaimStrings{string(claims.Role)},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%v", jwtId),
	})
	return t.SignedString(config.HmacSecret())
}

func AuthenticateJWT(token string) (userId int, name string, role model.Role, err error) {
	log := zap.L().Named("auth")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.HmacSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Error("failed to parse jwt", zap.Error(err))
		return -1, "unknown", model.RoleParticipant, err
	}
	log.Info("parsed jwt", zap.Any("claims", parsedToken.Claims))

	expiration, err := parsedToken.Claims.GetExpirationTime()
	if err != nil {
		log.Error("failed to get expiration", zap.Error(err))
		return -1, "unknown", model.RoleParticipant, err
	}

	if expiration.Before(time.Now()) {
		err := fmt.Errorf("expired")
		log.Error("token is expired", zap.Error(err))
		return -1, "unknown", model.RoleParticipant, err
	}

	userIdStr, err := parsedToken.Claims.GetSubject()
	if err != nil {
		err := fmt.Errorf("invalid token")
		log.Error("token misses user", zap.Error(err))
		return -1, "unknown", model.RoleParticipant, err
	}

	_userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid userId")
		log.Error("cannot parse subjet", zap.Error(err), zap.String("userId", userIdStr))
		return -1, "unknown", model.RoleParticipant, err
	}

	aud, err := parsedToken.Claims.GetAudience()
	if err != nil {
		return int(_userId), "unknown", model.RoleParticipant, nil
	}

	username, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return int(_userId), "unknown", model.RoleParticipant, nil
	}

	role = model.Role(aud[0])
	log.Info("parsed jwt", zap.Any("user", _userId), zap.Any("role", role))
	return int(_userId), username, role, err
}
