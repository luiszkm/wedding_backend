// file: internal/gift/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/gift/application"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/storage"
)

type GiftHandler struct {
	service        *application.GiftService
	storageService storage.FileStorage
}

func NewGiftHandler(service *application.GiftService, storageService storage.FileStorage) *GiftHandler {
	return &GiftHandler{service: service, storageService: storageService}
}

// RespondError writes a JSON error response with a given code, message, and HTTP status.
func RespondError(w http.ResponseWriter, r *http.Request, code string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

// Respond writes a JSON response with a given payload and HTTP status.
func Respond(w http.ResponseWriter, r *http.Request, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *GiftHandler) HandleCriarPresente(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		RespondError(w, r, "CORPO_GRANDE", "Requisição muito grande.", http.StatusBadRequest)
		return
	}

	presenteJSON := r.FormValue("presente")
	var reqDTO CriarPresenteRequestDTO
	if err := json.Unmarshal([]byte(presenteJSON), &reqDTO); err != nil {
		RespondError(w, r, "DADOS_INVALIDOS", "Os dados do presente estão malformados.", http.StatusBadRequest)
		return
	}

	var fotoFinalURL string
	file, fileHeader, err := r.FormFile("foto")
	if err != nil && err != http.ErrMissingFile {
		RespondError(w, r, "ERRO_ARQUIVO", "Erro ao processar o arquivo.", http.StatusBadRequest)
		return
	}

	if file != nil {
		defer file.Close()
		uploadedURL, _, err := h.storageService.Upload(r.Context(), file, fileHeader.Header.Get("Content-Type"), fileHeader.Size)
		if err != nil {
			log.Printf("ERRO de upload: %v", err)
			RespondError(w, r, "UPLOAD_FALHOU", "Não foi possível enviar a imagem.", http.StatusInternalServerError)
			return
		}
		fotoFinalURL = uploadedURL
	} else {
		fotoFinalURL = reqDTO.FotoURL
	}

	detalhesDominio := domain.DetalhesPresente{
		Tipo:       reqDTO.Detalhes.Tipo,
		LinkDaLoja: reqDTO.Detalhes.LinkDaLoja,
	}

	novoPresente, err := h.service.CriarNovoPresente(r.Context(), idCasamento, reqDTO.Nome,
		reqDTO.Descricao, fotoFinalURL, reqDTO.EhFavorito, reqDTO.Categoria,
		detalhesDominio)
	if err != nil {
		// ... tratamento de erros de negócio ...
		RespondError(w, r, "ERRO_INTERNO", "Falha ao criar presente.", http.StatusInternalServerError)
		return
	}

	Respond(w, r, CriarPresenteResponseDTO{IDPresente: novoPresente.ID().String()}, http.StatusCreated)
}

func (h *GiftHandler) HandleListarPresentesPublicos(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	// 1. Chama o serviço de aplicação
	presentes, err := h.service.ListarPresentesDisponiveis(r.Context(), idCasamento)
	if err != nil {
		log.Printf("ERRO: %v\n", err)
		RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a lista de presentes.", http.StatusInternalServerError)
		return
	}

	// 2. Mapeia os objetos de domínio para DTOs de resposta
	respDTO := make([]PresentePublicoDTO, len(presentes))
	for i, p := range presentes {
		respDTO[i] = PresentePublicoDTO{
			ID:         p.ID().String(),
			Nome:       p.Nome(),
			Descricao:  p.Descricao(),
			FotoURL:    p.FotoURL(),
			EhFavorito: p.EhFavorito(),
			Categoria:  p.Categoria(),
			Detalhes: DetalhesPresenteDTO{
				Tipo:       p.Detalhes().Tipo,
				LinkDaLoja: p.Detalhes().LinkDaLoja,
			},
		}
	}

	// 3. Responde com sucesso
	Respond(w, r, respDTO, http.StatusOK)
}

func (h *GiftHandler) HandleFinalizarSelecao(w http.ResponseWriter, r *http.Request) {
	var reqDTO FinalizarSelecaoRequestDTO
	// ... decodificar corpo da requisição ...

	ids, err := parseUUIDs(reqDTO.IDsDosPresentes)
	if err != nil {
		// ... responder com erro 400 ...
	}

	selecao, err := h.service.FinalizarSelecaoDePresentes(r.Context(), reqDTO.ChaveDeAcesso, ids)
	if err != nil {
		var conflitoErr *domain.ErrPresentesConflitantes
		if errors.As(err, &conflitoErr) {
			respConflito := ConflitoSelecaoDTO{
				Codigo:                "CONFLITO_DE_SELECAO",
				Mensagem:              "Um ou mais itens na sua lista já foram selecionados.",
				PresentesConflitantes: stringUUIDs(conflitoErr.PresentesIDs),
			}
			Respond(w, r, respConflito, http.StatusConflict)
			return
		}
		// ... tratar outros erros (404 para chave de acesso, 500, etc) ...
		return
	}

	presentesDTO := make([]PresenteConfirmadoDTO, len(selecao.PresentesConfirmados()))
	for i, p := range selecao.PresentesConfirmados() {
		presentesDTO[i] = PresenteConfirmadoDTO{ID: p.ID.String(), Nome: p.Nome}
	}

	respDTO := SelecaoConfirmadaDTO{
		IDSelecao:            selecao.ID().String(),
		Mensagem:             "Sua seleção foi confirmada com sucesso. Obrigado!",
		PresentesConfirmados: presentesDTO,
	}
	Respond(w, r, respDTO, http.StatusCreated)
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	// Pré-aloca o slice com a capacidade necessária para evitar realocações.
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			// Se qualquer ID for inválido, a operação inteira falha.
			return nil, fmt.Errorf("id inválido na lista: %s", idStr)
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

// stringUUIDs converte um slice de UUIDs de volta para um slice de strings.
// Usado para construir a resposta de erro em caso de conflito.
func stringUUIDs(uuids []uuid.UUID) []string {
	ids := make([]string, len(uuids))
	for i, u := range uuids {
		ids[i] = u.String()
	}
	return ids
}
