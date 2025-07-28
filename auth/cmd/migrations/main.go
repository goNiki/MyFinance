package main

import (
	"auth/internal/config"
	"auth/internal/storage"
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	var migrationPath string
	flag.StringVar(&migrationPath, "migration_path", "", "path to migrations")
	flag.Parse()
	if migrationPath == "" {
		panic("migrationPath is required")
	}

	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		panic(err)
	}

	cfg := config.MustLoad()
	storageInstance, err := storage.New(&cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	db := storageInstance.DB

	_, err = db.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Migration applied")
}
