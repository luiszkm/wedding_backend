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
	"github.com/stripe/stripe-go/v79"

	guestApp "github.com/luiszkm/wedding_backend/internal/guest/application"
	guestInfra "github.com/luiszkm/wedding_backend/internal/guest/infrastructure"
	guestREST "github.com/luiszkm/wedding_backend/internal/guest/interfaces/rest"

	giftApp "github.com/luiszkm/wedding_backend/internal/gift/application"
	giftInfra "github.com/luiszkm/wedding_backend/internal/gift/infrastructure"
	giftREST "github.com/luiszkm/wedding_backend/internal/gift/interfaces/rest"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/storage"

	mbApp "github.com/luiszkm/wedding_backend/internal/messageboard/application"
	mbInfra "github.com/luiszkm/wedding_backend/internal/messageboard/infrastructure"
	mbREST "github.com/luiszkm/wedding_backend/internal/messageboard/interfaces/rest"

	galleryApp "github.com/luiszkm/wedding_backend/internal/gallery/application"
	galleryInfra "github.com/luiszkm/wedding_backend/internal/gallery/infrastructure"
	galleryREST "github.com/luiszkm/wedding_backend/internal/gallery/interfaces/rest"

	iamApp "github.com/luiszkm/wedding_backend/internal/iam/application"
	iamInfra "github.com/luiszkm/wedding_backend/internal/iam/infrastructure"
	iamREST "github.com/luiszkm/wedding_backend/internal/iam/interfaces/rest"

	eventApp "github.com/luiszkm/wedding_backend/internal/event/application"
	eventInfra "github.com/luiszkm/wedding_backend/internal/event/infrastructure"
	eventREST "github.com/luiszkm/wedding_backend/internal/event/interfaces/rest"

	billingApp "github.com/luiszkm/wedding_backend/internal/billing/application"
	billingInfra "github.com/luiszkm/wedding_backend/internal/billing/infrastructure"
	billingREST "github.com/luiszkm/wedding_backend/internal/billing/interfaces/rest"

	// pageTemplateApp "github.com/luiszkm/wedding_backend/internal/pagetemplate/application"
	// pageTemplateREST "github.com/luiszkm/wedding_backend/internal/pagetemplate/interfaces/rest"
	// "github.com/luiszkm/wedding_backend/internal/platform/template"
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
	jwtSecret := os.Getenv("JWT_SECRET")
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	// Lemos o segredo do webhook aqui, uma única vez.
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
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
	// --- Inicialização dos Serviços ---
	jwtService := auth.NewJWTService(jwtSecret)
	paymentGateway := billingInfra.NewStripeGateway(stripe.Key)

	// ...

	// --- Repositórios ---
	guestRepo := guestInfra.NewPostgresGroupRepository(dbpool)
	presenteRepo := giftInfra.NewPostgresPresenteRepository(dbpool)
	selecaoRepo := giftInfra.NewPostgresSelecaoRepository(dbpool) // Novo repo
	recadoRepo := mbInfra.NewPostgresRecadoRepository(dbpool)
	fotoRepo := galleryInfra.NewPostgresFotoRepository(dbpool)
	usuarioRepo := iamInfra.NewPostgresUsuarioRepository(dbpool)
	eventRepo := eventInfra.NewPostgresEventoRepository(dbpool)
	planoRepo := billingInfra.NewPostgresPlanoRepository(dbpool)
	billingRepo := billingInfra.NewPostgresAssinaturaRepository(dbpool)

	// --- Serviços de Aplicação ---
	guestService := guestApp.NewGuestService(guestRepo)
	presenteService := giftApp.NewGiftService(presenteRepo, selecaoRepo, eventRepo)
	recadoService := mbApp.NewMessageBoardService(recadoRepo, guestRepo, eventRepo)
	galleryService := galleryApp.NewGalleryService(fotoRepo, storageSvc)
	iamService := iamApp.NewIAMService(usuarioRepo, jwtService)
	eventService := eventApp.NewEventService(eventRepo)
	billingService := billingApp.NewBillingService(planoRepo, billingRepo, paymentGateway)

	// --- Handlers ---
	guestHandler := guestREST.NewGuestHandler(guestService)
	presenteHandler := giftREST.NewGiftHandler(presenteService, storageSvc)
	recadoHandler := mbREST.NewMessageBoardHandler(recadoService)
	galleryHandler := galleryREST.NewGalleryHandler(galleryService)
	iamHandler := iamREST.NewIAMHandler(iamService)
	eventHandler := eventREST.NewEventHandler(eventService)
	billingHandler := billingREST.NewBillingHandler(billingService, stripeWebhookSecret)

	// --- Roteador e Rotas ---
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	authMiddleware := auth.Authenticator(jwtService)

	r.Route("/v1", func(r chi.Router) {
		// --- Rotas Públicas ---
		r.Post("/usuarios/registrar", iamHandler.HandleRegistrar)
		r.Post("/usuarios/login", iamHandler.HandleLogin)
		r.Get("/casamentos/{idCasamento}/recados/publico", recadoHandler.HandleListarRecadosPublicos)
		r.Get("/casamentos/{idCasamento}/presentes-publico", presenteHandler.HandleListarPresentesPublicos)
		r.Post("/rsvps", guestHandler.HandleConfirmarPresenca)
		r.Get("/planos", billingHandler.HandleListarPlanos)            // Nova rota pública
		r.Post("/webhooks/stripe", billingHandler.HandleStripeWebhook) // <-- Rota do Webhook

		// ... outras rotas públicas
		// --- Rotas Protegidas ---
		// Todas as rotas dentro deste grupo exigirão um token JWT válido.
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/casamentos/{idCasamento}/grupos-de-convidados", guestHandler.HandleCriarGrupoDeConvidados)
			r.Get("/acesso-convidado", guestHandler.HandleObterGrupoPorChaveDeAcesso)
			r.Put("/grupos-de-convidados/{idGrupo}", guestHandler.HandleRevisarGrupo)
			// rota de presentes
			r.Post("/casamentos/{idCasamento}/presentes", presenteHandler.HandleCriarPresente)
			r.Post("/selecoes-de-presente", presenteHandler.HandleFinalizarSelecao)
			//  rota de Recados
			r.Post("/recados", recadoHandler.HandleDeixarRecado)
			r.Get("/casamentos/{idCasamento}/recados/admin", recadoHandler.HandleListarRecadosAdmin)
			r.Patch("/recados/{idRecado}", recadoHandler.HandleModerarRecado)
			// rota de Galeria
			r.Post("/casamentos/{idCasamento}/fotos", galleryHandler.HandleFazerUpload)
			r.Get("/casamentos/{idCasamento}/fotos/publico", galleryHandler.HandleListarFotosPublicas)
			r.Post("/fotos/{idFoto}/favoritar", galleryHandler.HandleAlternarFavorito)
			r.Post("/fotos/{idFoto}/rotulos", galleryHandler.HandleAdicionarRotulo)
			r.Delete("/fotos/{idFoto}/rotulos/{nomeDoRotulo}", galleryHandler.HandleRemoverRotulo)
			r.Delete("/fotos/{idFoto}", galleryHandler.HandleDeletarFoto)

			// rotas de eventos
			r.Post("/eventos", eventHandler.HandleCriarEvento)
			r.Post("/assinaturas", billingHandler.HandleCriarAssinatura)

			// r.Get("/eventos/{urlSlug}", eventHandler.HandleObterEventoPorSlug)
			// r.Get("/eventos", eventHandler.HandleListarEventosPorUsuario)

		})
	})

	log.Printf("Servidor iniciado na porta %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
