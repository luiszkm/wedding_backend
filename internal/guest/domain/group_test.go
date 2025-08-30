// file: internal/guest/domain/group_test.go
package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGrupoDeConvidados(t *testing.T) {
	idCasamento := uuid.New()

	t.Run("deve criar um grupo com sucesso com dados válidos", func(t *testing.T) {
		nomes := []string{"Convidado 1", "Convidado 2"}
		chave := "padrinhos"

		grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)

		assert.NoError(t, err)
		assert.NotNil(t, grupo)
		assert.Equal(t, chave, grupo.ChaveDeAcesso())
		assert.Equal(t, idCasamento, grupo.IDCasamento())
		assert.Len(t, grupo.Convidados(), 2)
	})

	t.Run("deve retornar erro se a chave de acesso for vazia", func(t *testing.T) {
		nomes := []string{"Convidado 1"}

		_, err := NewGrupoDeConvidados(idCasamento, "", nomes)

		assert.Error(t, err)
		assert.Equal(t, ErrChaveDeAcessoObrigatoria, err)
	})

	t.Run("deve retornar erro se não houver nomes de convidados", func(t *testing.T) {
		nomes := []string{}
		chave := "padrinhos"

		_, err := NewGrupoDeConvidados(idCasamento, chave, nomes)

		assert.Error(t, err)
		assert.Equal(t, ErrPeloMenosUmConvidado, err)
	})
}

func TestGrupoDeConvidados_ConfirmarPresenca(t *testing.T) {
	idCasamento := uuid.New()
	nomes := []string{"João", "Maria"}
	chave := "familia"

	grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)
	assert.NoError(t, err)

	convidados := grupo.Convidados()
	convidado1ID := convidados[0].ID()
	convidado2ID := convidados[1].ID()

	t.Run("deve confirmar presença com sucesso para múltiplos convidados", func(t *testing.T) {
		respostas := []RespostaRSVP{
			{ConvidadoID: convidado1ID, Status: StatusRSVPConfirmado},
			{ConvidadoID: convidado2ID, Status: StatusRSVPRecusado},
		}

		err := grupo.ConfirmarPresenca(respostas)

		assert.NoError(t, err)
		assert.Equal(t, StatusRSVPConfirmado, grupo.Convidados()[0].StatusRSVP())
		assert.Equal(t, StatusRSVPRecusado, grupo.Convidados()[1].StatusRSVP())
	})

	t.Run("deve retornar erro para status inválido", func(t *testing.T) {
		respostas := []RespostaRSVP{
			{ConvidadoID: convidado1ID, Status: "INVALIDO"},
		}

		err := grupo.ConfirmarPresenca(respostas)

		assert.Error(t, err)
		assert.Equal(t, ErrStatusRSVPInvalido, err)
	})

	t.Run("deve retornar erro para convidado que não pertence ao grupo", func(t *testing.T) {
		convidadoInexistente := uuid.New()
		respostas := []RespostaRSVP{
			{ConvidadoID: convidadoInexistente, Status: StatusRSVPConfirmado},
		}

		err := grupo.ConfirmarPresenca(respostas)

		assert.Error(t, err)
		assert.Equal(t, ErrConvidadoNaoEncontradoNoGrupo, err)
	})
}

func TestGrupoDeConvidados_Revisar(t *testing.T) {
	idCasamento := uuid.New()
	nomes := []string{"Ana", "Pedro"}
	chave := "amigos"

	grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)
	assert.NoError(t, err)

	convidadoOriginalID := grupo.Convidados()[0].ID()

	t.Run("deve revisar grupo com sucesso mantendo convidados existentes", func(t *testing.T) {
		novaChave := "amigos-proximos"
		convidadosRevisao := []ConvidadoParaRevisao{
			{ID: convidadoOriginalID, Nome: "Ana Silva"}, // Renomeando
			{ID: uuid.Nil, Nome: "Carlos"},               // Novo convidado
		}

		err := grupo.Revisar(novaChave, convidadosRevisao)

		assert.NoError(t, err)
		assert.Equal(t, novaChave, grupo.ChaveDeAcesso())
		assert.Len(t, grupo.Convidados(), 2)
		assert.Equal(t, "Ana Silva", grupo.Convidados()[0].Nome())
		assert.Equal(t, "Carlos", grupo.Convidados()[1].Nome())
	})

	t.Run("deve retornar erro se chave de acesso for vazia", func(t *testing.T) {
		convidadosRevisao := []ConvidadoParaRevisao{
			{ID: uuid.Nil, Nome: "Teste"},
		}

		err := grupo.Revisar("", convidadosRevisao)

		assert.Error(t, err)
		assert.Equal(t, ErrChaveDeAcessoObrigatoria, err)
	})

	t.Run("deve retornar erro se não houver convidados", func(t *testing.T) {
		convidadosRevisao := []ConvidadoParaRevisao{}

		err := grupo.Revisar("nova-chave", convidadosRevisao)

		assert.Error(t, err)
		assert.Equal(t, ErrPeloMenosUmConvidado, err)
	})

	t.Run("deve retornar erro se tentar editar convidado que não pertence ao grupo", func(t *testing.T) {
		convidadoInexistente := uuid.New()
		convidadosRevisao := []ConvidadoParaRevisao{
			{ID: convidadoInexistente, Nome: "Inexistente"},
		}

		err := grupo.Revisar("nova-chave", convidadosRevisao)

		assert.Error(t, err)
		assert.Equal(t, ErrConvidadoNaoEncontradoNoGrupo, err)
	})
}

func TestGrupoDeConvidados_PodeSerRemovido(t *testing.T) {
	idCasamento := uuid.New()
	nomes := []string{"João", "Maria"}
	chave := "familia"

	grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)
	assert.NoError(t, err)

	t.Run("deve permitir remoção quando todos convidados estão pendentes", func(t *testing.T) {
		err := grupo.PodeSerRemovido()
		assert.NoError(t, err)
	})

	t.Run("deve impedir remoção quando há convidados confirmados", func(t *testing.T) {
		// Confirmar um convidado
		convidados := grupo.Convidados()
		respostas := []RespostaRSVP{
			{ConvidadoID: convidados[0].ID(), Status: StatusRSVPConfirmado},
		}

		err := grupo.ConfirmarPresenca(respostas)
		assert.NoError(t, err)

		// Agora deve impedir remoção
		err = grupo.PodeSerRemovido()
		assert.Error(t, err)
		assert.Equal(t, ErrNaoPodeRemoverGrupoComRSVP, err)
	})

	t.Run("deve impedir remoção quando há convidados recusados", func(t *testing.T) {
		// Criar novo grupo
		grupo2, err := NewGrupoDeConvidados(idCasamento, "outro-grupo", []string{"Pedro"})
		assert.NoError(t, err)

		// Recusar convidado
		convidados := grupo2.Convidados()
		respostas := []RespostaRSVP{
			{ConvidadoID: convidados[0].ID(), Status: StatusRSVPRecusado},
		}

		err = grupo2.ConfirmarPresenca(respostas)
		assert.NoError(t, err)

		// Agora deve impedir remoção
		err = grupo2.PodeSerRemovido()
		assert.Error(t, err)
		assert.Equal(t, ErrNaoPodeRemoverGrupoComRSVP, err)
	})
}
