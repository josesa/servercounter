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

	"github.com/josesa/servercounter/internal/counter"
	"github.com/josesa/servercounter/internal/server"
	"github.com/josesa/servercounter/internal/service"
)

type fakeTime struct{}

func (ft fakeTime) Now() time.Time {
	fake := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	return fake
}

func main() {
	// Initializing counter
	c := counter.New(
		counter.WithWindowSize(10),
		counter.WithFlushInterval(12*time.Second),
		// counter.WithTime(fakeTime{}),
	)

	// Creating HitCounter Service
	counterService, err := service.New(c, "data.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Creating the HTTP Server
	ws := server.New(counterService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /request", ws.Request)
	server := http.Server{
		Addr:    ":8080",
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
	err = counterService.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
}
