package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewItineraryItem(t *testing.T) {
	eventID := uuid.New()
	horario := time.Now().Add(2 * time.Hour)
	titulo := "Cerimônia"
	descricao := "Cerimônia religiosa na igreja"

	tests := []struct {
		name        string
		eventID     uuid.UUID
		horario     time.Time
		titulo      string
		descricao   *string
		expectedErr error
	}{
		{
			name:        "deve criar item válido",
			eventID:     eventID,
			horario:     horario,
			titulo:      titulo,
			descricao:   &descricao,
			expectedErr: nil,
		},
		{
			name:        "deve falhar com evento nulo",
			eventID:     uuid.Nil,
			horario:     horario,
			titulo:      titulo,
			descricao:   &descricao,
			expectedErr: ErrEventoObrigatorio,
		},
		{
			name:        "deve falhar com horário zero",
			eventID:     eventID,
			horario:     time.Time{},
			titulo:      titulo,
			descricao:   &descricao,
			expectedErr: ErrHorarioObrigatorio,
		},
		{
			name:        "deve falhar com título vazio",
			eventID:     eventID,
			horario:     horario,
			titulo:      "",
			descricao:   &descricao,
			expectedErr: ErrTituloAtividadeObrigatorio,
		},
		{
			name:        "deve falhar com título apenas com espaços",
			eventID:     eventID,
			horario:     horario,
			titulo:      "   ",
			descricao:   &descricao,
			expectedErr: ErrTituloAtividadeObrigatorio,
		},
		{
			name:        "deve falhar com título muito longo",
			eventID:     eventID,
			horario:     horario,
			titulo:      string(make([]byte, 256)), // 256 caracteres
			descricao:   &descricao,
			expectedErr: ErrTituloAtividadeMuitoLongo,
		},
		{
			name:        "deve criar item com descrição nil",
			eventID:     eventID,
			horario:     horario,
			titulo:      titulo,
			descricao:   nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewItineraryItem(tt.eventID, tt.horario, tt.titulo, tt.descricao)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.NotEqual(t, uuid.Nil, item.ID())
				assert.Equal(t, tt.eventID, item.IDEvento())
				assert.Equal(t, tt.horario, item.Horario())
				assert.Equal(t, tt.titulo, item.TituloAtividade())
				assert.Equal(t, tt.descricao, item.DescricaoAtividade())
				assert.False(t, item.CreatedAt().IsZero())
				assert.False(t, item.UpdatedAt().IsZero())
			}
		})
	}
}

func TestItineraryItem_Update(t *testing.T) {
	eventID := uuid.New()
	horarioOriginal := time.Now().Add(2 * time.Hour)
	tituloOriginal := "Cerimônia Original"
	descricaoOriginal := "Descrição original"

	novoHorario := time.Now().Add(3 * time.Hour)
	novoTitulo := "Cerimônia Atualizada"
	novaDescricao := "Nova descrição"

	tests := []struct {
		name        string
		horario     time.Time
		titulo      string
		descricao   *string
		expectedErr error
	}{
		{
			name:        "deve atualizar item válido",
			horario:     novoHorario,
			titulo:      novoTitulo,
			descricao:   &novaDescricao,
			expectedErr: nil,
		},
		{
			name:        "deve falhar com horário zero",
			horario:     time.Time{},
			titulo:      novoTitulo,
			descricao:   &novaDescricao,
			expectedErr: ErrHorarioObrigatorio,
		},
		{
			name:        "deve falhar com título vazio",
			horario:     novoHorario,
			titulo:      "",
			descricao:   &novaDescricao,
			expectedErr: ErrTituloAtividadeObrigatorio,
		},
		{
			name:        "deve falhar com título muito longo",
			horario:     novoHorario,
			titulo:      string(make([]byte, 256)),
			descricao:   &novaDescricao,
			expectedErr: ErrTituloAtividadeMuitoLongo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cria uma nova instância para cada teste
			item, err := NewItineraryItem(eventID, horarioOriginal, tituloOriginal, &descricaoOriginal)
			assert.NoError(t, err)
			
			updatedAtOriginal := item.UpdatedAt()
			time.Sleep(time.Millisecond) // Para garantir diferença no timestamp

			err = item.Update(tt.horario, tt.titulo, tt.descricao)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				// Valores devem permanecer os mesmos em caso de erro
				assert.Equal(t, horarioOriginal, item.Horario())
				assert.Equal(t, tituloOriginal, item.TituloAtividade())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.horario, item.Horario())
				assert.Equal(t, tt.titulo, item.TituloAtividade())
				assert.Equal(t, tt.descricao, item.DescricaoAtividade())
				assert.True(t, item.UpdatedAt().After(updatedAtOriginal))
			}
		})
	}
}

func TestHydrateItineraryItem(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	horario := time.Now().Add(2 * time.Hour)
	titulo := "Cerimônia"
	descricao := "Cerimônia religiosa"
	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()

	item := HydrateItineraryItem(id, eventID, horario, titulo, &descricao, createdAt, updatedAt)

	assert.NotNil(t, item)
	assert.Equal(t, id, item.ID())
	assert.Equal(t, eventID, item.IDEvento())
	assert.Equal(t, horario, item.Horario())
	assert.Equal(t, titulo, item.TituloAtividade())
	assert.Equal(t, &descricao, item.DescricaoAtividade())
	assert.Equal(t, createdAt, item.CreatedAt())
	assert.Equal(t, updatedAt, item.UpdatedAt())
}