package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/guest/application"
	"github.com/luiszkm/wedding_backend/internal/guest/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) Save(ctx context.Context, group *domain.GrupoDeConvidados) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockGroupRepository) FindByAccessKey(ctx context.Context, accessKey string) (*domain.GrupoDeConvidados, error) {
	args := m.Called(ctx, accessKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GrupoDeConvidados), args.Error(1)
}

func (m *MockGroupRepository) Update(ctx context.Context, userID uuid.UUID, group *domain.GrupoDeConvidados) error {
	args := m.Called(ctx, userID, group)
	return args.Error(0)
}

func (m *MockGroupRepository) FindByID(ctx context.Context, userID, groupID uuid.UUID) (*domain.GrupoDeConvidados, error) {
	args := m.Called(ctx, userID, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GrupoDeConvidados), args.Error(1)
}

func (m *MockGroupRepository) UpdateRSVP(ctx context.Context, group *domain.GrupoDeConvidados) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockGroupRepository) FindAllByEventID(ctx context.Context, userID, eventID uuid.UUID, statusFilter string) ([]*domain.GrupoDeConvidados, error) {
	args := m.Called(ctx, userID, eventID, statusFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.GrupoDeConvidados), args.Error(1)
}

func (m *MockGroupRepository) Delete(ctx context.Context, userID, groupID uuid.UUID) error {
	args := m.Called(ctx, userID, groupID)
	return args.Error(0)
}

func (m *MockGroupRepository) GetRSVPStats(ctx context.Context, userID, eventID uuid.UUID) (*domain.RSVPStats, error) {
	args := m.Called(ctx, userID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RSVPStats), args.Error(1)
}

func TestGuestHandler_HandleCriarGrupoDeConvidados(t *testing.T) {
	mockRepo := &MockGroupRepository{}
	service := application.NewGuestService(mockRepo)
	handler := NewGuestHandler(service)

	eventoID := uuid.New()
	userID := uuid.New()

	t.Run("deve criar grupo com sucesso", func(t *testing.T) {
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.GrupoDeConvidados")).Return(nil).Once()

		requestBody := CriarGrupoRequestDTO{
			ChaveDeAcesso:      "padrinhos",
			NomesDosConvidados: []string{"João", "Maria"},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/eventos/"+eventoID.String()+"/grupos-de-convidados", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("idCasamento", eventoID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.HandleCriarGrupoDeConvidados(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response CriarGrupoResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.IDGrupo)

		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro para dados inválidos", func(t *testing.T) {
		requestBody := CriarGrupoRequestDTO{
			ChaveDeAcesso:      "",
			NomesDosConvidados: []string{},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/eventos/"+eventoID.String()+"/grupos-de-convidados", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("idCasamento", eventoID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.HandleCriarGrupoDeConvidados(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGuestHandler_HandleObterGrupoPorChaveDeAcesso(t *testing.T) {
	mockRepo := &MockGroupRepository{}
	service := application.NewGuestService(mockRepo)
	handler := NewGuestHandler(service)

	t.Run("deve retornar grupo por chave de acesso", func(t *testing.T) {
		eventoID := uuid.New()
		grupoID := uuid.New()
		convidado1ID := uuid.New()
		convidado2ID := uuid.New()

		convidados := []*domain.Convidado{
			domain.HydrateConvidado(convidado1ID, "João", "PENDENTE"),
			domain.HydrateConvidado(convidado2ID, "Maria", "PENDENTE"),
		}
		grupo := domain.HydrateGroup(grupoID, eventoID, "padrinhos", convidados)

		mockRepo.On("FindByAccessKey", mock.Anything, "padrinhos").Return(grupo, nil).Once()

		req := httptest.NewRequest("GET", "/acesso-convidado?chave=padrinhos", nil)
		w := httptest.NewRecorder()

		handler.HandleObterGrupoPorChaveDeAcesso(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GrupoParaConfirmacaoDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, grupoID.String(), response.IDGrupo)
		assert.Len(t, response.Convidados, 2)
		assert.Equal(t, "João", response.Convidados[0].Nome)
		assert.Equal(t, "Maria", response.Convidados[1].Nome)

		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando grupo não encontrado", func(t *testing.T) {
		mockRepo.On("FindByAccessKey", mock.Anything, "inexistente").Return(nil, domain.ErrGrupoNaoEncontrado).Once()

		req := httptest.NewRequest("GET", "/acesso-convidado?chave=inexistente", nil)
		w := httptest.NewRecorder()

		handler.HandleObterGrupoPorChaveDeAcesso(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestGuestHandler_HandleConfirmarPresenca(t *testing.T) {
	mockRepo := &MockGroupRepository{}
	service := application.NewGuestService(mockRepo)
	handler := NewGuestHandler(service)

	t.Run("deve confirmar presença com sucesso", func(t *testing.T) {
		eventoID := uuid.New()
		grupoID := uuid.New()
		convidado1ID := uuid.New()
		convidado2ID := uuid.New()

		convidados := []*domain.Convidado{
			domain.HydrateConvidado(convidado1ID, "João", "PENDENTE"),
			domain.HydrateConvidado(convidado2ID, "Maria", "PENDENTE"),
		}
		grupo := domain.HydrateGroup(grupoID, eventoID, "padrinhos", convidados)

		mockRepo.On("FindByAccessKey", mock.Anything, "padrinhos").Return(grupo, nil).Once()
		mockRepo.On("UpdateRSVP", mock.Anything, mock.AnythingOfType("*domain.GrupoDeConvidados")).Return(nil).Once()

		requestBody := ConfirmarPresencaRequestDTO{
			ChaveDeAcesso: "padrinhos",
			Respostas: []RespostaRSVPDTO{
				{IDConvidado: convidado1ID.String(), Status: "CONFIRMADO"},
				{IDConvidado: convidado2ID.String(), Status: "RECUSADO"},
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/rsvps", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandleConfirmarPresenca(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro para status inválido", func(t *testing.T) {
		eventoID := uuid.New()
		grupoID := uuid.New()
		convidado1ID := uuid.New()

		convidados := []*domain.Convidado{
			domain.HydrateConvidado(convidado1ID, "João", "PENDENTE"),
		}
		grupo := domain.HydrateGroup(grupoID, eventoID, "padrinhos", convidados)

		mockRepo.On("FindByAccessKey", mock.Anything, "padrinhos").Return(grupo, nil).Once()

		requestBody := ConfirmarPresencaRequestDTO{
			ChaveDeAcesso: "padrinhos",
			Respostas: []RespostaRSVPDTO{
				{IDConvidado: convidado1ID.String(), Status: "INVALIDO"},
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/rsvps", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandleConfirmarPresenca(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
