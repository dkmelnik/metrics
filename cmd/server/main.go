package main

import (
	"github.com/dkmelnik/metrics/internal/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("server is running!")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s := server.NewServer(":8080")

	if err := s.Run(); err != nil {
		return err
	}

	return nil
}
