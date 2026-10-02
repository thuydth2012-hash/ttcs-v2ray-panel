package main

import (
	"log"

	"ttcs/database"
	"ttcs/xray"
)

func main() {

	if err := database.Init(
		database.DefaultDBPath,
	); err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	if err := xray.WriteConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println(
		"Xray config written successfully",
	)
}
