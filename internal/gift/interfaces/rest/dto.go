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
	FotoURL    string              `json:"fotoUrl,omitempty"`
	Categoria  string              `json:"categoria"`
	Detalhes   DetalhesPresenteDTO `json:"detalhes"`
	Tipo       string              `json:"tipo"` // INTEGRAL ou FRACIONADO
	// Campos para presentes fracionados
	ValorTotal  *float64 `json:"valorTotal,omitempty"`
	NumeroCotas *int     `json:"numeroCotas,omitempty"`
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
	Tipo       string              `json:"tipo"`
	Status     string              `json:"status"`
	// Campos para presentes fracionados
	ValorTotal        *float64 `json:"valorTotal,omitempty"`
	ValorCota         *float64 `json:"valorCota,omitempty"`
	CotasTotais       *int     `json:"cotasTotais,omitempty"`
	CotasDisponiveis  *int     `json:"cotasDisponiveis,omitempty"`
	CotasSelecionadas *int     `json:"cotasSelecionadas,omitempty"`
}

type ItemSelecaoDTO struct {
	IDPresente string `json:"idPresente"`
	Quantidade int    `json:"quantidade"`
}

type FinalizarSelecaoRequestDTO struct {
	ChaveDeAcesso string           `json:"chaveDeAcesso"`
	Itens         []ItemSelecaoDTO `json:"itens"`
}

// DTO legacy mantido para compatibilidade
type FinalizarSelecaoLegacyRequestDTO struct {
	ChaveDeAcesso   string   `json:"chaveDeAcesso"`
	IDsDosPresentes []string `json:"idsDosPresentes"`
}

type SelecaoConfirmadaDTO struct {
	IDSelecao            string                  `json:"idSelecao"`
	Mensagem             string                  `json:"mensagem"`
	ValorTotal           float64                 `json:"valorTotal"`
	PresentesConfirmados []PresenteConfirmadoDTO `json:"presentesConfirmados"`
}

type PresenteConfirmadoDTO struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Quantidade int      `json:"quantidade"`
	ValorCota  *float64 `json:"valorCota,omitempty"`
	ValorTotal *float64 `json:"valorTotal,omitempty"`
}

type ConflitoSelecaoDTO struct {
	Codigo                string   `json:"codigo"`
	Mensagem              string   `json:"mensagem"`
	PresentesConflitantes []string `json:"presentesConflitantes"`
}
