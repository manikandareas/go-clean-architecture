package pkg

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/spf13/viper"
)

const (
	ACCESS_TOKEN_KEY  string = "jwt.accessToken"
	REFRESH_TOKEN_KEY string = "jwt.refreshToken"
)

type JwtService struct {
	config *viper.Viper
}

func NewJwtService(config *viper.Viper) *JwtService {
	return &JwtService{config: config}
}

func (j *JwtService) GenerateJwtToken(claims *model.JwtClaims, secretKey string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, *claims).SignedString([]byte(j.config.GetString(secretKey)))
}

func (j *JwtService) VerifyJwtToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.GetString(secretKey)), nil
	})
}

func (j *JwtService) DecodeJwtToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := j.VerifyJwtToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
