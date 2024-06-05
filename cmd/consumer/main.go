package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	errs "github.com/vfcoelho/go-studies/go-api/internal/errors"
	"github.com/vfcoelho/go-studies/go-api/internal/events"
	"github.com/vfcoelho/go-studies/go-api/src/dtos"
	"github.com/vfcoelho/go-studies/go-api/src/handlers"
	"github.com/vfcoelho/go-studies/go-api/src/repositories"
)

func main() {
	consumer := events.NewConsumer[dtos.Record]()

	handlers := []events.Handler{errorRecover, parseMessage, processMessage} //REVIEW: decorator stack of handlers similar to the middleware pattern

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
func parseMessage(ctx *events.ConsumerCtx) error { //REVIEW: standardized parser to prevent code duplication in workers
	var message dtos.Record
	err := json.Unmarshal(ctx.GetMessage(), &message)
	if err != nil {
		return err
	}
	ctx.SetValue("message", message)
	return ctx.Next()
}

func errorRecover(ctx *events.ConsumerCtx) error { //REVIEW: error handling middleware for workers
	err := ctx.Next()
	if err != nil {
		var customErr *errs.Error
		switch {
		case errors.As(err, &customErr):
			stringErr, err := json.Marshal(customErr)
			if err != nil {
				return fmt.Errorf("error marshalling custom error: %w", err)
			}
			fmt.Println(string(stringErr))
			return nil
		default:
			return err
		}
	}
	return nil
}
