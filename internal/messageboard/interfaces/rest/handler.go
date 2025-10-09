// file: internal/messageboard/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	guestDomain "github.com/luiszkm/wedding_backend/internal/guest/domain"
	"github.com/luiszkm/wedding_backend/internal/messageboard/application"
	"github.com/luiszkm/wedding_backend/internal/messageboard/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type MessageBoardHandler struct {
	service *application.MessageBoardService
}

func NewMessageBoardHandler(service *application.MessageBoardService) *MessageBoardHandler {
	return &MessageBoardHandler{service: service}
}

func (h *MessageBoardHandler) HandleDeixarRecado(w http.ResponseWriter, r *http.Request) {
	var reqDTO DeixarRecadoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	eventID, err := uuid.Parse(reqDTO.IDEvento)
	if err != nil {
		web.RespondError(w, r, "DADOS_INVALIDOS", "ID do evento inválido: "+reqDTO.IDEvento, http.StatusBadRequest)
		return
	}

	err = h.service.DeixarNovoRecado(r.Context(), eventID, reqDTO.ChaveDeAcesso, reqDTO.NomeDoAutor, reqDTO.Texto)
	if err != nil {
		// Verifica se o erro foi porque a chave de acesso não foi encontrada.
		if errors.Is(err, guestDomain.ErrGrupoNaoEncontrado) {
			web.RespondError(w, r, "CHAVE_INVALIDA", "A chave de acesso fornecida é inválida.", http.StatusNotFound)
			return
		}
		// Outros erros podem ser de validação ou internos.
		log.Printf("ERRO ao deixar recado: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível postar seu recado.", http.StatusInternalServerError)
		return
	}

	// Conforme a documentação, retornamos 202 Accepted.
	w.WriteHeader(http.StatusAccepted)
}

func (h *MessageBoardHandler) HandleListarRecadosAdmin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	recados, err := h.service.ListarRecadosParaAdmin(r.Context(), userID, idCasamento)
	if err != nil {
		log.Printf("ERRO ao listar recados para admin: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a lista de recados.", http.StatusInternalServerError)
		return
	}

	// Mapeia os objetos de domínio para DTOs de resposta
	respDTO := make([]RecadoAdminDTO, len(recados))
	for i, recado := range recados {
		respDTO[i] = RecadoAdminDTO{
			ID:            recado.ID().String(),
			NomeDoAutor:   recado.NomeDoAutor(),
			Texto:         recado.Texto(),
			Status:        recado.Status(),
			EhFavorito:    recado.EhFavorito(),
			DataDeCriacao: recado.DataDeCriacao(),
		}
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}

func (h *MessageBoardHandler) HandleModerarRecado(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	recadoID, err := uuid.Parse(chi.URLParam(r, "idRecado"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do recado é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO ModerarRecadoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	cmd := application.ModeracaoCommand{
		Status:     reqDTO.Status,
		EhFavorito: reqDTO.EhFavorito,
	}

	recadoAtualizado, err := h.service.ModerarRecado(r.Context(), userID, recadoID, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrRecadoNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Recado não encontrado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrStatusInvalidoParaModeracao) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao moderar recado: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	respDTO := RecadoAdminDTO{ // Reutilizando o DTO de admin para a resposta
		ID:            recadoAtualizado.ID().String(),
		NomeDoAutor:   recadoAtualizado.NomeDoAutor(),
		Texto:         recadoAtualizado.Texto(),
		Status:        recadoAtualizado.Status(),
		EhFavorito:    recadoAtualizado.EhFavorito(),
		DataDeCriacao: recadoAtualizado.DataDeCriacao(),
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}
func (h *MessageBoardHandler) HandleListarRecadosPublicos(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	recados, err := h.service.ListarRecadosPublicos(r.Context(), idCasamento)
	if err != nil {
		log.Printf("ERRO ao listar recados públicos: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a lista de recados.", http.StatusInternalServerError)
		return
	}

	// Mapeia os objetos de domínio para DTOs de resposta pública
	respDTO := make([]RecadoPublicoDTO, len(recados))
	for i, recado := range recados {
		respDTO[i] = RecadoPublicoDTO{
			ID:            recado.ID().String(),
			NomeDoAutor:   recado.NomeDoAutor(),
			Texto:         recado.Texto(),
			EhFavorito:    recado.EhFavorito(),
			DataDeCriacao: recado.DataDeCriacao(),
		}
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}
