package main

import (
	"log"

	gwook "github.com/okulik/gwook/internal"
	"github.com/okulik/gwook/internal/service"
	"github.com/okulik/gwook/internal/settings"
)

func main() {
	settings, err := settings.Load()
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(settings, gwook.NewSvixAdapter(settings))
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
