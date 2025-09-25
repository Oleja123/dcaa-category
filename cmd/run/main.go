package main

import (
	"context"
	"log"
	"net/http"
	"time"

	categoryservice "github.com/Oleja123/dcaa-category/internal/application/category"
	categoryhandler "github.com/Oleja123/dcaa-category/internal/handler/category"
	categorydb "github.com/Oleja123/dcaa-category/internal/infrastructure/category/db"
	"github.com/Oleja123/dcaa-property/pkg/client/postgresql"
	"github.com/Oleja123/dcaa-property/pkg/config"
)

func main() {
	config, err := config.LoadConfig("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := postgresql.NewClient(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	repository := categorydb.NewRepository(client)

	service := categoryservice.NewService(repository)

	handler := categoryhandler.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/categories", handler.Handle)
	mux.HandleFunc("/categories/{id}", handler.HandleWithId)

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	s.ListenAndServe()
}
