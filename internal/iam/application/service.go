// file: internal/iam/application/service.go
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/luiszkm/wedding_backend/internal/iam/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
)

type IAMService struct {
	repo       domain.UsuarioRepository
	jwtService *auth.JWTService // Nova dependência
}

func NewIAMService(repo domain.UsuarioRepository, jwtService *auth.JWTService) *IAMService {
	return &IAMService{repo: repo, jwtService: jwtService}
}

func (s *IAMService) RegistrarNovoUsuario(ctx context.Context, nome, email, telefone, senha string) (*domain.Usuario, error) {
	// 1. Verifica se o e-mail já existe
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		// Se não deu erro, significa que encontrou um usuário.
		return nil, domain.ErrEmailJaExiste
	}
	// Se o erro for "não encontrado", ótimo, podemos continuar.
	// Qualquer outro erro é um problema técnico.
	if !errors.Is(err, domain.ErrUsuarioNaoEncontrado) {
		return nil, fmt.Errorf("erro ao verificar e-mail existente: %w", err)
	}

	// 2. Cria o novo agregado de usuário (que já faz o hash da senha)
	novoUsuario, err := domain.NewUsuario(nome, email, telefone, senha)
	if err != nil {
		return nil, err
	}

	// 3. Salva o novo usuário no banco de dados
	if err := s.repo.Save(ctx, novoUsuario); err != nil {
		return nil, fmt.Errorf("falha ao salvar novo usuário no repositório: %w", err)
	}

	return novoUsuario, nil
}

func (s *IAMService) Login(ctx context.Context, email, senha string) (string, error) {
	// 1. Busca o usuário pelo email.
	usuario, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUsuarioNaoEncontrado) {
			// Não diferenciar "usuário não existe" de "senha incorreta" por segurança.
			return "", domain.ErrCredenciaisInvalidas // Crie este erro em domain/usuario.go
		}
		return "", err // Outro erro técnico
	}

	// 2. Verifica se a senha fornecida corresponde ao hash armazenado.
	if !usuario.VerificarSenha(senha) {
		return "", domain.ErrCredenciaisInvalidas
	}

	// 3. Gera o token JWT.
	token, err := s.jwtService.GenerateToken(usuario.ID())
	if err != nil {
		return "", fmt.Errorf("falha ao gerar token de acesso: %w", err)
	}

	return token, nil
}

// LoginWithUserInfo retorna tanto o token quanto informações do usuário
func (s *IAMService) LoginWithUserInfo(ctx context.Context, email, senha string) (string, *domain.Usuario, error) {
	// 1. Busca o usuário pelo email.
	usuario, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUsuarioNaoEncontrado) {
			// Não diferenciar "usuário não existe" de "senha incorreta" por segurança.
			return "", nil, domain.ErrCredenciaisInvalidas
		}
		return "", nil, err // Outro erro técnico
	}

	// 2. Verifica se a senha fornecida corresponde ao hash armazenado.
	if !usuario.VerificarSenha(senha) {
		return "", nil, domain.ErrCredenciaisInvalidas
	}

	// 3. Gera o token JWT.
	token, err := s.jwtService.GenerateToken(usuario.ID())
	if err != nil {
		return "", nil, fmt.Errorf("falha ao gerar token de acesso: %w", err)
	}

	return token, usuario, nil
}
