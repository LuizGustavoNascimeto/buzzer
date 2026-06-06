package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	jwks keyfunc.Keyfunc
}

func New(region, userPoolID string) (*Validator, error) {
	jwksURL := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region, userPoolID,
	)

	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar JWKS: %w", err)
	}

	return &Validator{jwks: k}, nil
}

func (v *Validator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims inválidos")
	}

	return claims, nil
}
