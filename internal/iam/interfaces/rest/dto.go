// file: internal/iam/interfaces/rest/dto.go
package rest

type RegistrarRequestDTO struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone,omitempty"`
	Senha    string `json:"senha"`
}

type UsuarioResponseDTO struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type LoginRequestDTO struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"accessToken"`
}
