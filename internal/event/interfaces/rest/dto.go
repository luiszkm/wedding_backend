// file: internal/event/interfaces/rest/dto.go
package rest

import "time"

type CriarEventoRequestDTO struct {
	Nome    string    `json:"nomeDoEvento"`
	Data    time.Time `json:"dataDoEvento"`
	Tipo    string    `json:"tipo"`
	UrlSlug string    `json:"urlSlug"`
}

type CriarEventoResponseDTO struct {
	IDEvento string `json:"idEvento"`
}

type EventoResponseDTO struct {
	ID      string    `json:"id"`
	Nome    string    `json:"nome"`
	Data    time.Time `json:"data"`
	Tipo    string    `json:"tipo"`
	UrlSlug string    `json:"urlSlug"`
}
