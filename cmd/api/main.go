// file: cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	guestApp "github.com/luiszkm/wedding_backend/internal/guest/application"
	guestInfra "github.com/luiszkm/wedding_backend/internal/guest/infrastructure"
	guestREST "github.com/luiszkm/wedding_backend/internal/guest/interfaces/rest"

	giftApp "github.com/luiszkm/wedding_backend/internal/gift/application"
	giftInfra "github.com/luiszkm/wedding_backend/internal/gift/infrastructure"
	giftREST "github.com/luiszkm/wedding_backend/internal/gift/interfaces/rest"
	"github.com/luiszkm/wedding_backend/internal/platform/storage"

	mbApp "github.com/luiszkm/wedding_backend/internal/messageboard/application"
	mbInfra "github.com/luiszkm/wedding_backend/internal/messageboard/infrastructure"
	mbREST "github.com/luiszkm/wedding_backend/internal/messageboard/interfaces/rest"

	galleryApp "github.com/luiszkm/wedding_backend/internal/gallery/application"
	galleryInfra "github.com/luiszkm/wedding_backend/internal/gallery/infrastructure"
	galleryREST "github.com/luiszkm/wedding_backend/internal/gallery/interfaces/rest"
)

func main() {
	port := ":3000"
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado.")
	}
	// --- Configuração ---
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_URL")
	dbURL := os.Getenv("DATABASE_URL")

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Incapaz de conectar ao banco de dados: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// --- Composição da Raiz (Wiring) ---
	storageSvc, err := storage.NewR2Storage(context.Background(), accountID, accessKeyID, secretAccessKey, bucketName, publicURL) // <-- PASSANDO A NOVA VARIÁVEL
	if err != nil {
		log.Fatalf("Falha ao inicializar o serviço de storage R2: %v", err)
	}

	// --- Repositórios ---
	guestRepo := guestInfra.NewPostgresGroupRepository(dbpool)
	presenteRepo := giftInfra.NewPostgresPresenteRepository(dbpool)
	selecaoRepo := giftInfra.NewPostgresSelecaoRepository(dbpool) // Novo repo
	recadoRepo := mbInfra.NewPostgresRecadoRepository(dbpool)
	fotoRepo := galleryInfra.NewPostgresFotoRepository(dbpool)

	// --- Serviços de Aplicação ---
	guestService := guestApp.NewGuestService(guestRepo)
	presenteService := giftApp.NewGiftService(presenteRepo, selecaoRepo) // Injetando novo repo
	recadoService := mbApp.NewMessageBoardService(recadoRepo, guestRepo)
	galleryService := galleryApp.NewGalleryService(fotoRepo, storageSvc)

	// --- Handlers ---
	guestHandler := guestREST.NewGuestHandler(guestService)
	presenteHandler := giftREST.NewGiftHandler(presenteService, storageSvc)
	recadoHandler := mbREST.NewMessageBoardHandler(recadoService)
	galleryHandler := galleryREST.NewGalleryHandler(galleryService)

	// --- Roteador e Rotas ---
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		// Rota para criar grupo (já existente)
		r.Post("/casamentos/{idCasamento}/grupos-de-convidados", guestHandler.HandleCriarGrupoDeConvidados)
		// Nova rota para acesso do convidado
		r.Get("/acesso-convidado", guestHandler.HandleObterGrupoPorChaveDeAcesso)
		// Nova rota para submissão de RSVP em lote.
		r.Post("/rsvps", guestHandler.HandleConfirmarPresenca)
		// rota para editar grupo de convidados
		r.Put("/grupos-de-convidados/{idGrupo}", guestHandler.HandleRevisarGrupo)
		// rota para criar  presentes
		r.Post("/casamentos/{idCasamento}/presentes", presenteHandler.HandleCriarPresente)
		// rota para listar presentes
		r.Get("/casamentos/{idCasamento}/presentes-publico", presenteHandler.HandleListarPresentesPublicos)
		// rota para finalizar seleção de presentes
		r.Post("/selecoes-de-presente", presenteHandler.HandleFinalizarSelecao)
		//  rota de Recados
		r.Post("/recados", recadoHandler.HandleDeixarRecado)
		r.Get("/casamentos/{idCasamento}/recados/admin", recadoHandler.HandleListarRecadosAdmin)
		r.Patch("/recados/{idRecado}", recadoHandler.HandleModerarRecado)
		r.Get("/casamentos/{idCasamento}/recados/publico", recadoHandler.HandleListarRecadosPublicos)
		// rota de Galeria
		r.Post("/casamentos/{idCasamento}/fotos", galleryHandler.HandleFazerUpload)
		r.Get("/casamentos/{idCasamento}/fotos/publico", galleryHandler.HandleListarFotosPublicas)
		r.Post("/fotos/{idFoto}/favoritar", galleryHandler.HandleAlternarFavorito)
		r.Post("/fotos/{idFoto}/rotulos", galleryHandler.HandleAdicionarRotulo)
		r.Delete("/fotos/{idFoto}/rotulos/{nomeDoRotulo}", galleryHandler.HandleRemoverRotulo)
		r.Delete("/fotos/{idFoto}", galleryHandler.HandleDeletarFoto)

	})

	log.Printf("Servidor iniciado na porta %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
