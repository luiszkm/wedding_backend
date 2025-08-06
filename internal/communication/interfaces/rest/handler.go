package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/communication/application"
	"github.com/luiszkm/wedding_backend/internal/communication/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type CommunicationHandler struct {
	service *application.CommunicationService
}

func NewCommunicationHandler(service *application.CommunicationService) *CommunicationHandler {
	return &CommunicationHandler{service: service}
}

func (h *CommunicationHandler) HandleCriarComunicado(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}

	idEvento, err := uuid.Parse(chi.URLParam(r, "idEvento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do evento é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO CriarComunicadoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	comunicado, err := h.service.CriarComunicado(r.Context(), userID, idEvento, reqDTO.Titulo, reqDTO.Mensagem)
	if err != nil {
		if errors.Is(err, domain.ErrTituloObrigatorio) || errors.Is(err, domain.ErrTituloMuitoLongo) || errors.Is(err, domain.ErrMensagemObrigatoria) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao criar comunicado: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível criar o comunicado.", http.StatusInternalServerError)
		return
	}

	respDTO := ComunicadoResponseDTO{
		ID:             comunicado.ID().String(),
		IDEvento:       comunicado.IDEvento().String(),
		Titulo:         comunicado.Titulo(),
		Mensagem:       comunicado.Mensagem(),
		DataPublicacao: comunicado.DataPublicacao(),
	}

	web.Respond(w, r, respDTO, http.StatusCreated)
}

func (h *CommunicationHandler) HandleEditarComunicado(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}

	idComunicado, err := uuid.Parse(chi.URLParam(r, "idComunicado"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do comunicado é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO EditarComunicadoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	err = h.service.EditarComunicado(r.Context(), userID, idComunicado, reqDTO.Titulo, reqDTO.Mensagem)
	if err != nil {
		if errors.Is(err, domain.ErrComunicadoNaoEncontrado) {
			web.RespondError(w, r, "COMUNICADO_NAO_ENCONTRADO", "Comunicado não encontrado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTituloObrigatorio) || errors.Is(err, domain.ErrTituloMuitoLongo) || errors.Is(err, domain.ErrMensagemObrigatoria) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao editar comunicado: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível editar o comunicado.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CommunicationHandler) HandleDeletarComunicado(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}

	idComunicado, err := uuid.Parse(chi.URLParam(r, "idComunicado"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do comunicado é inválido.", http.StatusBadRequest)
		return
	}

	err = h.service.DeletarComunicado(r.Context(), userID, idComunicado)
	if err != nil {
		if errors.Is(err, domain.ErrComunicadoNaoEncontrado) {
			web.RespondError(w, r, "COMUNICADO_NAO_ENCONTRADO", "Comunicado não encontrado.", http.StatusNotFound)
			return
		}
		log.Printf("ERRO ao deletar comunicado: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível deletar o comunicado.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CommunicationHandler) HandleListarComunicados(w http.ResponseWriter, r *http.Request) {
	idEvento, err := uuid.Parse(chi.URLParam(r, "idEvento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do evento é inválido.", http.StatusBadRequest)
		return
	}

	comunicados, err := h.service.ListarComunicadosPorEvento(r.Context(), idEvento)
	if err != nil {
		log.Printf("ERRO ao listar comunicados: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar comunicados.", http.StatusInternalServerError)
		return
	}

	respDTO := make([]ComunicadoResponseDTO, len(comunicados))
	for i, comunicado := range comunicados {
		respDTO[i] = ComunicadoResponseDTO{
			ID:             comunicado.ID().String(),
			IDEvento:       comunicado.IDEvento().String(),
			Titulo:         comunicado.Titulo(),
			Mensagem:       comunicado.Mensagem(),
			DataPublicacao: comunicado.DataPublicacao(),
		}
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}
