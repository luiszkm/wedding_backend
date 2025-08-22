// file: internal/platform/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/luiszkm/wedding_backend/internal/platform/web" // Ajuste o path se necessário
)

// Authenticator é um middleware Chi para validar o token JWT.
// Suporta tanto cookies quanto Authorization Bearer headers.
func Authenticator(jwtSvc *JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			
			// 1. Tentar obter token do Authorization Header primeiro
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			} else {
				// 2. Se não encontrar no header, tentar no cookie
				cookie, err := r.Cookie("accessToken")
				if err != nil {
					if err == http.ErrNoCookie {
						web.RespondError(w, r, "TOKEN_AUSENTE", "Token de autorização não fornecido.", http.StatusUnauthorized)
						return
					}
					web.RespondError(w, r, "REQUISICAO_INVALIDA", "Erro ao ler cookie.", http.StatusBadRequest)
					return
				}
				tokenString = cookie.Value
			}

			// 3. Valida o token.
			userID, err := jwtSvc.ValidateToken(tokenString)
			if err != nil {
				web.RespondError(w, r, "TOKEN_INVALIDO", "O token fornecido é inválido ou expirou.", http.StatusUnauthorized)
				return
			}

			// 4. Se o token for válido, injeta o userID no contexto da requisição.
			ctx := context.WithValue(r.Context(), UserContextKey, userID)

			// 5. Chama o próximo handler na cadeia, agora com o novo contexto.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
