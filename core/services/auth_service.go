package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang-api-hexagonal/config"
	"golang-api-hexagonal/core/domain"
	"time"
)

// AuthService service to authenticate
type AuthService struct {
	log  *zap.SugaredLogger
	conf config.Oauth
}

// NewAuthService create new auth service
func NewAuthService(log *zap.SugaredLogger, conf config.Oauth) *AuthService {
	return &AuthService{
		log:  log,
		conf: conf,
	}
}

// CreateOauthToken create the oauth token
func (as *AuthService) CreateOauthToken(request *domain.Auth, traceID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"username": request.Username, "role": domain.SuperRole, "iss": domain.Issuer, "exp": time.Now().Add(time.Hour * 1).Unix()})

	tokenString, err := token.SignedString([]byte(as.conf.Secret))
	if err != nil {
		as.log.With("traceId", traceID).Errorf("Error to sign token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// ParseOauthToken parse the oauth token to claims
func (as *AuthService) ParseOauthToken(tokenString string) (*domain.AuthClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(as.conf.Secret), nil
	})

	if err != nil || !token.Valid {
		as.log.Errorf("invalid token: %v", err)
		return nil, errors.New("invalid token")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != domain.Issuer {
		as.log.Errorf("invalid token issuer: %v", err)
		return nil, errors.New("invalid token issuer")
	}

	expiration, err := token.Claims.GetExpirationTime()
	if err != nil || expiration.Time.Before(time.Now()) {
		as.log.Errorf("token expired: %v", err)
		return nil, errors.New("token expired")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &domain.AuthClaims{
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
		}, nil
	} else {
		as.log.Errorf("claims not found")
		return nil, errors.New("claims not found")
	}
}
