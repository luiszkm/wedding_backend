// file: internal/messageboard/interfaces/rest/dto.go
package rest

import "time"

// DeixarRecadoRequestDTO é o corpo da requisição para postar um recado.
type DeixarRecadoRequestDTO struct {
	ChaveDeAcesso string `json:"chaveDeAcesso"`
	NomeDoAutor   string `json:"nomeDoAutor"`
	Texto         string `json:"texto"`
}

type RecadoAdminDTO struct {
	ID            string    `json:"id"`
	NomeDoAutor   string    `json:"nomeDoAutor"`
	Texto         string    `json:"texto"`
	Status        string    `json:"status"`
	EhFavorito    bool      `json:"ehFavorito"`
	DataDeCriacao time.Time `json:"dataDeCriacao"`
}
type ModerarRecadoRequestDTO struct {
	Status     *string `json:"status,omitempty"`
	EhFavorito *bool   `json:"ehFavorito,omitempty"`
}
type RecadoPublicoDTO struct {
	ID            string    `json:"id"`
	NomeDoAutor   string    `json:"nomeDoAutor"`
	Texto         string    `json:"texto"`
	EhFavorito    bool      `json:"ehFavorito"`
	DataDeCriacao time.Time `json:"dataDeCriacao"`
}
