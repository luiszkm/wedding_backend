// file: internal/platform/storage/filestorage.go
package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3" // <-- ESTA LINHA FOI CORRIGIDA
	"github.com/google/uuid"
)

// FileStorage define a interface para um serviço de armazenamento de arquivos.
type FileStorage interface {
	Upload(ctx context.Context, file io.Reader, fileType string, fileSize int64) (url string, key string, err error)
	Delete(ctx context.Context, key string) error
}

// R2Storage é a implementação para o Cloudflare R2.
type R2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

// NewR2Storage cria uma nova instância do R2Storage.
func NewR2Storage(ctx context.Context, accountID, accessKeyID, secretAccessKey, bucketName, publicURL string) (FileStorage, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	customHTTPClient := &http.Client{
		Transport: customTransport,
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
		config.WithHTTPClient(customHTTPClient), // <-- Usando nosso cliente HTTP customizado
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar configuração r2: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &R2Storage{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

// Upload envia um arquivo para o R2 e retorna sua URL pública.
func (s *R2Storage) Upload(ctx context.Context, file io.Reader, fileType string, fileSize int64) (string, string, error) {
	key := fmt.Sprintf("wedding/%d-%s", time.Now().UnixNano(), uuid.New().String())

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String(fileType),
		ContentLength: &fileSize,
	})
	if err != nil {
		return "", "", fmt.Errorf("falha ao fazer upload para o r2: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.publicURL, key)
	// Retorna a URL pública, a chave única do arquivo e nenhum erro.
	return url, key, nil
}
func (s *R2Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("falha ao deletar objeto no r2: %w", err)
	}
	return nil
}

type MockStorage struct{}

func NewMockStorage() FileStorage {
	return &MockStorage{}
}

// Upload no mock agora também retorna uma chave fake.
func (s *MockStorage) Upload(ctx context.Context, file io.Reader, fileType string, fileSize int64) (string, string, error) {
	fakeURL := "https://fake-bucket.s3.amazonaws.com/presentes/placeholder.jpg"
	fakeKey := "presentes/placeholder.jpg"
	return fakeURL, fakeKey, nil
}
func (s *MockStorage) Delete(ctx context.Context, key string) error {
	// No mock, apenas logamos a ação e retornamos sucesso.
	log.Printf("[MOCK] Deletando arquivo com a chave: %s", key)
	return nil
}
