package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/vfcoelho/go-project-pocs/internal/events"
	"github.com/vfcoelho/go-project-pocs/internal/http"
	"github.com/vfcoelho/go-project-pocs/src/dtos"
)

func main() {
	producer := events.NewProducer[dtos.Record]()

	app := fiber.New()

	http.SetupRouter(app, producer)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Gracefully shutting down...")
	app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	producer.Close()

	fmt.Println("Fiber was successful shutdown.")
}
