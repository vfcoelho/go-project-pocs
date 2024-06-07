package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vfcoelho/go-project-pocs/internal/errors"
)

const (
	RECORD_ALREADY_EXISTS_ERROR errors.ErrorCode = "record_already_exists"
	RECORD_NOT_FOUND_ERROR      errors.ErrorCode = "record_not_found"
)

var MAPPING = map[errors.ErrorCode]int{ //REVIEW: possible errors and their default codes would be defined in a separate file in order to be used by the whole project
	RECORD_NOT_FOUND_ERROR:      fiber.StatusNotFound,
	RECORD_ALREADY_EXISTS_ERROR: fiber.StatusConflict,
}
