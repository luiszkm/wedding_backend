// file: internal/gift/interfaces/rest/dto.go
package rest

// DetalhesPresenteDTO representa os detalhes polimórficos na requisição.
type DetalhesPresenteDTO struct {
	Tipo       string `json:"tipo"`
	LinkDaLoja string `json:"linkDaLoja,omitempty"`
}

// CriarPresenteRequestDTO é o contrato de entrada para criar um presente.
type CriarPresenteRequestDTO struct {
	Nome       string              `json:"nome"`
	Descricao  string              `json:"descricao"`
	EhFavorito bool                `json:"ehFavorito"`
	FotoURL    string              `json:"fotoUrl,omitempty"` // Usado se não houver upload
	Categoria  string              `json:"categoria"`         // Novo campo
	Detalhes   DetalhesPresenteDTO `json:"detalhes"`
}

// CriarPresenteResponseDTO é a resposta de sucesso.
type CriarPresenteResponseDTO struct {
	IDPresente string `json:"idPresente"`
}

type PresentePublicoDTO struct {
	ID         string              `json:"id"`
	Nome       string              `json:"nome"`
	Descricao  string              `json:"descricao"`
	FotoURL    string              `json:"fotoUrl"`
	EhFavorito bool                `json:"ehFavorito"`
	Categoria  string              `json:"categoria"`
	Detalhes   DetalhesPresenteDTO `json:"detalhes"`
}

type FinalizarSelecaoRequestDTO struct {
	ChaveDeAcesso   string   `json:"chaveDeAcesso"`
	IDsDosPresentes []string `json:"idsDosPresentes"`
}

type SelecaoConfirmadaDTO struct {
	IDSelecao            string                  `json:"idSelecao"`
	Mensagem             string                  `json:"mensagem"`
	PresentesConfirmados []PresenteConfirmadoDTO `json:"presentesConfirmados"`
}

type PresenteConfirmadoDTO struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type ConflitoSelecaoDTO struct {
	Codigo                string   `json:"codigo"`
	Mensagem              string   `json:"mensagem"`
	PresentesConflitantes []string `json:"presentesConflitantes"`
}
