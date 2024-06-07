package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorsTestSuite struct {
	suite.Suite
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

func (suite *ErrorsTestSuite) TestCustomErrorIs() {

	sentinelError := errors.New("test error")
	otherSentinelError := errors.New("other test error")
	const TEST_CODE ErrorCode = "TEST_CODE"
	const OTHER_TEST_CODE ErrorCode = "OTHER_TEST_CODE"

	type testCase struct {
		Error  error
		Target error
	}

	TestCase := func(name string, useCase testCase) (string, func()) {
		return name, func() {
			assert.ErrorIs(suite.T(), useCase.Error, useCase.Target)
		}
	}

	suite.Run(TestCase("custom error is wrapped error", testCase{
		Error:  NewError(sentinelError),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error is original error even when wrapped by fmt package", testCase{
		Error:  NewError(fmt.Errorf("inner wrapper: %w", sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("fmt package wrapped custom error is original error", testCase{
		Error:  fmt.Errorf("outer wrapper: %w", NewError(sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error wrapped custom error is original error", testCase{
		Error:  NewError(NewError(sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error with code wrapped custom error is original error", testCase{
		Error:  NewError(NewError(sentinelError), WithCode(TEST_CODE)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error wrapped custom error with code is original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE))),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error with code wrapped custom error is custom error with code wrapper regardless of original error ", testCase{
		Error:  NewError(NewError(sentinelError), WithCode(TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error wrapped custom error wrapper with code is error wrapper with code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE))),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error with another code wrapped custom error wrapper with code is error wrapper with code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE)), WithCode(OTHER_TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error with another code wrapped custom error wrapper with code is error wrapper with another code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE)), WithCode(OTHER_TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(OTHER_TEST_CODE)),
	}))
	suite.Run(TestCase("fmt package wrapped custom error wrapper with code is custom wrapper with code regardless of original error", testCase{
		Error:  fmt.Errorf("raw error wrap: %w", NewError(sentinelError, WithCode(TEST_CODE))),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
}

func (suite *ErrorsTestSuite) TestCustomErrorAs() {

	sentinelError := errors.New("test error")
	otherSentinelError := errors.New("other test error")
	const TEST_CODE ErrorCode = "TEST_CODE"
	const OTHER_TEST_CODE ErrorCode = "OTHER_TEST_CODE"

	type testCase struct {
		Error  error
		Target any
	}

	TestCase := func(name string, useCase testCase) (string, func()) {
		return name, func() {
			assert.ErrorAs(suite.T(), useCase.Error, &useCase.Target)
		}
	}

	suite.Run(TestCase("custom error as wrapped error", testCase{
		Error:  NewError(sentinelError),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error as original error even when wrapped by fmt package", testCase{
		Error:  NewError(fmt.Errorf("inner wrapper: %w", sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("fmt package wrapped custom error as original error", testCase{
		Error:  fmt.Errorf("outer wrapper: %w", NewError(sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error wrapped custom error as original error", testCase{
		Error:  NewError(NewError(sentinelError)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error with code wrapped custom error as original error", testCase{
		Error:  NewError(NewError(sentinelError), WithCode(TEST_CODE)),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error wrapped custom error with code as original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE))),
		Target: sentinelError,
	}))
	suite.Run(TestCase("custom error with code wrapped custom error as custom error with code wrapper regardless of original error ", testCase{
		Error:  NewError(NewError(sentinelError), WithCode(TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error wrapped custom error wrapper with code as error wrapper with code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE))),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error with another code wrapped custom error wrapper with code as error wrapper with code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE)), WithCode(OTHER_TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("custom error with another code wrapped custom error wrapper with code as error wrapper with another code regardless of original error", testCase{
		Error:  NewError(NewError(sentinelError, WithCode(TEST_CODE)), WithCode(OTHER_TEST_CODE)),
		Target: NewError(otherSentinelError, WithCode(OTHER_TEST_CODE)),
	}))
	suite.Run(TestCase("fmt package wrapped custom error wrapper with code as custom wrapper with code regardless of original error", testCase{
		Error:  fmt.Errorf("outer wrapper: %w", NewError(sentinelError, WithCode(TEST_CODE))),
		Target: NewError(otherSentinelError, WithCode(TEST_CODE)),
	}))
	suite.Run(TestCase("01", testCase{
		Error:  fmt.Errorf("outer wrapper: %w", fmt.Errorf("inner wrapper: %w", sentinelError)),
		Target: lo.ToPtr(fmt.Errorf("other wrapper: %w", otherSentinelError)),
	}))

}
