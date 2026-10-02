package main

import (
	"fmt"
	"log"
	"os"
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

	var inbounds []database.Inbound

	if err := database.DB.
		Find(&inbounds).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Found %d inbound(s) in database\n\n",
		len(inbounds),
	)

	config, err := xray.BuildConfig(inbounds)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(
		"/tmp/ttcs-xray.json",
		config,
		0644,
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"Xray config written to /tmp/ttcs-xray.json",
	)
}
