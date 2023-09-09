package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/auth"
	"url-shortener/internal/config"
	"url-shortener/internal/db"
	"url-shortener/internal/github"
	"url-shortener/internal/server"
	"url-shortener/internal/shorten"
	"url-shortener/internal/storage/shortening"
	"url-shortener/internal/storage/user"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("unable to load configuration")
	}
}

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer dbCancel()

	mgoClient, err := db.Connect(dbCtx, config.Get().DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	mgoDB := mgoClient.Client().Database(config.Get().DB.Database)

	var (
		shorteningStorage = shortening.NewMongoDB(mgoDB)
		userStorage       = user.NewMongoDB(mgoDB)
		shortener         = shorten.NewService(shorteningStorage)
		githubClient      = github.NewClient()
		authenticator     = auth.NewService(
			githubClient,
			userStorage,
			config.Get().Github.ClientID,
			config.Get().Github.ClientSecret,
		)
		srv = server.New(shortener, authenticator)
	)

	srv.AddCloser(mgoClient.Close)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(config.Get().ListenAddr(), srv); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error running server: %v", err)
		}
	}()

	log.Println("server started")
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}

	log.Println("server stopped")
}
