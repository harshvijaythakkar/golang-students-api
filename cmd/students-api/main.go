package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harshvijaythakkar/golang-students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()


	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	// setup server
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("address", cfg.Addr))

	// chan to handle graful shutdwon
	done := make(chan os.Signal, 1)

	// send signal in done chan
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// run parallel threads
	go func ()  {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	} ()

	// receieve msg, wait for os signal msg
	<-done

	slog.Info("Shutting down the server")

	// context to handle shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	// call shutdown method with 5 sec context
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	
	slog.Info("Server shutdown successfully")
}

