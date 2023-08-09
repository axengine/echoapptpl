package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var (
	defaultSecretKey = []byte("WpQ3TaXgkRNj9PAP")
	secret_key       = defaultSecretKey
	issuer           = "JWT"
)

func Init(key string, issuer string) {
	secret_key = []byte(key)
	issuer = issuer
}

type CustomClaims struct {
	Uid   uint `json:"uid"`
	Role  uint `json:"role"` //0-默认 1-admin
	Level uint `json:"level"`
	jwt.RegisteredClaims
}

func GenToken(uid uint, role uint, level uint) (string, error) {
	claim := CustomClaims{
		uid,
		role,
		level,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), //过期时间
			Issuer:    issuer,                                                 //签发人
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret_key))
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	}
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

func GetClaimFromContext(c echo.Context) CustomClaims {
	return *c.Get("claims").(*CustomClaims)
}
