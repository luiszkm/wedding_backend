package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTituloObrigatorio       = errors.New("o título é obrigatório")
	ErrTituloMuitoLongo        = errors.New("o título deve ter no máximo 255 caracteres")
	ErrMensagemObrigatoria     = errors.New("a mensagem é obrigatória")
	ErrComunicadoNaoEncontrado = errors.New("comunicado não encontrado")
)

type Comunicado struct {
	id             uuid.UUID
	idEvento       uuid.UUID
	titulo         string
	mensagem       string
	dataPublicacao time.Time
}

func NewComunicado(idEvento uuid.UUID, titulo, mensagem string) (*Comunicado, error) {
	titulo = strings.TrimSpace(titulo)
	mensagem = strings.TrimSpace(mensagem)

	if titulo == "" {
		return nil, ErrTituloObrigatorio
	}
	if len(titulo) > 255 {
		return nil, ErrTituloMuitoLongo
	}
	if mensagem == "" {
		return nil, ErrMensagemObrigatoria
	}

	return &Comunicado{
		id:             uuid.New(),
		idEvento:       idEvento,
		titulo:         titulo,
		mensagem:       mensagem,
		dataPublicacao: time.Now(),
	}, nil
}

func HydrateComunicado(id, idEvento uuid.UUID, titulo, mensagem string, dataPublicacao time.Time) *Comunicado {
	return &Comunicado{
		id:             id,
		idEvento:       idEvento,
		titulo:         titulo,
		mensagem:       mensagem,
		dataPublicacao: dataPublicacao,
	}
}

func (c *Comunicado) Editar(novoTitulo, novaMensagem string) error {
	novoTitulo = strings.TrimSpace(novoTitulo)
	novaMensagem = strings.TrimSpace(novaMensagem)

	if novoTitulo == "" {
		return ErrTituloObrigatorio
	}
	if len(novoTitulo) > 255 {
		return ErrTituloMuitoLongo
	}
	if novaMensagem == "" {
		return ErrMensagemObrigatoria
	}

	c.titulo = novoTitulo
	c.mensagem = novaMensagem
	return nil
}

func (c *Comunicado) ID() uuid.UUID             { return c.id }
func (c *Comunicado) IDEvento() uuid.UUID       { return c.idEvento }
func (c *Comunicado) Titulo() string            { return c.titulo }
func (c *Comunicado) Mensagem() string          { return c.mensagem }
func (c *Comunicado) DataPublicacao() time.Time { return c.dataPublicacao }
