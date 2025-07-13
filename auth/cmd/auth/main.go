package main

import (
	"auth/internal/config"
	"auth/internal/handlers/register"
	mwLogger "auth/internal/middlware/logger"
	"auth/internal/storage"
	"auth/internal/utils/logger"
	"auth/internal/utils/logger/sl"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	//TODO init config
	cfg := config.MustLoad()

	//TODO init loger
	log := logger.InitLogger(cfg.Env)
	log.Info("Init Logger")
	//TODO inti storage
	storage, err := storage.New(&cfg.DBConfig)
	if err != nil {
		log.Error("failed to init db ", sl.Err(err))
		return
	}
	fmt.Println(storage)
	log.Info("Connected to PostgreSQL")

	//TODO init router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Post("/user", register.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	//TODO init Server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}
