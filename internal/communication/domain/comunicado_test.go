package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewComunicado(t *testing.T) {
	idEvento := uuid.New()

	t.Run("deve criar comunicado com sucesso com dados válidos", func(t *testing.T) {
		titulo := "Lembrete importante"
		mensagem := "A cerimônia começa às 19h. Não se atrasem!"

		comunicado, err := NewComunicado(idEvento, titulo, mensagem)

		assert.NoError(t, err)
		assert.NotNil(t, comunicado)
		assert.Equal(t, idEvento, comunicado.IDEvento())
		assert.Equal(t, titulo, comunicado.Titulo())
		assert.Equal(t, mensagem, comunicado.Mensagem())
		assert.NotEqual(t, uuid.Nil, comunicado.ID())
		assert.False(t, comunicado.DataPublicacao().IsZero())
	})

	t.Run("deve retornar erro se o título for vazio", func(t *testing.T) {
		mensagem := "Mensagem válida"

		_, err := NewComunicado(idEvento, "", mensagem)

		assert.Error(t, err)
		assert.Equal(t, ErrTituloObrigatorio, err)
	})

	t.Run("deve retornar erro se o título for apenas espaços", func(t *testing.T) {
		mensagem := "Mensagem válida"

		_, err := NewComunicado(idEvento, "   ", mensagem)

		assert.Error(t, err)
		assert.Equal(t, ErrTituloObrigatorio, err)
	})

	t.Run("deve retornar erro se o título for muito longo", func(t *testing.T) {
		titulo := strings.Repeat("a", 256) // 256 caracteres
		mensagem := "Mensagem válida"

		_, err := NewComunicado(idEvento, titulo, mensagem)

		assert.Error(t, err)
		assert.Equal(t, ErrTituloMuitoLongo, err)
	})

	t.Run("deve aceitar título com exatamente 255 caracteres", func(t *testing.T) {
		titulo := strings.Repeat("a", 255) // 255 caracteres
		mensagem := "Mensagem válida"

		comunicado, err := NewComunicado(idEvento, titulo, mensagem)

		assert.NoError(t, err)
		assert.Equal(t, titulo, comunicado.Titulo())
	})

	t.Run("deve retornar erro se a mensagem for vazia", func(t *testing.T) {
		titulo := "Título válido"

		_, err := NewComunicado(idEvento, titulo, "")

		assert.Error(t, err)
		assert.Equal(t, ErrMensagemObrigatoria, err)
	})

	t.Run("deve retornar erro se a mensagem for apenas espaços", func(t *testing.T) {
		titulo := "Título válido"

		_, err := NewComunicado(idEvento, titulo, "   ")

		assert.Error(t, err)
		assert.Equal(t, ErrMensagemObrigatoria, err)
	})

	t.Run("deve trimmar espaços do título e mensagem", func(t *testing.T) {
		titulo := "  Título com espaços  "
		mensagem := "  Mensagem com espaços  "

		comunicado, err := NewComunicado(idEvento, titulo, mensagem)

		assert.NoError(t, err)
		assert.Equal(t, "Título com espaços", comunicado.Titulo())
		assert.Equal(t, "Mensagem com espaços", comunicado.Mensagem())
	})
}

func TestComunicado_Editar(t *testing.T) {
	idEvento := uuid.New()
	comunicado, _ := NewComunicado(idEvento, "Título original", "Mensagem original")

	t.Run("deve editar comunicado com sucesso", func(t *testing.T) {
		novoTitulo := "Título atualizado"
		novaMensagem := "Mensagem atualizada"

		err := comunicado.Editar(novoTitulo, novaMensagem)

		assert.NoError(t, err)
		assert.Equal(t, novoTitulo, comunicado.Titulo())
		assert.Equal(t, novaMensagem, comunicado.Mensagem())
	})

	t.Run("deve retornar erro se novo título for vazio", func(t *testing.T) {
		err := comunicado.Editar("", "Nova mensagem")

		assert.Error(t, err)
		assert.Equal(t, ErrTituloObrigatorio, err)
	})

	t.Run("deve retornar erro se novo título for muito longo", func(t *testing.T) {
		tituloLongo := strings.Repeat("a", 256)

		err := comunicado.Editar(tituloLongo, "Nova mensagem")

		assert.Error(t, err)
		assert.Equal(t, ErrTituloMuitoLongo, err)
	})

	t.Run("deve retornar erro se nova mensagem for vazia", func(t *testing.T) {
		err := comunicado.Editar("Novo título", "")

		assert.Error(t, err)
		assert.Equal(t, ErrMensagemObrigatoria, err)
	})

	t.Run("deve trimmar espaços ao editar", func(t *testing.T) {
		err := comunicado.Editar("  Título trimado  ", "  Mensagem trimada  ")

		assert.NoError(t, err)
		assert.Equal(t, "Título trimado", comunicado.Titulo())
		assert.Equal(t, "Mensagem trimada", comunicado.Mensagem())
	})
}

func TestHydrateComunicado(t *testing.T) {
	t.Run("deve hidratar comunicado corretamente", func(t *testing.T) {
		id := uuid.New()
		idEvento := uuid.New()
		titulo := "Título teste"
		mensagem := "Mensagem teste"
		tempComunicado, _ := NewComunicado(uuid.New(), "temp", "temp")
		dataPublicacao := tempComunicado.DataPublicacao()

		comunicadoHidratado := HydrateComunicado(id, idEvento, titulo, mensagem, dataPublicacao)

		assert.Equal(t, id, comunicadoHidratado.ID())
		assert.Equal(t, idEvento, comunicadoHidratado.IDEvento())
		assert.Equal(t, titulo, comunicadoHidratado.Titulo())
		assert.Equal(t, mensagem, comunicadoHidratado.Mensagem())
		assert.Equal(t, dataPublicacao, comunicadoHidratado.DataPublicacao())
	})
}
