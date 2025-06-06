// file: internal/guest/interfaces/rest/response.go
package rest

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse é a estrutura padronizada para erros da API.
type ErrorResponse struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

// Respond converte um payload Go para JSON e o escreve na resposta HTTP.
func Respond(w http.ResponseWriter, r *http.Request, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "could not encode response", http.StatusInternalServerError)
		}
	}
}

// RespondError envia uma resposta de erro JSON padronizada.
func RespondError(w http.ResponseWriter, r *http.Request, codigo string, mensagem string, statusCode int) {
	errResponse := ErrorResponse{
		Codigo:   codigo,
		Mensagem: mensagem,
	}
	Respond(w, r, errResponse, statusCode)
}