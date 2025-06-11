// file: internal/iam/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/luiszkm/wedding_backend/internal/iam/application"
	"github.com/luiszkm/wedding_backend/internal/iam/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type IAMHandler struct {
	service *application.IAMService
}

func NewIAMHandler(service *application.IAMService) *IAMHandler {
	return &IAMHandler{service: service}
}

func (h *IAMHandler) HandleRegistrar(w http.ResponseWriter, r *http.Request) {
	var reqDTO RegistrarRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	_, err := h.service.RegistrarNovoUsuario(r.Context(), reqDTO.Nome, reqDTO.Email, reqDTO.Telefone, reqDTO.Senha)
	if err != nil {
		// Trata erros de negócio específicos
		if errors.Is(err, domain.ErrEmailJaExiste) {
			web.RespondError(w, r, "EMAIL_EM_USO", err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		// Outros erros podem ser de validação do domínio ou erros técnicos.
		web.RespondError(w, r, "ERRO_REGISTRO", err.Error(), http.StatusBadRequest)
		return
	}

	// Conforme a documentação, apenas retornamos sucesso com status 201 Created.
	web.Respond(w, r, nil, http.StatusCreated)
}

func (h *IAMHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var reqDTO LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "Corpo da requisição malformado.", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), reqDTO.Email, reqDTO.Senha)
	if err != nil {
		if errors.Is(err, domain.ErrCredenciaisInvalidas) {
			web.RespondError(w, r, "AUTENTICACAO_FALHOU", err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("ERRO no login: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar login.", http.StatusInternalServerError)
		return
	}
	isSecure := os.Getenv("APP_ENV") == "production"

	// --- LÓGICA DE COOKIE ---

	// 1. Define a data de expiração do cookie (ex: 7 dias, igual ao JWT)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// 2. Cria o cookie com os atributos de segurança
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken", // Nome do cookie
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,                 // <-- Impede o acesso via JavaScript (proteção XSS)
		Secure:   isSecure,             // <-- Envia o cookie apenas sobre HTTPS (use 'false' apenas em localhost sem TLS)
		SameSite: http.SameSiteLaxMode, // <-- Proteção contra CSRF
		Path:     "/",                  // O cookie será válido para todo o site
	})

	// 3. A resposta agora não precisa mais enviar o token no corpo.
	// Podemos enviar uma mensagem de sucesso ou os dados do usuário.
	web.Respond(w, r, map[string]string{"status": "login bem-sucedido"}, http.StatusOK)
}
