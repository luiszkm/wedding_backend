// file: internal/gallery/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/gallery/application"
	"github.com/luiszkm/wedding_backend/internal/gallery/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type GalleryHandler struct {
	service *application.GalleryService
}

func NewGalleryHandler(service *application.GalleryService) *GalleryHandler {
	return &GalleryHandler{service: service}
}

func (h *GalleryHandler) HandleFazerUpload(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento")) // Assumindo que o ID do casamento virá do token/sessão no futuro
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil { // Limite de 50MB
		web.RespondError(w, r, "CORPO_GRANDE", "Requisição muito grande.", http.StatusBadRequest)
		return
	}

	// Pega o campo de texto 'rotulo'
	rotulo := r.FormValue("rotulo")

	// Pega os múltiplos arquivos do campo 'imagens[]'
	files := r.MultipartForm.File["imagens[]"]
	if len(files) == 0 {
		web.RespondError(w, r, "ARQUIVO_AUSENTE", "Nenhum arquivo de imagem foi enviado.", http.StatusBadRequest)
		return
	}

	ids, err := h.service.FazerUploadDeFotos(r.Context(), idCasamento, rotulo, files)
	if err != nil {
		log.Printf("ERRO no upload de fotos: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar upload.", http.StatusInternalServerError)
		return
	}

	respDTO := UploadFotosResponseDTO{
		IDsDasFotosCriadas: stringUUIDs(ids), // reusar a função helper de conversão
	}

	web.Respond(w, r, respDTO, http.StatusCreated)
}
func (h *GalleryHandler) HandleListarFotosPublicas(w http.ResponseWriter, r *http.Request) {
	idCasamento, err := uuid.Parse(chi.URLParam(r, "idCasamento"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do casamento é inválido.", http.StatusBadRequest)
		return
	}

	// Lê o parâmetro de consulta opcional '?rotulo='
	filtroRotulo := r.URL.Query().Get("rotulo")

	fotos, err := h.service.ListarFotosPublicas(r.Context(), idCasamento, filtroRotulo)
	if err != nil {
		// Tratar erros de negócio como 'rótulo inválido'
		if errors.Is(err, domain.ErrRotuloInvalido) {
			web.RespondError(w, r, "PARAMETRO_INVALIDO", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao listar fotos públicas: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a galeria.", http.StatusInternalServerError)
		return
	}

	// Mapeia os objetos de domínio para DTOs
	respDTO := make([]FotoPublicaDTO, len(fotos))
	for i, f := range fotos {
		// Converte os rótulos do tipo domain.Rotulo para string
		stringRotulos := make([]string, len(f.Rotulos()))
		for j, r := range f.Rotulos() {
			stringRotulos[j] = string(r)
		}
		respDTO[i] = FotoPublicaDTO{
			ID:         f.ID().String(),
			URLPublica: f.URLPublica(),
			EhFavorito: f.EhFavorito(),
			Rotulos:    stringRotulos,
		}
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}

func (h *GalleryHandler) HandleAlternarFavorito(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	fotoID, err := uuid.Parse(chi.URLParam(r, "idFoto"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID da foto é inválido.", http.StatusBadRequest)
		return
	}
	fotoAtualizada, err := h.service.AlternarFavoritoFoto(r.Context(), userID, fotoID)

	if err != nil {
		if errors.Is(err, domain.ErrFotoNaoEncontrada) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Foto não encontrada.", http.StatusNotFound)
			return
		}
		log.Printf("ERRO ao alternar favorito: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	// Reutiliza o DTO público para a resposta
	stringRotulos := make([]string, len(fotoAtualizada.Rotulos()))
	for i, rotulo := range fotoAtualizada.Rotulos() {
		stringRotulos[i] = string(rotulo)
	}
	respDTO := FotoPublicaDTO{
		ID:         fotoAtualizada.ID().String(),
		URLPublica: fotoAtualizada.URLPublica(),
		EhFavorito: fotoAtualizada.EhFavorito(),
		Rotulos:    stringRotulos,
	}

	web.Respond(w, r, respDTO, http.StatusOK)
}
func (h *GalleryHandler) HandleAdicionarRotulo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	fotoID, err := uuid.Parse(chi.URLParam(r, "idFoto"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID da foto é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO AdicionarRotuloRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "Corpo da requisição malformado.", http.StatusBadRequest)
		return
	}

	err = h.service.AdicionarRotulo(r.Context(), userID, fotoID, reqDTO.NomeDoRotulo)
	if err != nil {
		// Adicionar tratamento para erros de negócio (foto não encontrada, rótulo inválido)
		log.Printf("ERRO ao adicionar rótulo: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GalleryHandler) HandleRemoverRotulo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	fotoID, err := uuid.Parse(chi.URLParam(r, "idFoto"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID da foto é inválido.", http.StatusBadRequest)
		return
	}
	nomeRotulo := chi.URLParam(r, "nomeDoRotulo")

	err = h.service.RemoverRotulo(r.Context(), userID, fotoID, nomeRotulo)
	if err != nil {
		// Adicionar tratamento para erros de negócio
		log.Printf("ERRO ao remover rótulo: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *GalleryHandler) HandleDeletarFoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}
	fotoID, err := uuid.Parse(chi.URLParam(r, "idFoto"))
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID da foto é inválido.", http.StatusBadRequest)
		return
	}

	err = h.service.DeletarFoto(r.Context(), userID, fotoID)
	if err != nil {
		if errors.Is(err, domain.ErrFotoNaoEncontrada) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Foto não encontrada.", http.StatusNotFound)
			return
		}
		log.Printf("ERRO ao deletar foto: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func stringUUIDs(uuids []uuid.UUID) []string {
	ids := make([]string, len(uuids))
	for i, u := range uuids {
		ids[i] = u.String()
	}
	return ids
}
