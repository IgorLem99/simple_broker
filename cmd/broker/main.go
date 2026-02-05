package main

import (
	"log"

	"github.com/IgorLem99/simple_broker/internal/broker"
	"github.com/IgorLem99/simple_broker/internal/config"
	"github.com/IgorLem99/simple_broker/internal/server"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	b := broker.New(cfg)
	srv := server.New(cfg.Addr, b)

	log.Printf("starting server on %s", cfg.Addr)
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
