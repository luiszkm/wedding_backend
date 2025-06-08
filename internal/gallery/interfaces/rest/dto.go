// file: internal/gallery/interfaces/rest/dto.go
package rest

// UploadFotosResponseDTO é a resposta de sucesso para o upload.
type UploadFotosResponseDTO struct {
	IDsDasFotosCriadas []string `json:"idsDasFotosCriadas"`
}
type FotoPublicaDTO struct {
	ID         string   `json:"id"`
	URLPublica string   `json:"urlPublica"`
	EhFavorito bool     `json:"ehFavorito"`
	Rotulos    []string `json:"rotulos"`
}
type AdicionarRotuloRequestDTO struct {
	NomeDoRotulo string `json:"nomeDoRotulo"`
}
