package main

import (
	"log"

	"github.com/DmitriiKumancev/mongoapi/internal/controllers"
	"github.com/DmitriiKumancev/mongoapi/internal/repository"
	"github.com/DmitriiKumancev/mongoapi/pkg/database"
	"github.com/DmitriiKumancev/mongoapi/pkg/router"
)

func main() {

	uri := "mongodb://root:pass@localhost:27017"
	mongoDb, err := database.NewMongoDatabase(uri)
	if err != nil {
		log.Fatal(err.Error())
	}

	repo := repository.New(mongoDb)

	controller := controllers.New(repo)

	r := router.Initialize(controller)

	r.Start(":3030")
}
