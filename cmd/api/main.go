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
	"github.com/luiszkm/wedding_backend/internal/guest/application"
	"github.com/luiszkm/wedding_backend/internal/guest/infrastructure"
	"github.com/luiszkm/wedding_backend/internal/guest/interfaces/rest"
)

func main() {
	dbURL := "postgres://user:password@localhost:5432/wedding_db?sslmode=disable"
	port := ":3000"

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Incapaz de conectar ao banco de dados: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// --- Composição da Raiz (Wiring) ---
	// 1. Instancia a implementação da infraestrutura
	guestRepo := infrastructure.NewPostgresGroupRepository(dbpool)
	// 2. Instancia o serviço de aplicação, injetando a implementação do repositório
	guestService := application.NewGuestService(guestRepo)
	// 3. Instancia o handler, injetando o serviço
	guestHandler := rest.NewGuestHandler(guestService)

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
		// Nova rota para editar grupo de convidados
		r.Put("/grupos-de-convidados/{idGrupo}", guestHandler.HandleRevisarGrupo)

	})

	log.Printf("Servidor iniciado na porta %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
