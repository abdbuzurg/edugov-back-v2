package main

import (
	"context"
	"edugov-back-v2/internal/config"
	"edugov-back-v2/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appDB "edugov-back-v2/internal/db"
	appHTTP "edugov-back-v2/internal/http"
	handlers "edugov-back-v2/internal/http/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func main() {
	addr := viper.GetString("server.addr")
	dbURL := viper.GetString("db.url")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := appDB.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	store := appDB.NewStore(pool)
	service := service.New(store, validator.New())
	h := handlers.New(service)
	app := appHTTP.NewRouter(h)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server started at port", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}

	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("Server shutdown gracefully")
}

func init() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}
}
