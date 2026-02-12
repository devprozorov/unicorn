package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
}

type AccessClaims struct {
	UserID string   `json:"sub"`
	Type   string   `json:"typ"`
	AMR    []string `json:"amr"`
	jwt.RegisteredClaims
}

type MFAClaims struct {
	UserID string `json:"sub"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{secret: []byte(secret)}
}

func (t *TokenService) NewAccessToken(userID, typ string, amr []string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		Type:   typ,
		AMR:    amr,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
}

func (t *TokenService) NewMFAToken(userID, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := MFAClaims{
		UserID: userID,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
}

func (t *TokenService) ParseAccess(token string) (*AccessClaims, error) {
	var claims AccessClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil { return nil, err }
	if claims.UserID == "" || claims.Type == "" { return nil, errors.New("bad claims") }
	return &claims, nil
}

func (t *TokenService) ParseMFA(token string) (*MFAClaims, error) {
	var claims MFAClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil { return nil, err }
	if claims.UserID == "" || claims.Type == "" { return nil, errors.New("bad claims") }
	return &claims, nil
}
