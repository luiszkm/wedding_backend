// file: internal/iam/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/iam/domain"
)

type PostgresUsuarioRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUsuarioRepository(db *pgxpool.Pool) domain.UsuarioRepository {
	return &PostgresUsuarioRepository{db: db}
}

func (r *PostgresUsuarioRepository) Save(ctx context.Context, usuario *domain.Usuario) error {
	sql := `INSERT INTO usuarios (id, nome, email, telefone, password_hash) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, sql,
		usuario.ID(),
		usuario.Nome(),
		usuario.Email(),
		usuario.Telefone(),
		usuario.PasswordHash(),
	)
	if err != nil {
		return fmt.Errorf("falha ao salvar usuário: %w", err)
	}
	return nil
}

func (r *PostgresUsuarioRepository) FindByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	sql := `SELECT id, nome, email, telefone, password_hash FROM usuarios WHERE email = $1`
	row := r.db.QueryRow(ctx, sql, email)

	var id uuid.UUID
	var nome, uEmail, telefone, hash string

	err := row.Scan(&id, &nome, &uEmail, &telefone, &hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Não é um erro técnico, simplesmente não encontrou.
			return nil, domain.ErrUsuarioNaoEncontrado
		}
		return nil, fmt.Errorf("falha ao buscar usuário por email: %w", err)
	}

	return domain.HydrateUsuario(id, nome, uEmail, telefone, hash), nil
}
