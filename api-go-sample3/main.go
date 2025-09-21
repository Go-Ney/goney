package main

import (
	"fmt"
	"log"
	"os"
	"api-go-sample3/config"
	"api-go-sample3/pkg/core"
)

func main() {
	cfg := config.Load()
	app := core.NewApplication(cfg)

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Servidor Go-ney iniciado en puerto %s\n", port)
	fmt.Printf("🌐 Visita: http://localhost:%s\n", port)
	fmt.Printf("🩺 Health: http://localhost:%s/api/v1/health\n", port)

	log.Fatal(app.Listen(":" + port))
}
