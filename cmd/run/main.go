package main

import (
	"context"
	"log"
	"net/http"
	"time"

	categoryhandler "github.com/Oleja123/dcaa-category/internal/handler/category"
	categorydb "github.com/Oleja123/dcaa-category/internal/repository/categorymock/db"
	categoryservice "github.com/Oleja123/dcaa-category/internal/service/category"
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

	mux.HandleFunc("/categories", handler.FindAll)
	mux.HandleFunc("/categories/{id}", handler.FindOne)
	mux.HandleFunc("/categories/create", handler.Create)
	mux.HandleFunc("/categories/update", handler.Update)
	mux.HandleFunc("/categories/delete/{id}", handler.Delete)

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	s.ListenAndServe()
}
