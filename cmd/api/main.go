// file: cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	communicationApp "github.com/luiszkm/wedding_backend/internal/communication/application"
	communicationInfra "github.com/luiszkm/wedding_backend/internal/communication/infrastructure"
	communicationREST "github.com/luiszkm/wedding_backend/internal/communication/interfaces/rest"

	itineraryApp "github.com/luiszkm/wedding_backend/internal/itinerary/application"
	itineraryInfra "github.com/luiszkm/wedding_backend/internal/itinerary/infrastructure"
	itineraryREST "github.com/luiszkm/wedding_backend/internal/itinerary/interfaces/rest"

	pageTemplateApp "github.com/luiszkm/wedding_backend/internal/pagetemplate/application"
	pageTemplateREST "github.com/luiszkm/wedding_backend/internal/pagetemplate/interfaces/rest"
	"github.com/luiszkm/wedding_backend/internal/platform/template"
)

func main() {
	port := ":8080"
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

	// CORS configuration
	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsAllowedOrigins == "" {
		corsAllowedOrigins = "http://localhost:3000,http://localhost:3001,https://localhost:3000"
	}
	corsAllowedMethods := os.Getenv("CORS_ALLOWED_METHODS")
	if corsAllowedMethods == "" {
		corsAllowedMethods = "GET,POST,PUT,DELETE,OPTIONS"
	}
	corsAllowedHeaders := os.Getenv("CORS_ALLOWED_HEADERS")
	if corsAllowedHeaders == "" {
		corsAllowedHeaders = "Accept,Authorization,Content-Type,X-CSRF-Token,X-Requested-With"
	}
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
	templateEngine := template.NewGoTemplateEngine("templates")

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
	communicationRepo := communicationInfra.NewPostgresComunicadoRepository(dbpool)
	itineraryRepo := itineraryInfra.NewPostgresItineraryRepository(dbpool)

	// --- Serviços de Aplicação ---
	guestService := guestApp.NewGuestService(guestRepo)
	presenteService := giftApp.NewGiftService(presenteRepo, selecaoRepo, eventRepo)
	recadoService := mbApp.NewMessageBoardService(recadoRepo, guestRepo, eventRepo)
	galleryService := galleryApp.NewGalleryService(fotoRepo, storageSvc)
	iamService := iamApp.NewIAMService(usuarioRepo, jwtService)
	eventService := eventApp.NewEventService(eventRepo)
	billingService := billingApp.NewBillingService(planoRepo, billingRepo, paymentGateway)
	communicationService := communicationApp.NewCommunicationService(communicationRepo, eventRepo)
	itineraryService := itineraryApp.NewItineraryService(itineraryRepo)
	pageTemplateService := pageTemplateApp.NewPageTemplateService(templateEngine, eventRepo, guestRepo, presenteRepo, recadoRepo, fotoRepo)

	// --- Handlers ---
	guestHandler := guestREST.NewGuestHandler(guestService)
	presenteHandler := giftREST.NewGiftHandler(presenteService, storageSvc)
	recadoHandler := mbREST.NewMessageBoardHandler(recadoService)
	galleryHandler := galleryREST.NewGalleryHandler(galleryService)
	iamHandler := iamREST.NewIAMHandler(iamService)
	eventHandler := eventREST.NewEventHandler(eventService)
	billingHandler := billingREST.NewBillingHandler(billingService, stripeWebhookSecret)
	communicationHandler := communicationREST.NewCommunicationHandler(communicationService)
	itineraryHandler := itineraryREST.NewItineraryHandler(itineraryService)
	pageTemplateHandler := pageTemplateREST.NewPageTemplateHandler(pageTemplateService)

	// --- Roteador e Rotas ---
	r := chi.NewRouter()

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     strings.Split(corsAllowedOrigins, ","),
		AllowedMethods:     strings.Split(corsAllowedMethods, ","),
		AllowedHeaders:     strings.Split(corsAllowedHeaders, ","),
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	authMiddleware := auth.Authenticator(jwtService)

	r.Route("/v1", func(r chi.Router) {
		// --- Rotas Públicas ---
		r.Post("/usuarios/registrar", iamHandler.HandleRegistrar)
		r.Post("/usuarios/login", iamHandler.HandleLogin)
		r.Get("/eventos/{idCasamento}/recados/publico", recadoHandler.HandleListarRecadosPublicos)
		r.Get("/eventos/{idCasamento}/presentes-publico", presenteHandler.HandleListarPresentesPublicos)
		r.Get("/eventos/{idEvento}/comunicados", communicationHandler.HandleListarComunicados)
		r.Get("/eventos/{idEvento}/roteiro", itineraryHandler.HandleGetItinerary)         // Rota pública do roteiro
		r.Get("/eventos/{urlSlug}/pagina", pageTemplateHandler.HandleRenderPublicPage)    // Página pública do evento
		r.Get("/templates/disponiveis", pageTemplateHandler.HandleListAvailableTemplates) // Templates disponíveis
		r.Post("/rsvps", guestHandler.HandleConfirmarPresenca)
		r.Get("/planos", billingHandler.HandleListarPlanos)            // Nova rota pública
		r.Post("/webhooks/stripe", billingHandler.HandleStripeWebhook) // <-- Rota do Webhook

		// ... outras rotas públicas
		// --- Rotas Protegidas ---
		// Todas as rotas dentro deste grupo exigirão um token JWT válido.
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/eventos/{idCasamento}/grupos-de-convidados", guestHandler.HandleCriarGrupoDeConvidados)
			r.Get("/acesso-convidado", guestHandler.HandleObterGrupoPorChaveDeAcesso)
			r.Put("/grupos-de-convidados/{idGrupo}", guestHandler.HandleRevisarGrupo)
			// rota de presentes
			r.Post("/eventos/{idCasamento}/presentes", presenteHandler.HandleCriarPresente)
			r.Post("/selecoes-de-presente", presenteHandler.HandleFinalizarSelecao)
			//  rota de Recados
			r.Post("/recados", recadoHandler.HandleDeixarRecado)
			r.Get("/eventos/{idCasamento}/recados/admin", recadoHandler.HandleListarRecadosAdmin)
			r.Patch("/recados/{idRecado}", recadoHandler.HandleModerarRecado)
			// rota de Comunicados
			r.Post("/eventos/{idEvento}/comunicados", communicationHandler.HandleCriarComunicado)
			r.Put("/comunicados/{idComunicado}", communicationHandler.HandleEditarComunicado)
			r.Delete("/comunicados/{idComunicado}", communicationHandler.HandleDeletarComunicado)
			// rotas de Roteiro/Itinerary (autenticadas)
			r.Post("/eventos/{idEvento}/roteiro", itineraryHandler.HandleCreateItineraryItem)
			r.Put("/roteiro/{idItemRoteiro}", itineraryHandler.HandleUpdateItineraryItem)
			r.Delete("/roteiro/{idItemRoteiro}", itineraryHandler.HandleDeleteItineraryItem)
			// rota de Galeria
			r.Post("/eventos/{idCasamento}/fotos", galleryHandler.HandleFazerUpload)
			r.Get("/eventos/{idCasamento}/fotos/publico", galleryHandler.HandleListarFotosPublicas)
			r.Post("/fotos/{idFoto}/favoritar", galleryHandler.HandleAlternarFavorito)
			r.Post("/fotos/{idFoto}/rotulos", galleryHandler.HandleAdicionarRotulo)
			r.Delete("/fotos/{idFoto}/rotulos/{nomeDoRotulo}", galleryHandler.HandleRemoverRotulo)
			r.Delete("/fotos/{idFoto}", galleryHandler.HandleDeletarFoto)

			// rotas de eventos
			r.Post("/eventos", eventHandler.HandleCriarEvento)
			r.Post("/assinaturas", billingHandler.HandleCriarAssinatura)

			r.Get("/eventos/{urlSlug}", eventHandler.HandleObterEventoPorSlug)
			r.Get("/eventos", eventHandler.HandleListarEventosPorUsuario)

			// rotas de templates
			r.Put("/eventos/{eventId}/template", pageTemplateHandler.HandleUpdateEventTemplate)
			r.Get("/templates/{templateId}", pageTemplateHandler.HandleGetTemplateMetadata)
			r.Post("/templates/preview", pageTemplateHandler.HandlePreviewTemplate)

		})
	})

	log.Printf("Servidor iniciado na porta %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
