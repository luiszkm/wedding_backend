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
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/storage"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type GiftHandler struct {
	service        *application.GiftService
	storageService storage.FileStorage
}

func NewGiftHandler(service *application.GiftService, storageService storage.FileStorage) *GiftHandler {
	return &GiftHandler{service: service, storageService: storageService}
}

func (h *GiftHandler) HandleCriarPresente(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	idEvento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		web.RespondError(w, r, "CORPO_GRANDE", "Requisição muito grande.", http.StatusBadRequest)
		return
	}

	presenteJSON := r.FormValue("presente")
	var reqDTO CriarPresenteRequestDTO
	if err := json.Unmarshal([]byte(presenteJSON), &reqDTO); err != nil {
		web.RespondError(w, r, "DADOS_INVALIDOS", "Os dados do presente estão malformados.", http.StatusBadRequest)
		return
	}

	// Validar campos para tipo de presente
	if reqDTO.Tipo != domain.TipoPresenteIntegral && reqDTO.Tipo != domain.TipoPresenteFracionado {
		web.RespondError(w, r, "TIPO_INVALIDO", "Tipo de presente deve ser INTEGRAL ou FRACIONADO.", http.StatusBadRequest)
		return
	}

	if reqDTO.Tipo == domain.TipoPresenteFracionado {
		if reqDTO.ValorTotal == nil || reqDTO.NumeroCotas == nil {
			web.RespondError(w, r, "CAMPOS_OBRIGATORIOS", "Valor total e número de cotas são obrigatórios para presentes fracionados.", http.StatusBadRequest)
			return
		}
	}

	var fotoFinalURL string
	file, fileHeader, err := r.FormFile("foto")
	if err != nil && err != http.ErrMissingFile {
		web.RespondError(w, r, "ERRO_ARQUIVO", "Erro ao processar o arquivo.", http.StatusBadRequest)
		return
	}

	if file != nil {
		defer file.Close()
		uploadedURL, _, err := h.storageService.Upload(r.Context(), file, fileHeader.Header.Get("Content-Type"), fileHeader.Size)
		if err != nil {
			log.Printf("ERRO de upload: %v", err)
			web.RespondError(w, r, "UPLOAD_FALHOU", "Não foi possível enviar a imagem.", http.StatusInternalServerError)
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

	var novoPresente *domain.Presente

	if reqDTO.Tipo == domain.TipoPresenteIntegral {
		novoPresente, err = h.service.CriarPresenteIntegral(
			r.Context(),
			userID,
			idEvento,
			reqDTO.Nome,
			reqDTO.Descricao,
			fotoFinalURL,
			reqDTO.Categoria,
			reqDTO.EhFavorito,
			detalhesDominio,
		)
	} else {
		novoPresente, err = h.service.CriarPresenteFracionado(
			r.Context(),
			userID,
			idEvento,
			reqDTO.Nome,
			reqDTO.Descricao,
			fotoFinalURL,
			reqDTO.Categoria,
			reqDTO.EhFavorito,
			detalhesDominio,
			*reqDTO.ValorTotal,
			*reqDTO.NumeroCotas,
		)
	}

	if err != nil {
		if errors.Is(err, domain.ErrNomePresenteObrigatorio) ||
			errors.Is(err, domain.ErrDetalhesInvalidos) ||
			errors.Is(err, domain.ErrValorTotalInvalido) ||
			errors.Is(err, domain.ErrNumeroCotasInvalido) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao criar presente: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao criar presente.", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, CriarPresenteResponseDTO{IDPresente: novoPresente.ID().String()}, http.StatusCreated)
}

func (h *GiftHandler) HandleListarPresentesPublicos(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	presentes, err := h.service.ListarPresentesDisponiveis(r.Context(), idCasamento)
	if err != nil {
		log.Printf("ERRO: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a lista de presentes.", http.StatusInternalServerError)
		return
	}

	respDTO := make([]PresentePublicoDTO, len(presentes))
	for i, p := range presentes {
		dto := PresentePublicoDTO{
			ID:         p.ID().String(),
			Nome:       p.Nome(),
			Descricao:  p.Descricao(),
			FotoURL:    p.FotoURL(),
			EhFavorito: p.EhFavorito(),
			Categoria:  p.Categoria(),
			Tipo:       p.Tipo(),
			Status:     p.Status(),
			Detalhes: DetalhesPresenteDTO{
				Tipo:       p.Detalhes().Tipo,
				LinkDaLoja: p.Detalhes().LinkDaLoja,
			},
		}

		if p.EhFracionado() {
			valorTotal := p.ValorTotal()
			valorCota := p.ObterValorCota()
			cotasTotais := len(p.Cotas())
			cotasDisponiveis := p.ContarCotasDisponiveis()
			cotasSelecionadas := p.ContarCotasSelecionadas()

			dto.ValorTotal = valorTotal
			dto.ValorCota = &valorCota
			dto.CotasTotais = &cotasTotais
			dto.CotasDisponiveis = &cotasDisponiveis
			dto.CotasSelecionadas = &cotasSelecionadas
		}

		respDTO[i] = dto
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}

func (h *GiftHandler) HandleFinalizarSelecao(w http.ResponseWriter, r *http.Request) {
	var reqDTO FinalizarSelecaoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	if len(reqDTO.Itens) == 0 {
		web.RespondError(w, r, "LISTA_VAZIA", "A lista de presentes não pode estar vazia.", http.StatusBadRequest)
		return
	}

	// Converter DTOs para domain objects
	itens := make([]application.ItemSelecao, len(reqDTO.Itens))
	for i, item := range reqDTO.Itens {
		idPresente, err := uuid.Parse(item.IDPresente)
		if err != nil {
			web.RespondError(w, r, "ID_INVALIDO", fmt.Sprintf("ID do presente inválido: %s", item.IDPresente), http.StatusBadRequest)
			return
		}

		if item.Quantidade <= 0 {
			web.RespondError(w, r, "QUANTIDADE_INVALIDA", "A quantidade deve ser positiva.", http.StatusBadRequest)
			return
		}

		itens[i] = application.ItemSelecao{
			IDPresente: idPresente,
			Quantidade: item.Quantidade,
		}
	}

	selecao, err := h.service.FinalizarSelecaoDePresentes(r.Context(), reqDTO.ChaveDeAcesso, itens)
	if err != nil {
		var conflitoErr *domain.ErrPresentesConflitantes
		if errors.As(err, &conflitoErr) {
			respConflito := ConflitoSelecaoDTO{
				Codigo:                "CONFLITO_DE_SELECAO",
				Mensagem:              "Um ou mais itens na sua lista já foram selecionados.",
				PresentesConflitantes: stringUUIDs(conflitoErr.PresentesIDs),
			}
			web.Respond(w, r, respConflito, http.StatusConflict)
			return
		}

		log.Printf("ERRO ao finalizar seleção: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao finalizar seleção.", http.StatusInternalServerError)
		return
	}

	presentesDTO := make([]PresenteConfirmadoDTO, len(selecao.PresentesConfirmados()))
	for i, p := range selecao.PresentesConfirmados() {
		dto := PresenteConfirmadoDTO{
			ID:         p.ID.String(),
			Nome:       p.Nome,
			Quantidade: p.Quantidade,
		}

		if p.ValorCota != nil {
			valorTotal := *p.ValorCota * float64(p.Quantidade)
			dto.ValorCota = p.ValorCota
			dto.ValorTotal = &valorTotal
		}

		presentesDTO[i] = dto
	}

	respDTO := SelecaoConfirmadaDTO{
		IDSelecao:            selecao.ID().String(),
		Mensagem:             "Sua seleção foi confirmada com sucesso. Obrigado!",
		ValorTotal:           selecao.CalcularValorTotal(),
		PresentesConfirmados: presentesDTO,
	}

	web.Respond(w, r, respDTO, http.StatusCreated)
}

// Método legacy para compatibilidade
func (h *GiftHandler) HandleFinalizarSelecaoLegacy(w http.ResponseWriter, r *http.Request) {
	var reqDTO FinalizarSelecaoLegacyRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	// Converter para o novo formato
	itens := make([]application.ItemSelecao, len(reqDTO.IDsDosPresentes))
	for i, idStr := range reqDTO.IDsDosPresentes {
		idPresente, err := uuid.Parse(idStr)
		if err != nil {
			web.RespondError(w, r, "ID_INVALIDO", fmt.Sprintf("ID do presente inválido: %s", idStr), http.StatusBadRequest)
			return
		}

		itens[i] = application.ItemSelecao{
			IDPresente: idPresente,
			Quantidade: 1, // Legacy sempre usa quantidade 1
		}
	}

	// Usar novo método
	novoReqDTO := FinalizarSelecaoRequestDTO{
		ChaveDeAcesso: reqDTO.ChaveDeAcesso,
		Itens:         make([]ItemSelecaoDTO, len(itens)),
	}

	for i, item := range itens {
		novoReqDTO.Itens[i] = ItemSelecaoDTO{
			IDPresente: item.IDPresente.String(),
			Quantidade: item.Quantidade,
		}
	}

	// Reprocessar com novo handler
	reqBody, _ := json.Marshal(novoReqDTO)
	r.Body = http.NoBody
	r.ContentLength = int64(len(reqBody))

	h.HandleFinalizarSelecao(w, r)
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("id inválido na lista: %s", idStr)
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

func stringUUIDs(uuids []uuid.UUID) []string {
	ids := make([]string, len(uuids))
	for i, u := range uuids {
		ids[i] = u.String()
	}
	return ids
}
