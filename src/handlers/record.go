package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vfcoelho/go-project-pocs/internal"
	errs "github.com/vfcoelho/go-project-pocs/internal/errors"
	"github.com/vfcoelho/go-project-pocs/internal/events"
	"github.com/vfcoelho/go-project-pocs/src/dtos"
)

type RecordRepository[T any] interface {
	Get(id uuid.UUID) (record T, err error)
	Add(record T) error
	Update(record T) error
}

type EventProducer[T any] interface {
	Send(event T) error
}

func Get(c *fiber.Ctx, recordRepository RecordRepository[*dtos.Record]) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err //REVIEW: the ErrorRecoverMiddleware can handle go raw errors, returning 500 and the error text to the caller
	}

	record, err := recordRepository.Get(id)
	if err != nil {
		return errs.NewError(err, errs.WithCode(internal.RECORD_NOT_FOUND_ERROR), errs.WithData(struct {
			ID uuid.UUID `json:"id"`
		}{ID: id})) //REVIEW: the custom error can be used to return data to the caller
	}

	return c.JSON(record)
}

func Post(c *fiber.Ctx, recordRepository RecordRepository[*dtos.Record], producer EventProducer[dtos.Record]) error {

	payload := dtos.NewRecord()
	if err := c.BodyParser(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("error parsing payload: %w", err).Error()) //REVIEW: the ErrorRecoverMiddleware can handle fiber errors, that define an error code. Tough advised to use the error mapping instead, it can be useful for backward compatibility
	}

	if err := recordRepository.Add(payload); err != nil {
		if errors.Is(err, errs.NewIsComparable(internal.RECORD_ALREADY_EXISTS_ERROR)) { //REVIEW: the custom error can be used to treat specific error codes that we might not want to return to the caller or cause the application to break loop
			return fmt.Errorf("error adding record: %w", err) //REVIEW: the ErrorRecoverMiddleware can still identify the underlying error if it was wrapped.
		}
		return err
	}
	err := producer.Send(*payload)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}

func Consume(c *events.ConsumerCtx, recordRepository RecordRepository[*dtos.Record]) error {

	record := c.GetValue("message").(dtos.Record)
	record.SetProcessed()

	err := recordRepository.Update(&record)
	err = errs.NewError(err, errs.WithData(record))

	return err
}
