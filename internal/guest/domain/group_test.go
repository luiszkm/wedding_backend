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