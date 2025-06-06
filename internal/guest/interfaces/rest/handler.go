// file: internal/guest/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/guest/application"
	"github.com/luiszkm/wedding_backend/internal/guest/domain"
)

type GuestHandler struct {
	service *application.GuestService
}

func NewGuestHandler(service *application.GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

func (h *GuestHandler) HandleCriarGrupoDeConvidados(w http.ResponseWriter, r *http.Request) {
	idCasamentoStr := chi.URLParam(r, "idCasamento")
	idCasamento, err := uuid.Parse(idCasamentoStr)
	if err != nil {
		RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO CriarGrupoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	idGrupo, err := h.service.CriarNovoGrupo(
		r.Context(),
		idCasamento,
		reqDTO.ChaveDeAcesso,
		reqDTO.NomesDosConvidados,
	)
	if err != nil {
		// Mapeia erros do domínio para respostas HTTP apropriadas
		if errors.Is(err, domain.ErrChaveDeAcessoObrigatoria) || errors.Is(err, domain.ErrPeloMenosUmConvidado) {
			RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		// Outros erros são tratados como internos
		RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	respDTO := CriarGrupoResponseDTO{IDGrupo: idGrupo.String()}
	Respond(w, r, respDTO, http.StatusCreated)
}

func (h *GuestHandler) HandleObterGrupoPorChaveDeAcesso(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o query parameter da URL.
	chaveDeAcesso := r.URL.Query().Get("chave")
	if chaveDeAcesso == "" {
		RespondError(w, r, "PARAMETRO_AUSENTE", "O parâmetro 'chave' é obrigatório.", http.StatusBadRequest)
		return
	}

	// 2. Chamar a camada de aplicação.
	grupo, err := h.service.ObterGrupoPorChaveDeAcesso(r.Context(), chaveDeAcesso)
	if err != nil {
		// Se o erro for "não encontrado", retornamos 404.
		if errors.Is(err, domain.ErrGrupoNaoEncontrado) {
			RespondError(w, r, "NAO_ENCONTRADO", "Nenhuma chave de acesso correspondente foi encontrada.", http.StatusNotFound)
			return
		}
		// Outros erros são internos.
		log.Printf("ERRO: %v\n", err)
		RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	// 3. Mapear o agregado de domínio para o DTO de resposta.
	convidadosDTO := make([]ConvidadoDTO, len(grupo.Convidados()))
	for i, c := range grupo.Convidados() {
		convidadosDTO[i] = ConvidadoDTO{
			ID:         c.ID().String(),
			Nome:       c.Nome(),
			StatusRSVP: c.StatusRSVP(),
		}
	}
	respDTO := GrupoParaConfirmacaoDTO{
		IDGrupo:    grupo.ID().String(),
		Convidados: convidadosDTO,
	}

	// 4. Responder com sucesso.
	Respond(w, r, respDTO, http.StatusOK)
}
