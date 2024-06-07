package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

//REVIEW: the errors package is built to be used as a potential standalone package for error handling in go projects

type ErrorCode string //REVIEW: defines type to create error codes

type Error struct { //REVIEW: defines a custom error type to add functionality fo debugging, logging and responding API calls
	Err  error     `json:"error,omitempty"`
	Code ErrorCode `json:"code,omitempty"`
	Data any       `json:"data,omitempty"`
}

type payloadError struct {
	Err  string    `json:"error,omitempty"`
	Code ErrorCode `json:"code,omitempty"`
	Data any       `json:"data,omitempty"`
}

func (he Error) Error() string { //REVIEW: custom error is a go error
	return fmt.Errorf("%v: %w", he.Code, he.Err).Error()
}

func (he Error) Is(target error) bool { //REVIEW: Is method evaluates error code in order to match different data and wrapped errors
	var targetError Error
	if errors.As(target, &targetError) {
		return targetError.Code == he.Code
	}
	return errors.Is(he.Err, target)
}

func (he Error) MarshalJSON() ([]byte, error) { //REVIEW: also marshalable to json for later logging
	return json.Marshal(payloadError{
		Err:  he.Err.Error(),
		Code: he.Code,
		Data: he.Data,
	})
}

func (he Error) Unwrap() error {
	return he.Err
}

func (he *Error) UnmarshalJSON(data []byte) error { //REVIEW: also unmarshalable from json for API usage
	var s payloadError
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*he = NewError(errors.New(s.Err), WithCode(s.Code), WithData(s.Data))
	return nil
}

func NewIsComparable(code ErrorCode) Error {
	return NewError(nil, WithCode(code))
}

type ErrorOption func(*Error) //REVIEW: provides with-builder methods for convenience
func WithCode(code ErrorCode) ErrorOption {
	return func(he *Error) {
		he.Code = code
	}
}
func WithData(data any) ErrorOption {
	return func(he *Error) {
		he.Data = data
	}
}
func NewError(err error, opts ...ErrorOption) Error {
	he := Error{Err: err}
	for _, opt := range opts {
		opt(&he)
	}
	return he
}
