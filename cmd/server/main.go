package main

import (
	"id-generator/internal/app"
	"id-generator/internal/delivery/httphandler"
	"id-generator/internal/infra/snowflake"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	gen, err := snowflake.NewSnowflakeGenerator(1)
	if err != nil {
		log.Fatalf("error initializing snowflake: %v", err)
	}
	service := app.NewIDService(gen)
	handler := httphandler.NewHandler(service)

	mux.HandleFunc("/message-id", handler.GetID)

	log.Println("Message ID Service running on port:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
