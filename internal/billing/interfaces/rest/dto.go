// file: internal/billing/interfaces/rest/dto.go
package rest

// PlanoDTO representa a visão pública de um plano de assinatura.
type PlanoDTO struct {
	ID                  string `json:"id"`
	Nome                string `json:"nome"`
	PrecoEmCentavos     int    `json:"precoEmCentavos"`
	NumeroMaximoEventos int    `json:"numeroMaximoEventos"`
	DuracaoEmDias       int    `json:"duracaoEmDias"`
}

type CriarAssinaturaRequestDTO struct {
	IDPlano string `json:"idPlano"`
}

type CriarAssinaturaResponseDTO struct {
	IDAssinatura string `json:"idAssinatura"`
	Status       string `json:"status"`
	CheckoutURL  string `json:"checkoutUrl"`
}
