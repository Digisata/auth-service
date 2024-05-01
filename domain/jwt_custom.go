package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type (
	JwtCustomClaims struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		jwt.StandardClaims
	}

	JwtCustomRefreshClaims struct {
		ID string `json:"id"`
		jwt.StandardClaims
	}
)
