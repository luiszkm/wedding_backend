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
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type GuestHandler struct {
	service *application.GuestService
}

func NewGuestHandler(service *application.GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

func (h *GuestHandler) HandleCriarGrupoDeConvidados(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		// Este erro não deveria acontecer se o middleware estiver funcionando.
		web.RespondError(w, r, "ERRO_CONTEXTO", "Não foi possível obter o ID do usuário.", http.StatusInternalServerError)
		return
	}
	idCasamentoStr := chi.URLParam(r, "idCasamento")
	idCasamento, err := uuid.Parse(idCasamentoStr)
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO CriarGrupoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
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
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		// Outros erros são tratados como internos
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	respDTO := CriarGrupoResponseDTO{IDGrupo: idGrupo.String()}
	web.Respond(w, r, respDTO, http.StatusCreated)
}

func (h *GuestHandler) HandleObterGrupoPorChaveDeAcesso(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o query parameter da URL.
	chaveDeAcesso := r.URL.Query().Get("chave")
	if chaveDeAcesso == "" {
		web.RespondError(w, r, "PARAMETRO_AUSENTE", "O parâmetro 'chave' é obrigatório.", http.StatusBadRequest)
		return
	}

	// 2. Chamar a camada de aplicação.
	grupo, err := h.service.ObterGrupoPorChaveDeAcesso(r.Context(), chaveDeAcesso)
	if err != nil {
		// Se o erro for "não encontrado", retornamos 404.
		if errors.Is(err, domain.ErrGrupoNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Nenhuma chave de acesso correspondente foi encontrada.", http.StatusNotFound)
			return
		}
		// Outros erros são internos.
		log.Printf("ERRO: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
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
	web.Respond(w, r, respDTO, http.StatusOK)
}

func (h *GuestHandler) HandleConfirmarPresenca(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar o corpo da requisição.
	var reqDTO ConfirmarPresencaRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	// 2. Converter o DTO da camada de interface para o tipo do domínio.
	respostasDominio := make([]domain.RespostaRSVP, len(reqDTO.Respostas))
	for i, rsvpDTO := range reqDTO.Respostas {
		convidadoID, err := uuid.Parse(rsvpDTO.IDConvidado)
		if err != nil {
			web.RespondError(w, r, "DADOS_INVALIDOS", "ID de convidado inválido: "+rsvpDTO.IDConvidado, http.StatusBadRequest)
			return
		}
		respostasDominio[i] = domain.RespostaRSVP{
			ConvidadoID: convidadoID,
			Status:      rsvpDTO.Status,
		}
	}

	// 3. Chamar o serviço de aplicação.
	err := h.service.ConfirmarPresencaGrupo(r.Context(), reqDTO.ChaveDeAcesso, respostasDominio)
	if err != nil {
		if errors.Is(err, domain.ErrGrupoNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Chave de acesso não encontrada.", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrStatusRSVPInvalido) || errors.Is(err, domain.ErrConvidadoNaoEncontradoNoGrupo) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	// 4. Responder com sucesso.
	w.WriteHeader(http.StatusNoContent)
}

func (h *GuestHandler) HandleRevisarGrupo(w http.ResponseWriter, r *http.Request) {
	grupoID, err := uuid.Parse(chi.URLParam(r, "idGrupo"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do grupo é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO RevisarGrupoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	convidadosDominio := make([]domain.ConvidadoParaRevisao, len(reqDTO.Convidados))
	for i, cDTO := range reqDTO.Convidados {
		convidadoID := uuid.Nil
		if cDTO.ID != nil && *cDTO.ID != "" {
			parsedID, err := uuid.Parse(*cDTO.ID)
			if err != nil {
				web.RespondError(w, r, "DADOS_INVALIDOS", "ID de convidado inválido: "+*cDTO.ID, http.StatusBadRequest)
				return
			}
			convidadoID = parsedID
		}
		convidadosDominio[i] = domain.ConvidadoParaRevisao{
			ID:   convidadoID,
			Nome: cDTO.Nome,
		}
	}

	err = h.service.RevisarGrupo(r.Context(), grupoID, reqDTO.ChaveDeAcesso, convidadosDominio)
	// ... (Lógica de tratamento de erro similar aos outros handlers) ...
	// ...
	// Exemplo:
	if err != nil {
		if errors.Is(err, domain.ErrGrupoNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Grupo não encontrado.", http.StatusNotFound)
			return
		}
		// ... outros erros de negócio
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
