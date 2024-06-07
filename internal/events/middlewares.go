package events

import (
	"encoding/json"
	"errors"
	"fmt"

	errs "github.com/vfcoelho/go-project-pocs/internal/errors"
)

func ParseMessage[T any](ctx *ConsumerCtx) error { //REVIEW: standardized parser to prevent code duplication in workers
	var message T
	err := json.Unmarshal(ctx.GetMessage(), &message)
	if err != nil {
		return err
	}
	ctx.SetValue("message", message)
	return ctx.Next()
}

func SetCodeErrorMappings(mappings map[errs.ErrorCode]bool) func(*ConsumerCtx) error { //REVIEW: fiber middleware to set error mappings and later be used by the error response middleware
	return func(ctx *ConsumerCtx) (err error) {
		ctx.SetValue("codeErrorMappings", mappings)
		return ctx.Next()
	}
}

func ErrorRecover(ctx *ConsumerCtx) error { //REVIEW: error handling middleware for workers
	err := ctx.Next()

	if err != nil {
		var customErr errs.Error
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
