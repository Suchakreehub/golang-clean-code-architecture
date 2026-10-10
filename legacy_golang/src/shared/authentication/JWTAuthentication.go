package authentication

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/max38/golang-clean-code-architecture/src/config"
	entityuser "github.com/max38/golang-clean-code-architecture/src/domain/entities/user"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	ApiKey  TokenType = "apikey"
)

type IJWTAuthentication interface {
	SignToken() string
}

type jwtAuthentication struct {
	mapClaims *jwtMapClaims
}

type jwtMapClaims struct {
	Claims *entityuser.UserTokenClaimsEntity `json:"claims"`
	jwt.RegisteredClaims
}

func JWTAuthentication(tokenType TokenType, claims *entityuser.UserTokenClaimsEntity) (IJWTAuthentication, error) {
	switch tokenType {
	case Access:
		return newAccessToken(claims), nil
	case Refresh:
		return newRefreshToken(claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func (a *jwtAuthentication) SignToken() string {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	var signedString string
	var keySign = *getSecretKey()
	signedString, _ = token.SignedString(keySign)
	return signedString
}

func newAccessToken(claims *entityuser.UserTokenClaimsEntity) IJWTAuthentication {
	return &jwtAuthentication{
		mapClaims: &jwtMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    getIssuerName(),
				Subject:   "access-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeDurationCalculation(config.Config.Int("GOAPP_JWT_ACCESS_EXPIRES")),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(claims *entityuser.UserTokenClaimsEntity) IJWTAuthentication {
	return &jwtAuthentication{
		mapClaims: &jwtMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    getIssuerName(),
				Subject:   "refresh-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeDurationCalculation(config.Config.Int("GOAPP_JWT_ACCESS_EXPIRES")),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

// --- Support Functions ---

func getIssuerName() string {
	var applicatiobName = config.Config.String("GOAPP_NAME")
	strings.Replace(applicatiobName, " ", "-", -1)
	applicatiobName = strings.ToLower(applicatiobName)
	return applicatiobName
}

func getSecretKey() *[]byte {
	var secretKey = []byte(config.Config.String("GOAPP_SECRET_KEY"))
	return &secretKey
}

func jwtTimeDurationCalculation(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func ParseToken(tokenString string) (*jwtMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return *getSecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*jwtMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

// User for Refresh Token when we need ratain expire time in refresh token
func RepeatToken(claims *entityuser.UserTokenClaimsEntity, expireTimestamp int64) IJWTAuthentication {
	return &jwtAuthentication{
		mapClaims: &jwtMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    getIssuerName(),
				Subject:   "refresh-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(expireTimestamp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
