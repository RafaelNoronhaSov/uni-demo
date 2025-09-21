package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"uni-demo/pkg/api"
	"uni-demo/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "uni-demo/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Uni-Demo API
// @version 1.0
// @description API for university demo project
// @host localhost:8080
// @BasePath /
func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := db.Connect(dbURL)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco: %v\n", err)
	}
	defer pool.Close()
	log.Println("Conexão realizada com sucesso!")

	if err := db.Migrate(pool); err != nil {
		log.Fatalf("Falha rodar as migrations: %v\n", err)
	}
	log.Println("Migrations executadas com sucesso!")

	if err := db.Seed(pool); err != nil {
		log.Fatalf("Falha ao popular banco: %v\n", err)
	}
	log.Println("Banco populado com sucesso!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiHandler := api.NewHandler(pool)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format(time.RFC3339)
		w.Write([]byte(fmt.Sprintf("ok - %s\n", timestamp)))
	})

	r.Get("/professor-hours", apiHandler.GetProfessorHours)
	r.Get("/room-schedules", apiHandler.GetRoomSchedules)

	// Swagger route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Inicializando servidor na porta 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v\n", err)
	}
}
