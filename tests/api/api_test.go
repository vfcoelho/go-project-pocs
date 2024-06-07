package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	errs "github.com/vfcoelho/go-project-pocs/internal/errors"
	"github.com/vfcoelho/go-project-pocs/internal/events"
	"github.com/vfcoelho/go-project-pocs/internal/http"
	"github.com/vfcoelho/go-project-pocs/src/dtos"
)

type ApiTestSuite struct {
	suite.Suite
	app *fiber.App
}

// TODO: add mocks for event producer
func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (suite *ApiTestSuite) SetupTest() {
	suite.app = fiber.New()

	producer := events.NewProducer[dtos.Record]()
	http.SetupRouter(suite.app, producer)
}

func (suite *ApiTestSuite) TestCreateRecord() {

	type testResponse struct {
		errs.Error
	}
	type record struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	type testCase struct {
		Payload  record
		WantCode int
		Want     testResponse
	}

	TestCase := func(name string, useCase testCase) (string, func()) {
		return name, func() {
			payload, _ := json.Marshal(useCase.Payload)
			route := "/v1/record"
			req := httptest.NewRequest("POST", route, bytes.NewBuffer(payload))
			req.Header.Add("Content-Type", "application/json")

			resp, _ := suite.app.Test(req, -1)
			bodyString, _ := io.ReadAll(resp.Body)
			var body testResponse
			err := json.Unmarshal(bodyString, &body)
			_ = err

			assert.Equal(suite.T(), useCase.WantCode, resp.StatusCode)
			assert.Equal(suite.T(), useCase.Want, body)
		}
	}

	suite.Run(TestCase("should fail while creating record with nil id", testCase{
		Payload: record{
			Id:   uuid.Nil.String(),
			Name: "Dummy Record",
		},
		WantCode: 500,
		Want: testResponse{
			Error: errs.Error{Err: errors.New("id cannot be nil")},
		},
	}))
	suite.Run(TestCase("should fail while creating record with invalid id", testCase{
		Payload: record{
			Id:   "5d2ca371-f623-4aac-abb0-ddc31f44d00-",
			Name: "Dummy Record",
		},
		WantCode: 400,
		Want: testResponse{
			Error: errs.Error{Err: errors.New("error parsing payload: invalid UUID format")},
		},
	}))
	suite.Run(TestCase("should succeed while creating record with valid id", testCase{
		Payload: record{
			Id:   "5d2ca371-f623-4aac-abb0-ddc31f44d002",
			Name: "Dummy Record",
		},
		WantCode: 201,
		Want:     testResponse{},
	}))
	suite.Run(TestCase("should fail while creating record with repeated id", testCase{
		Payload: record{
			Id:   "5d2ca371-f623-4aac-abb0-ddc31f44d002",
			Name: "Dummy Record",
		},
		WantCode: 409,
		Want: testResponse{
			Error: errs.Error{
				Err:  errors.New("id already exists"),
				Code: "record_already_exists",
			},
		},
	}))
}
