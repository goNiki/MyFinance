package main

import (
	"auth/internal/config"
	"auth/internal/storage"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {

	var migrationPath string

	flag.StringVar(&migrationPath, "migration_path", "", "path to migrations")
	flag.Parse()
	if migrationPath == "" {
		panic("migrationPath is required")
	}

	sqlByted, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		panic(err)
	}

	cfg := config.MustLoad()

	storage, err := storage.New(&cfg.DBConfig)
	db := storage.DB

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(context.Background(), string(sqlByted))
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Migration applied")
}
