package rest

import "time"

type CriarComunicadoRequestDTO struct {
	Titulo   string `json:"titulo"`
	Mensagem string `json:"mensagem"`
}

type EditarComunicadoRequestDTO struct {
	Titulo   string `json:"titulo"`
	Mensagem string `json:"mensagem"`
}

type ComunicadoResponseDTO struct {
	ID             string    `json:"id"`
	IDEvento       string    `json:"idEvento"`
	Titulo         string    `json:"titulo"`
	Mensagem       string    `json:"mensagem"`
	DataPublicacao time.Time `json:"dataPublicacao"`
}
