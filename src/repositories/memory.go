package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	internal "github.com/vfcoelho/go-studies/go-api/internal"
	errs "github.com/vfcoelho/go-studies/go-api/internal/errors"
)

type RecordInterface interface {
	ID() uuid.UUID
	SetID(uuid.UUID)
}

var recordNotFound = errors.New("record not found")

type Records[T RecordInterface] map[uuid.UUID]T

type MemoryRepository[T RecordInterface] struct {
	records Records[T]
}

func NewMemoryRepository[T RecordInterface]() *MemoryRepository[T] {
	result := new(MemoryRepository[T])
	result.records = make(Records[T])
	return result
}

func (mr *MemoryRepository[T]) Get(id uuid.UUID) (record T, err error) {
	record, ok := mr.records[id]
	if !ok {
		err = recordNotFound
	}
	return
}

func (mr *MemoryRepository[T]) Add(record T) (err error) {
	if _, ok := mr.records[record.ID()]; ok {
		return errs.NewError(errors.New("id already exists"), errs.WithCode(internal.RECORD_ALREADY_EXISTS_ERROR))
	}
	if record.ID() == uuid.Nil {
		return errors.New("id cannot be nil")
	}
	record.SetID(lo.Ternary(record.ID() == uuid.Nil, uuid.New(), record.ID()))
	mr.records[record.ID()] = record
	fmt.Printf("Record added: %v\n", record)
	return
}

func (mr *MemoryRepository[T]) Update(record T) (err error) {
	if _, ok := mr.records[record.ID()]; !ok {
		return recordNotFound
	}

	mr.records[record.ID()] = record
	fmt.Printf("Record updated: %v\n", record)
	return
}
