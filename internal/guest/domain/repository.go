// file: internal/guest/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

// RSVPStats representa as estatísticas de RSVP para um evento
type RSVPStats struct {
	TotalGrupos           int
	TotalConvidados       int
	ConvidadosConfirmados int
	ConvidadosRecusados   int
	ConvidadosPendentes   int
	PercentualConfirmado  float64
	PercentualRecusado    float64
	PercentualPendente    float64
}

type GroupRepository interface {
	Save(ctx context.Context, group *GrupoDeConvidados) error
	FindByAccessKey(ctx context.Context, eventID uuid.UUID, accessKey string) (*GrupoDeConvidados, error)
	Update(ctx context.Context, userID uuid.UUID, group *GrupoDeConvidados) error        // <-- userID adicionado
	FindByID(ctx context.Context, userID, groupID uuid.UUID) (*GrupoDeConvidados, error) // <-- userID adicionado
	UpdateRSVP(ctx context.Context, group *GrupoDeConvidados) error
	FindAllByEventID(ctx context.Context, userID, eventID uuid.UUID, statusFilter string) ([]*GrupoDeConvidados, error)
	Delete(ctx context.Context, userID, groupID uuid.UUID) error
	GetRSVPStats(ctx context.Context, userID, eventID uuid.UUID) (*RSVPStats, error)
}
