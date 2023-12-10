package main

import (
	"errors"
	"flag"
	"github.com/dkmelnik/metrics/internal/handlers"
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
	addr := flag.String("a", "localhost:8080", "in the form host:port. If empty, :8080 is used")

	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("unknown flags or parameters")
	}

	r := handlers.ConfigureRouter()
	s := server.NewServer(*addr, r)

	if err := s.Run(); err != nil {
		return err
	}

	return nil
}
