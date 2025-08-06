package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPresenteIntegral(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{
		Tipo:       TipoDetalheProdutoExterno,
		LinkDaLoja: "https://exemplo.com/produto",
	}

	t.Run("deve criar presente integral com sucesso", func(t *testing.T) {
		presente, err := NewPresenteIntegral(idCasamento, "Cafeteira", "Descrição", "https://foto.com", true, "COZINHA", detalhes)

		assert.NoError(t, err)
		assert.NotNil(t, presente)
		assert.Equal(t, idCasamento, presente.IDCasamento())
		assert.Equal(t, "Cafeteira", presente.Nome())
		assert.Equal(t, TipoPresenteIntegral, presente.Tipo())
		assert.True(t, presente.EhIntegral())
		assert.False(t, presente.EhFracionado())
		assert.Equal(t, StatusDisponivel, presente.Status())
		assert.Nil(t, presente.ValorTotal())
		assert.Nil(t, presente.Cotas())
		assert.Equal(t, 0, presente.ContarCotasDisponiveis())
	})

	t.Run("deve retornar erro se nome for vazio", func(t *testing.T) {
		_, err := NewPresenteIntegral(idCasamento, "", "Descrição", "https://foto.com", false, "COZINHA", detalhes)
		assert.Error(t, err)
		assert.Equal(t, ErrNomePresenteObrigatorio, err)
	})

	t.Run("deve retornar erro se detalhes forem inválidos", func(t *testing.T) {
		detalhesInvalidos := DetalhesPresente{
			Tipo:       TipoDetalheProdutoExterno,
			LinkDaLoja: "", // Link vazio para produto externo
		}
		_, err := NewPresenteIntegral(idCasamento, "Presente", "Desc", "", false, "CAT", detalhesInvalidos)
		assert.Error(t, err)
		assert.Equal(t, ErrDetalhesInvalidos, err)
	})
}

func TestNewPresenteFracionado(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{
		Tipo:       TipoDetalheProdutoExterno,
		LinkDaLoja: "https://exemplo.com/produto",
	}

	t.Run("deve criar presente fracionado com sucesso", func(t *testing.T) {
		valorTotal := 1000.0
		numeroCotas := 5

		presente, err := NewPresenteFracionado(idCasamento, "Geladeira", "Descrição", "https://foto.com", true, "COZINHA", detalhes, valorTotal, numeroCotas)

		assert.NoError(t, err)
		assert.NotNil(t, presente)
		assert.Equal(t, "Geladeira", presente.Nome())
		assert.Equal(t, TipoPresenteFracionado, presente.Tipo())
		assert.False(t, presente.EhIntegral())
		assert.True(t, presente.EhFracionado())
		assert.Equal(t, StatusDisponivel, presente.Status())
		assert.NotNil(t, presente.ValorTotal())
		assert.Equal(t, valorTotal, *presente.ValorTotal())
		assert.NotNil(t, presente.Cotas())
		assert.Len(t, presente.Cotas(), numeroCotas)
		assert.Equal(t, numeroCotas, presente.ContarCotasDisponiveis())
		assert.Equal(t, 0, presente.ContarCotasSelecionadas())

		// Verificar valor das cotas
		valorCotaEsperado := 200.0 // 1000 / 5
		assert.Equal(t, valorCotaEsperado, presente.ObterValorCota())

		// Verificar que todas as cotas foram criadas corretamente
		for i, cota := range presente.Cotas() {
			assert.Equal(t, presente.ID(), cota.IDPresente())
			assert.Equal(t, i+1, cota.NumeroCota())
			assert.Equal(t, valorCotaEsperado, cota.ValorCota())
			assert.True(t, cota.EstaDisponivel())
		}
	})

	t.Run("deve retornar erro se valor total for inválido", func(t *testing.T) {
		_, err := NewPresenteFracionado(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes, 0, 5)
		assert.Error(t, err)
		assert.Equal(t, ErrValorTotalInvalido, err)

		_, err = NewPresenteFracionado(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes, -100, 5)
		assert.Error(t, err)
		assert.Equal(t, ErrValorTotalInvalido, err)
	})

	t.Run("deve retornar erro se número de cotas for inválido", func(t *testing.T) {
		_, err := NewPresenteFracionado(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes, 1000, 1)
		assert.Error(t, err)
		assert.Equal(t, ErrNumeroCotasInvalido, err)

		_, err = NewPresenteFracionado(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes, 1000, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNumeroCotasInvalido, err)
	})
}

func TestPresente_SelecionarCotas(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{Tipo: TipoDetalheProdutoExterno, LinkDaLoja: "https://test.com"}
	presente, _ := NewPresenteFracionado(idCasamento, "Geladeira", "Desc", "", false, "CAT", detalhes, 1000, 5)
	idSelecao := uuid.New()

	t.Run("deve selecionar cotas com sucesso", func(t *testing.T) {
		err := presente.SelecionarCotas(2, idSelecao)

		assert.NoError(t, err)
		assert.Equal(t, StatusParcialmenteSelecionado, presente.Status())
		assert.Equal(t, 3, presente.ContarCotasDisponiveis())
		assert.Equal(t, 2, presente.ContarCotasSelecionadas())

		// Verificar que as duas primeiras cotas foram selecionadas
		assert.True(t, presente.Cotas()[0].EstaSelecionada())
		assert.True(t, presente.Cotas()[1].EstaSelecionada())
		assert.True(t, presente.Cotas()[2].EstaDisponivel())
	})

	t.Run("deve retornar erro para presente integral", func(t *testing.T) {
		presenteIntegral, _ := NewPresenteIntegral(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes)
		err := presenteIntegral.SelecionarCotas(1, idSelecao)

		assert.Error(t, err)
		assert.Equal(t, ErrPresenteNaoFracionado, err)
	})

	t.Run("deve retornar erro se não houver cotas suficientes", func(t *testing.T) {
		err := presente.SelecionarCotas(4, uuid.New()) // Apenas 3 disponíveis

		assert.Error(t, err)
		assert.Equal(t, ErrCotasIndisponiveis, err)
	})

	t.Run("deve selecionar todas as cotas restantes", func(t *testing.T) {
		err := presente.SelecionarCotas(3, uuid.New())

		assert.NoError(t, err)
		assert.Equal(t, StatusSelecionado, presente.Status())
		assert.Equal(t, 0, presente.ContarCotasDisponiveis())
		assert.Equal(t, 5, presente.ContarCotasSelecionadas())
	})

	t.Run("deve retornar erro se presente já estiver completamente selecionado", func(t *testing.T) {
		err := presente.SelecionarCotas(1, uuid.New())

		assert.Error(t, err)
		assert.Equal(t, ErrPresenteJaSelecionado, err)
	})
}

func TestPresente_SelecionarIntegral(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{Tipo: TipoDetalheProdutoExterno, LinkDaLoja: "https://test.com"}
	presente, _ := NewPresenteIntegral(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes)
	idSelecao := uuid.New()

	t.Run("deve selecionar presente integral com sucesso", func(t *testing.T) {
		err := presente.SelecionarIntegral(idSelecao)

		assert.NoError(t, err)
		assert.Equal(t, StatusSelecionado, presente.Status())
	})

	t.Run("deve retornar erro para presente fracionado", func(t *testing.T) {
		presenteFracionado, _ := NewPresenteFracionado(idCasamento, "Geladeira", "Desc", "", false, "CAT", detalhes, 1000, 5)
		err := presenteFracionado.SelecionarIntegral(idSelecao)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operação válida apenas para presentes integrais")
	})

	t.Run("deve retornar erro se presente já estiver selecionado", func(t *testing.T) {
		err := presente.SelecionarIntegral(uuid.New())

		assert.Error(t, err)
		assert.Equal(t, ErrPresenteJaSelecionado, err)
	})
}

func TestPresente_LiberarSelecao(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{Tipo: TipoDetalheProdutoExterno, LinkDaLoja: "https://test.com"}

	t.Run("deve liberar seleção de presente integral", func(t *testing.T) {
		presente, _ := NewPresenteIntegral(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes)
		idSelecao := uuid.New()
		presente.SelecionarIntegral(idSelecao)

		err := presente.LiberarSelecao(idSelecao)

		assert.NoError(t, err)
		assert.Equal(t, StatusDisponivel, presente.Status())
	})

	t.Run("deve liberar seleção de cotas do presente fracionado", func(t *testing.T) {
		presente, _ := NewPresenteFracionado(idCasamento, "Geladeira", "Desc", "", false, "CAT", detalhes, 1000, 5)
		idSelecao := uuid.New()
		presente.SelecionarCotas(3, idSelecao)

		assert.Equal(t, StatusParcialmenteSelecionado, presente.Status())
		assert.Equal(t, 3, presente.ContarCotasSelecionadas())

		err := presente.LiberarSelecao(idSelecao)

		assert.NoError(t, err)
		assert.Equal(t, StatusDisponivel, presente.Status())
		assert.Equal(t, 5, presente.ContarCotasDisponiveis())
		assert.Equal(t, 0, presente.ContarCotasSelecionadas())
	})
}

func TestPresente_ContadorCotas(t *testing.T) {
	idCasamento := uuid.New()
	detalhes := DetalhesPresente{Tipo: TipoDetalheProdutoExterno, LinkDaLoja: "https://test.com"}

	t.Run("presente integral deve retornar 0 para contadores de cota", func(t *testing.T) {
		presente, _ := NewPresenteIntegral(idCasamento, "Presente", "Desc", "", false, "CAT", detalhes)

		assert.Equal(t, 0, presente.ContarCotasDisponiveis())
		assert.Equal(t, 0, presente.ContarCotasSelecionadas())

		// Depois de selecionar
		presente.SelecionarIntegral(uuid.New())
		assert.Equal(t, 1, presente.ContarCotasSelecionadas()) // Para integral, 1 quando selecionado
	})

	t.Run("presente fracionado deve contar cotas corretamente", func(t *testing.T) {
		presente, _ := NewPresenteFracionado(idCasamento, "Geladeira", "Desc", "", false, "CAT", detalhes, 1000, 10)

		// Inicial
		assert.Equal(t, 10, presente.ContarCotasDisponiveis())
		assert.Equal(t, 0, presente.ContarCotasSelecionadas())

		// Após selecionar algumas cotas
		presente.SelecionarCotas(4, uuid.New())
		assert.Equal(t, 6, presente.ContarCotasDisponiveis())
		assert.Equal(t, 4, presente.ContarCotasSelecionadas())

		// Após selecionar mais cotas
		presente.SelecionarCotas(3, uuid.New())
		assert.Equal(t, 3, presente.ContarCotasDisponiveis())
		assert.Equal(t, 7, presente.ContarCotasSelecionadas())
	})
}
