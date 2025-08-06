package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCota(t *testing.T) {
	idPresente := uuid.New()

	t.Run("deve criar cota com sucesso com dados válidos", func(t *testing.T) {
		numeroCota := 1
		valorCota := 100.50

		cota, err := NewCota(idPresente, numeroCota, valorCota)

		assert.NoError(t, err)
		assert.NotNil(t, cota)
		assert.Equal(t, idPresente, cota.IDPresente())
		assert.Equal(t, numeroCota, cota.NumeroCota())
		assert.Equal(t, valorCota, cota.ValorCota())
		assert.Equal(t, StatusCotaDisponivel, cota.Status())
		assert.True(t, cota.EstaDisponivel())
		assert.False(t, cota.EstaSelecionada())
		assert.Nil(t, cota.IDSelecao())
		assert.NotEqual(t, uuid.Nil, cota.ID())
	})

	t.Run("deve retornar erro se número da cota for inválido", func(t *testing.T) {
		_, err := NewCota(idPresente, 0, 100.50)
		assert.Error(t, err)
		assert.Equal(t, ErrNumeroCotaInvalido, err)

		_, err = NewCota(idPresente, -1, 100.50)
		assert.Error(t, err)
		assert.Equal(t, ErrNumeroCotaInvalido, err)
	})

	t.Run("deve retornar erro se valor da cota for inválido", func(t *testing.T) {
		_, err := NewCota(idPresente, 1, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrValorCotaInvalido, err)

		_, err = NewCota(idPresente, 1, -10.50)
		assert.Error(t, err)
		assert.Equal(t, ErrValorCotaInvalido, err)
	})
}

func TestCota_Selecionar(t *testing.T) {
	idPresente := uuid.New()
	cota, _ := NewCota(idPresente, 1, 100.50)
	idSelecao := uuid.New()

	t.Run("deve selecionar cota disponível com sucesso", func(t *testing.T) {
		err := cota.Selecionar(idSelecao)

		assert.NoError(t, err)
		assert.Equal(t, StatusCotaSelecionada, cota.Status())
		assert.False(t, cota.EstaDisponivel())
		assert.True(t, cota.EstaSelecionada())
		assert.NotNil(t, cota.IDSelecao())
		assert.Equal(t, idSelecao, *cota.IDSelecao())
	})

	t.Run("deve retornar erro ao tentar selecionar cota já selecionada", func(t *testing.T) {
		outroIdSelecao := uuid.New()
		err := cota.Selecionar(outroIdSelecao)

		assert.Error(t, err)
		assert.Equal(t, ErrCotaJaSelecionada, err)
		// Status e seleção anterior devem permanecer inalterados
		assert.Equal(t, StatusCotaSelecionada, cota.Status())
		assert.Equal(t, idSelecao, *cota.IDSelecao())
	})
}

func TestCota_LiberarSelecao(t *testing.T) {
	idPresente := uuid.New()
	cota, _ := NewCota(idPresente, 1, 100.50)
	idSelecao := uuid.New()

	t.Run("deve retornar erro ao tentar liberar cota já disponível", func(t *testing.T) {
		err := cota.LiberarSelecao()

		assert.Error(t, err)
		assert.Equal(t, ErrCotaJaDisponivel, err)
		assert.Equal(t, StatusCotaDisponivel, cota.Status())
	})

	t.Run("deve liberar seleção de cota selecionada com sucesso", func(t *testing.T) {
		// Primeiro selecionar a cota
		cota.Selecionar(idSelecao)
		assert.True(t, cota.EstaSelecionada())

		// Depois liberar
		err := cota.LiberarSelecao()

		assert.NoError(t, err)
		assert.Equal(t, StatusCotaDisponivel, cota.Status())
		assert.True(t, cota.EstaDisponivel())
		assert.False(t, cota.EstaSelecionada())
		assert.Nil(t, cota.IDSelecao())
	})
}

func TestHydrateCota(t *testing.T) {
	t.Run("deve hidratar cota corretamente", func(t *testing.T) {
		id := uuid.New()
		idPresente := uuid.New()
		numeroCota := 5
		valorCota := 250.75
		status := StatusCotaSelecionada
		idSelecao := uuid.New()

		cota := HydrateCota(id, idPresente, numeroCota, valorCota, status, &idSelecao)

		assert.Equal(t, id, cota.ID())
		assert.Equal(t, idPresente, cota.IDPresente())
		assert.Equal(t, numeroCota, cota.NumeroCota())
		assert.Equal(t, valorCota, cota.ValorCota())
		assert.Equal(t, status, cota.Status())
		assert.NotNil(t, cota.IDSelecao())
		assert.Equal(t, idSelecao, *cota.IDSelecao())
		assert.False(t, cota.EstaDisponivel())
		assert.True(t, cota.EstaSelecionada())
	})

	t.Run("deve hidratar cota disponível sem ID de seleção", func(t *testing.T) {
		id := uuid.New()
		idPresente := uuid.New()

		cota := HydrateCota(id, idPresente, 1, 100.0, StatusCotaDisponivel, nil)

		assert.Equal(t, StatusCotaDisponivel, cota.Status())
		assert.Nil(t, cota.IDSelecao())
		assert.True(t, cota.EstaDisponivel())
		assert.False(t, cota.EstaSelecionada())
	})
}
