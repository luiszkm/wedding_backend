// file: internal/platform/auth/jwt.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserContextKey = contextKey("userID")

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: []byte(secret)}
}

// GenerateToken cria um novo token JWT para um usuário.
func (s *JWTService) GenerateToken(userID uuid.UUID) (string, error) {
	// "Claims" são as informações que carregamos dentro do token.
	claims := jwt.MapClaims{
		"sub": userID.String(),                           // "Subject", o ID do usuário.
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Data de expiração (ex: 7 dias)
		"iat": time.Now().Unix(),                         // "Issued At", quando o token foi criado.
	}

	// Cria o token com o algoritmo de assinatura HS256 e os claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a nossa chave secreta para gerar a string final.
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("falha ao assinar o token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Valida se o algoritmo de assinatura é o que esperamos (HS256).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("falha ao fazer parse do token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("claim 'sub' inválida no token")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("claim 'sub' não é um UUID válido: %w", err)
		}
		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("token inválido")
}
