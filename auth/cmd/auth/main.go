package main

import (
	"auth/internal/config"
	"auth/internal/storage"
	"auth/internal/utils/logger"
	"auth/internal/utils/logger/sl"
	"fmt"
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

	//TODO init Server

}
