// file: internal/iam/domain/usuario.go
package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailJaExiste        = errors.New("o e-mail informado já está em uso")
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas, verifique seu e-mail e senha")
)

type Usuario struct {
	id           uuid.UUID
	nome         string
	email        string
	telefone     string
	passwordHash string
}

// NewUsuario é a fábrica para criar um novo usuário.
func NewUsuario(nome, email, telefone, senhaPura string) (*Usuario, error) {
	if nome == "" || email == "" || senhaPura == "" {
		return nil, errors.New("nome, email e senha são obrigatórios")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senhaPura), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar hash da senha: %w", err)
	}

	return &Usuario{
		id:           uuid.New(),
		nome:         nome,
		email:        email,
		telefone:     telefone,
		passwordHash: string(hash),
	}, nil
}

// HydrateUsuario reconstrói um objeto a partir dos dados do banco.
func HydrateUsuario(id uuid.UUID, nome, email, telefone, hash string) *Usuario {
	return &Usuario{
		id:           id,
		nome:         nome,
		email:        email,
		telefone:     telefone,
		passwordHash: hash,
	}
}

// Getters para acesso seguro aos campos.
func (u *Usuario) ID() uuid.UUID        { return u.id }
func (u *Usuario) Nome() string         { return u.nome }
func (u *Usuario) Email() string        { return u.email }
func (u *Usuario) Telefone() string     { return u.telefone }
func (u *Usuario) PasswordHash() string { return u.passwordHash }

// VerificarSenha compara uma senha em texto puro com o hash armazenado.
func (u *Usuario) VerificarSenha(senhaPura string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(senhaPura))
	return err == nil
}
