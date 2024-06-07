package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vfcoelho/go-project-pocs/internal/events"
	"github.com/vfcoelho/go-project-pocs/src/dtos"
	"github.com/vfcoelho/go-project-pocs/src/handlers"
	"github.com/vfcoelho/go-project-pocs/src/repositories"
)

func main() {
	consumer := events.NewConsumer[dtos.Record]()

	handlers := []events.Handler{events.ErrorRecover, events.ParseMessage[dtos.Record], processMessage} //REVIEW: decorator stack of handlers similar to the middleware pattern

	go func() {
		if err := consumer.Consume(handlers...); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Running cleanup tasks...")

	consumer.Close()

	fmt.Println("Fiber was successful shutdown.")
}

func processMessage(ctx *events.ConsumerCtx) error {
	memoryRepository := repositories.NewMemoryRepository[*dtos.Record]() //FIXME: will never succeed because it's not using a shared memory between the producer and the consumer
	return handlers.Consume(ctx, memoryRepository)
}
