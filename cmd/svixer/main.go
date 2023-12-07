package main

import (
	"log"

	"github.com/okulik/gigs-svixer/internal/service"
	"github.com/okulik/gigs-svixer/internal/settings"
)

func main() {
	settings, err := settings.Load()
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(settings)
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
