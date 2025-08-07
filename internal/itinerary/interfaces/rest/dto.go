package rest

import "time"

// CreateItineraryItemRequestDTO é o corpo da requisição para criar um item do roteiro
type CreateItineraryItemRequestDTO struct {
	Horario           time.Time `json:"horario"`
	TituloAtividade   string    `json:"tituloAtividade"`
	DescricaoAtividade *string   `json:"descricaoAtividade,omitempty"`
}

// CreateItineraryItemResponseDTO é a resposta da criação de um item do roteiro
type CreateItineraryItemResponseDTO struct {
	ID string `json:"id"`
}

// UpdateItineraryItemRequestDTO é o corpo da requisição para atualizar um item do roteiro
type UpdateItineraryItemRequestDTO struct {
	Horario           time.Time `json:"horario"`
	TituloAtividade   string    `json:"tituloAtividade"`
	DescricaoAtividade *string   `json:"descricaoAtividade,omitempty"`
}

// ItineraryItemDTO representa um item do roteiro na resposta da API
type ItineraryItemDTO struct {
	ID                string    `json:"id"`
	IDEvento          string    `json:"idEvento"`
	Horario           time.Time `json:"horario"`
	TituloAtividade   string    `json:"tituloAtividade"`
	DescricaoAtividade *string   `json:"descricaoAtividade,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// GetItineraryResponseDTO é a resposta para a listagem de itens do roteiro
type GetItineraryResponseDTO struct {
	Items []ItineraryItemDTO `json:"items"`
}