// file: internal/gallery/application/service.go
package application

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/gallery/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/storage"
)

type GalleryService struct {
	fotoRepo domain.FotoRepository
	storage  storage.FileStorage
}

func NewGalleryService(fotoRepo domain.FotoRepository, storage storage.FileStorage) *GalleryService {
	return &GalleryService{fotoRepo: fotoRepo, storage: storage}
}

func (s *GalleryService) FazerUploadDeFotos(ctx context.Context, idCasamento uuid.UUID, rotulo string, files []*multipart.FileHeader) ([]uuid.UUID, error) {
	rotuloParaAdicionar := domain.Rotulo(strings.ToUpper(rotulo))
	if !rotuloParaAdicionar.IsValid() {
		return nil, domain.ErrRotuloInvalido
	}

	var novasFotos []*domain.Foto

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("falha ao abrir arquivo %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		// A chamada agora recebe três valores de retorno: url, key e err.
		url, key, err := s.storage.Upload(ctx, file, fileHeader.Header.Get("Content-Type"), fileHeader.Size)
		if err != nil {
			return nil, fmt.Errorf("falha no upload do arquivo %s: %w", fileHeader.Filename, err)
		}

		// A criação da entidade fica muito mais limpa e explícita.
		foto := domain.NewFoto(idCasamento, key, url)

		foto.AdicionarRotulo(domain.RotuloCasamento)
		foto.AdicionarRotulo(rotuloParaAdicionar)

		novasFotos = append(novasFotos, foto)
	}

	if err := s.fotoRepo.SalvarMultiplas(ctx, novasFotos); err != nil {
		return nil, fmt.Errorf("falha ao salvar metadados das fotos: %w", err)
	}

	ids := make([]uuid.UUID, len(novasFotos))
	for i, f := range novasFotos {
		ids[i] = f.ID()
	}

	return ids, nil
}
func (s *GalleryService) ListarFotosPublicas(ctx context.Context, casamentoID uuid.UUID, filtroRotulo string) ([]*domain.Foto, error) {
	rotulo := domain.Rotulo(strings.ToUpper(filtroRotulo))
	// Valida o rótulo apenas se um filtro foi de fato fornecido.
	if filtroRotulo != "" && !rotulo.IsValid() {
		return nil, domain.ErrRotuloInvalido
	}

	fotos, err := s.fotoRepo.ListarPublicasPorCasamento(ctx, casamentoID, rotulo)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista pública de fotos: %w", err)
	}
	return fotos, nil
}
func (s *GalleryService) AlternarFavoritoFoto(ctx context.Context, fotoID uuid.UUID) (*domain.Foto, error) {
	// 1. Carregar o agregado
	foto, err := s.fotoRepo.FindByID(ctx, fotoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar foto para favoritar: %w", err)
	}

	// 2. Executar a lógica de negócio no domínio
	foto.AlternarFavorito()

	// 3. Persistir as alterações
	if err := s.fotoRepo.Update(ctx, foto); err != nil {
		return nil, fmt.Errorf("falha ao salvar alteração de favorito: %w", err)
	}

	return foto, nil
}

func (s *GalleryService) AdicionarRotulo(ctx context.Context, fotoID uuid.UUID, nomeRotulo string) error {
	foto, err := s.fotoRepo.FindByID(ctx, fotoID)
	if err != nil {
		return fmt.Errorf("falha ao buscar foto para adicionar rótulo: %w", err)
	}

	rotulo := domain.Rotulo(strings.ToUpper(nomeRotulo))
	if err := foto.AdicionarRotulo(rotulo); err != nil {
		return err
	}

	return s.fotoRepo.Update(ctx, foto)
}

func (s *GalleryService) RemoverRotulo(ctx context.Context, fotoID uuid.UUID, nomeRotulo string) error {
	foto, err := s.fotoRepo.FindByID(ctx, fotoID)
	if err != nil {
		return fmt.Errorf("falha ao buscar foto para remover rótulo: %w", err)
	}

	rotulo := domain.Rotulo(strings.ToUpper(nomeRotulo))
	if err := foto.RemoverRotulo(rotulo); err != nil {
		return err
	}

	return s.fotoRepo.Update(ctx, foto)
}
func (s *GalleryService) DeletarFoto(ctx context.Context, fotoID uuid.UUID) error {
	// 1. Busca os metadados da foto para obter a chave de armazenamento (storageKey).
	foto, err := s.fotoRepo.FindByID(ctx, fotoID)
	if err != nil {
		return fmt.Errorf("falha ao buscar foto para deletar: %w", err)
	}

	// 2. Apaga o arquivo no bucket R2.
	if err := s.storage.Delete(ctx, foto.StorageKey()); err != nil {
		// Logamos o erro mas podemos decidir continuar para apagar o registro do banco mesmo assim.
		// Ou podemos parar aqui. Parar é mais seguro para evitar referências quebradas.
		log.Printf("AVISO: falha ao deletar arquivo do storage (%s), abortando deleção do registro no DB: %v", foto.StorageKey(), err)
		return fmt.Errorf("falha ao deletar arquivo no storage: %w", err)
	}

	// 3. Se o arquivo foi apagado do storage com sucesso, apaga o registro do banco.
	if err := s.fotoRepo.Delete(ctx, fotoID); err != nil {
		return fmt.Errorf("falha ao deletar registro da foto no banco: %w", err)
	}

	return nil
}
