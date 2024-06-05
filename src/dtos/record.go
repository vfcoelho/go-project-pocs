package dtos

import "github.com/google/uuid"

type Status string

const (
	pending   Status = "pending"
	processed Status = "processed"
)

type Record struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status Status    `json:"status"`
}

func NewRecord() *Record {
	return &Record{
		Status: pending,
	}
}

func (r *Record) SetProcessed() {
	r.Status = processed
}

func (r Record) ID() uuid.UUID {
	return r.Id
}

func (r *Record) SetID(id uuid.UUID) {
	r.Id = id
}
