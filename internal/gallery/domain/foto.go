// file: internal/gallery/domain/foto.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Rotulo é um tipo para representar as categorias de fotos de forma segura.
type Rotulo string

// Constantes para os rótulos válidos, conforme a documentação.
const (
	RotuloMain      Rotulo = "MAIN"
	RotuloCasamento Rotulo = "CASAMENTO"
	RotuloLuaDeMel  Rotulo = "LUADEMEL"
	RotuloHistoria  Rotulo = "HISTORIA"
	RotuloFamilia   Rotulo = "FAMILIA"
	RotuloOutros    Rotulo = "OUTROS"
)

var (
	ErrRotuloInvalido    = errors.New("rótulo inválido")
	ErrFotoNaoEncontrada = errors.New("foto não encontrada")
)

// IsValid verifica se uma string corresponde a um Rótulo válido.
func (r Rotulo) IsValid() bool {
	switch r {
	case RotuloMain, RotuloCasamento, RotuloLuaDeMel, RotuloHistoria, RotuloFamilia, RotuloOutros:
		return true
	}
	return false
}

// Foto é o agregado raiz do contexto de Galeria.
type Foto struct {
	id          uuid.UUID
	idCasamento uuid.UUID
	storageKey  string
	urlPublica  string
	ehFavorito  bool
	createdAt   time.Time
	rotulos     map[Rotulo]bool // Usamos um mapa para evitar rótulos duplicados
}

// NewFoto é a fábrica para criar uma Foto em estado válido.
func NewFoto(idCasamento uuid.UUID, storageKey, urlPublica string) *Foto {
	return &Foto{
		id:          uuid.New(),
		idCasamento: idCasamento,
		storageKey:  storageKey,
		urlPublica:  urlPublica,
		ehFavorito:  false,
		createdAt:   time.Now(),
		rotulos:     make(map[Rotulo]bool),
	}
}

// AdicionarRotulo adiciona um rótulo à foto.
func (f *Foto) AdicionarRotulo(rotulo Rotulo) error {
	if !rotulo.IsValid() {
		return ErrRotuloInvalido
	}
	f.rotulos[rotulo] = true
	return nil
}
func HydrateFoto(id, idCasamento uuid.UUID, storageKey, urlPublica string, ehFavorito bool, createdAt time.Time, rotulos []Rotulo) *Foto {
	rotulosMap := make(map[Rotulo]bool)
	for _, r := range rotulos {
		rotulosMap[r] = true
	}

	return &Foto{
		id:          id,
		idCasamento: idCasamento,
		storageKey:  storageKey,
		urlPublica:  urlPublica,
		ehFavorito:  ehFavorito,
		createdAt:   createdAt,
		rotulos:     rotulosMap,
	}
}
func (f *Foto) AlternarFavorito() {
	f.ehFavorito = !f.ehFavorito
}
func (f *Foto) RemoverRotulo(rotulo Rotulo) error {
	if !rotulo.IsValid() {
		return ErrRotuloInvalido
	}
	// Não há problema em tentar remover um rótulo que não existe.
	delete(f.rotulos, rotulo)
	return nil
}

// Getters
func (f *Foto) ID() uuid.UUID          { return f.id }
func (f *Foto) IDCasamento() uuid.UUID { return f.idCasamento }
func (f *Foto) StorageKey() string     { return f.storageKey }
func (f *Foto) URLPublica() string     { return f.urlPublica }
func (f *Foto) EhFavorito() bool       { return f.ehFavorito }
func (f *Foto) CreatedAt() time.Time   { return f.createdAt }
func (f *Foto) Rotulos() []Rotulo {
	rotulos := make([]Rotulo, 0, len(f.rotulos))
	for r := range f.rotulos {
		rotulos = append(rotulos, r)
	}
	return rotulos
}
