package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/josesa/servercounter/internal/config"
	"github.com/josesa/servercounter/internal/counter"
	"github.com/josesa/servercounter/internal/server"
	"github.com/josesa/servercounter/internal/service"
	"github.com/josesa/servercounter/internal/storage"
)

func main() {
	// Load configuration object
	config := config.NewFromEnv()

	// Initializing counter
	ct := counter.New(
		counter.WithWindowSize(config.CounterWindowSeconds),
		counter.WithFlushInterval(time.Duration(config.CounterFlushIntervalSeconds)*time.Second),
	)

	// Creating HitCounter Service
	storage := storage.NewFileStorage(config.StoragePath)
	counterService, err := service.New(ct, storage)
	if err != nil {
		log.Fatal(err)
	}

	// Creating the HTTP Server
	ws := server.New(counterService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /request", ws.Request)
	server := http.Server{
		Addr:    config.ServerAddress,
		Handler: mux,
	}

	// Create handler for terminating the webserver
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-terminate
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("SIGTERM received, terminating HTTP server")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// HTTP server has been stopped, ensure service is correctly terminated
	err = counterService.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
