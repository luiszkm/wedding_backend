// file: internal/guest/interfaces/rest/dto.go
package rest

// CriarGrupoRequestDTO é o contrato de entrada da API.
type CriarGrupoRequestDTO struct {
	ChaveDeAcesso      string   `json:"chaveDeAcesso"`
	NomesDosConvidados []string `json:"nomesDosConvidados"`
}

// CriarGrupoResponseDTO é o contrato de saída da API.
type CriarGrupoResponseDTO struct {
	IDGrupo string `json:"idGrupo"`
}

// GrupoParaConfirmacaoDTO representa os dados do grupo para o convidado confirmar presença.
type GrupoParaConfirmacaoDTO struct {
	IDGrupo    string         `json:"idGrupo"`
	Convidados []ConvidadoDTO `json:"convidados"`
}

// ConvidadoDTO representa um único convidado dentro do grupo.
type ConvidadoDTO struct {
	ID         string `json:"id"`
	Nome       string `json:"nome"`
	StatusRSVP string `json:"statusRSVP"` // e.g., "PENDENTE"
}
