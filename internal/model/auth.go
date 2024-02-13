package model

import "github.com/golang-jwt/jwt/v4"

type Auth struct {
	// Login user id
	ID    string
	Email string
	Name  string
}

type JwtClaims struct {
	jwt.RegisteredClaims
	User map[string]string
}

type BackendTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
