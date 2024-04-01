package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang-api-hexagonal/adapters/api/dto"
	"golang-api-hexagonal/core/domain"
	"golang-api-hexagonal/core/ports"
	"net/http"
	"strings"
)

// JWTVerify jwt verify token
type JWTVerify struct {
	log     *zap.SugaredLogger
	service ports.IAuthService
}

// NewJWTHandler create new JWT verify handler
func NewJWTHandler(log *zap.SugaredLogger, service ports.IAuthService) *JWTVerify {
	return &JWTVerify{
		log:     log,
		service: service,
	}
}

// JWTVerifyHandler handler jwt verify
func (jw *JWTVerify) JWTVerifyHandler() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, token, err := jw.parseTokenFromRequest(r)
			if err != nil {
				jw.log.Errorf("JWT parsing failed: %v", err)
				dto.RenderErrorResponse(r.Context(), w, http.StatusUnauthorized, err)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), domain.ClaimsKey, *claims))
			r = r.WithContext(context.WithValue(r.Context(), domain.JwtTokenKey, token))
			next.ServeHTTP(w, r)
		})
	}
}

func (jw *JWTVerify) parseTokenFromRequest(r *http.Request) (*domain.AuthClaims, string, error) {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		jw.log.Error("no security header")
		return nil, "", errors.New("no security header")
	}

	tokenString := strings.Split(header, "Bearer ")
	if len(tokenString) < 2 {
		jw.log.Error("no security header token")
		return nil, "", errors.New("no security header token")
	}

	authClaims, err := jw.service.ParseOauthToken(tokenString[1])
	if err != nil {
		return nil, "", err
	}
	return authClaims, header, nil
}
