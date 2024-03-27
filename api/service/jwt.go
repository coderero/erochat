package service

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	// RSAPrivateKey is the RSA private key.
	RSAPrivateKey *rsa.PrivateKey

	// RSAPublicKey is the RSA public key.
	RSAPublicKey *rsa.PublicKey

	// TokenDuration is the duration of the token.
	TokenDuration time.Duration

	// RefreshTokenDuration is the duration of the refresh token.
	RefreshTokenDuration time.Duration
}

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

func (t TokenType) String() string {
	return [...]string{
		"access",
		"refresh",
	}[t]
}

func (t TokenType) Duration(s *JWTService) time.Duration {
	switch t {
	case AccessToken:
		return s.TokenDuration
	case RefreshToken:
		return s.RefreshTokenDuration
	}
	return 0
}

// NewJWTService creates a new JWTService.
func NewJWTService(privateKey, publicKey []byte, tokenDuration, refreshTokenDuration time.Duration) (*JWTService, error) {
	var (
		rsaPrivateKey *rsa.PrivateKey
		rsaPublicKey  *rsa.PublicKey
		err           error
	)

	// Parse the private key.
	if rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey); err != nil {
		return nil, err
	}

	// Parse the public key.
	if rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey); err != nil {
		return nil, err
	}

	return &JWTService{
		RSAPrivateKey:        rsaPrivateKey,
		RSAPublicKey:         rsaPublicKey,
		TokenDuration:        tokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}, nil
}

// GenerateTokens generates a token and a refresh token.
func (s *JWTService) GenerateTokens(email string, userId uuid.UUID) (string, string, error) {
	var (
		token        string
		refreshToken string
		err          error
	)

	// Create a new token.
	token, err = s.createToken(email, userId, time.Now().Add(s.TokenDuration).Unix())
	if err != nil {
		return "", "", err
	}

	// Create a new refresh token.
	refreshToken, err = s.createToken(email, userId, time.Now().Add(s.RefreshTokenDuration).Unix())
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// ValidateToken validates a token.
func (s *JWTService) ValidateToken(tokenString string) (bool, error) {
	var (
		token *jwt.Token
		err   error
	)
	// Parse the token.
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.RSAPublicKey, nil
	})
	if err != nil {
		return false, err
	}
	// Validate the token.
	if !token.Valid {
		return false, nil
	}
	return true, nil
}

// GetClaims gets the claims from a token.
func (s *JWTService) GetClaims(tokenString string) (jwt.MapClaims, error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
		err    error
	)

	// Parse the token.
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.RSAPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Get the claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

// GenerateToken generates a token.
func (s *JWTService) GenerateToken(email string, userId uuid.UUID, tokenType TokenType) (string, error) {
	d := tokenType.Duration(s)
	return s.createToken(email, userId, time.Now().Add(d).Unix())
}

func (s *JWTService) RefreshToken(refreshToken string) (string, error) {
	// Validate the refresh token.
	if ok, err := s.ValidateToken(refreshToken); err != nil || !ok {
		return "", err
	}

	// Get the claims from the refresh token.
	claims, err := s.GetClaims(refreshToken)
	if err != nil {
		return "", err
	}

	// Get the email from the claims.
	email, ok := claims["sub"].(string)
	if !ok {
		return "", err
	}

	// Get the user id from the claims.
	userId, ok := claims["uid"].(string)
	if !ok {
		return "", err
	}

	// Parse the user id.
	uid, err := uuid.Parse(userId)
	if err != nil {
		return "", err
	}

	// Create a new token.
	return s.createToken(email, uid, time.Now().Add(s.TokenDuration).Unix())
}

// createToken creates a token.
func (s *JWTService) createToken(email string, userId uuid.UUID, duration int64) (string, error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
		err    error
	)

	// Create a new token.
	claims = jwt.MapClaims{
		"iss": "erosecurity",
		"uid": userId,
		"sub": email,
		"exp": time.Now().Add(time.Duration(duration) * time.Second).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Set the duration of the token.
	tokenString, err := token.SignedString(s.RSAPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
