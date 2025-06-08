// file: internal/platform/storage/filestorage.go
package storage

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3" // <-- ESTA LINHA FOI CORRIGIDA
	"github.com/google/uuid"
)

// FileStorage define a interface para um serviço de armazenamento de arquivos.
type FileStorage interface {
	Upload(ctx context.Context, file io.Reader, fileType string, fileSize int64) (string, error)
}

// R2Storage é a implementação para o Cloudflare R2.
type R2Storage struct {
	client     *s3.Client
	bucketName string
	accountID  string
	publicURL  string
}

// NewR2Storage cria uma nova instância do R2Storage.
func NewR2Storage(ctx context.Context, accountID, accessKeyID, secretAccessKey, bucketName, publicURL string) (FileStorage, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	// --- Bloco Novo: Criando um cliente HTTP customizado ---
	// Forçamos o uso de TLS 1.2, que é mais robusto e pode resolver problemas de handshake.
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	customHTTPClient := &http.Client{
		Transport: customTransport,
	}
	// --- Fim do Bloco Novo ---

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
		accountID:  accountID,
		publicURL:  publicURL,
	}, nil
}

// Upload envia um arquivo para o R2 e retorna sua URL pública.
func (s *R2Storage) Upload(ctx context.Context, file io.Reader, fileType string, fileSize int64) (string, error) {
	// Gera uma chave única para o arquivo para evitar colisões de nome.
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("erro de debug: falha ao ler o stream para a memória: %w", err)
	}
	// 2. Agora você pode "ver" os dados.
	log.Printf("[DEBUG] Tamanho do arquivo lido na memória: %d bytes", len(fileBytes))
	log.Printf("[DEBUG] Tipo de conteúdo detectado: %s", http.DetectContentType(fileBytes))

	// Cuidado ao imprimir o conteúdo se o arquivo for grande.
	// Imprimir os primeiros 100 bytes pode ser útil.
	if len(fileBytes) > 100 {
		log.Printf("[DEBUG] Primeiros 100 bytes (em hexadecimal): %x", fileBytes[:100])
	} else {
		log.Printf("[DEBUG] Conteúdo (em hexadecimal): %x", fileBytes)
	}

	// 3. Cria um NOVO leitor a partir dos bytes em memória para o upload.
	bodyReader := bytes.NewReader(fileBytes)

	key := fmt.Sprintf("presentes/%s", uuid.New().String())

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          bodyReader, // Usa o novo leitor
		ContentType:   aws.String(fileType),
		ContentLength: aws.Int64(int64(len(fileBytes))), // Usa o tamanho real lido
	})
	if err != nil {
		return "", fmt.Errorf("falha ao fazer upload para o r2: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.publicURL, key)
	return url, nil
}
