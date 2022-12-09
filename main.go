package main

import (
	"log"
	"promhsd/db"
	_ "promhsd/docs"
	_ "promhsd/storage/dynamo"
	_ "promhsd/storage/file"
)

// @title        PromHSD
// @version      0.0.1
// @description  prometheus http static config discovery service

// @contact.name  Rinat Almakhov
// @contact.url   https://github.com/Gasoid/

// @license.name  MIT License
// @license.url   https://github.com/Gasoid/promHSD/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /api/

var (
	dbService *db.Service
)

func main() {
	var err error
	storage := getStorage()
	storageArgs := getStorageArgs(storage)
	dbService, err = db.New(storage, storageArgs)
	if err != nil {
		log.Fatal("Can't initialize dbService")
	}
	r := setupRouter()
	r.Run()
}
