package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/justsushant/envbox/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEnvService struct {
	mock.Mock
}

func(m *MockEnvService) CreateEnv(client *client.Client, payload types.CreateEnvPayload) (string, error) {
	args := m.Called(client, payload)
	return args.String(0), args.Error(1)
}

func(m *MockEnvService) KillEnv(client *client.Client, id string) error {
	args := m.Called(client, id)
	return args.Error(0)
}

func(m *MockEnvService) GetAllEnvs() ([]types.GetImageResponse, error) {
	args := m.Called()
	return args.Get(0).([]types.GetImageResponse), args.Error(1)
}
func(m *MockEnvService) GetTerminal(*client.Client, string) (dockerTypes.HijackedResponse, error) {
	args := m.Called()
	return args.Get(0).(dockerTypes.HijackedResponse), args.Error(1)
}

func TestCreateEnvHandler(t *testing.T) {
	tt := []struct {
		name string
		payload types.CreateEnvPayload
		mockServiceOutput []interface{}
		expectedStatus int
		expectedResponse string
	}{
		{
			name: "happy path create env",
			payload: types.CreateEnvPayload{
				ImageID: 1,
			},
			mockServiceOutput: []interface{}{"testAccessLink", nil},
			expectedStatus: http.StatusOK,
			expectedResponse: `{
				"status": true,
				"message": "testAccessLink"
			}`,
			
		},
		{
			name: "unhappy path create env",
			payload: types.CreateEnvPayload{
				ImageID: 99,
			},
			mockServiceOutput: []interface{}{"", fmt.Errorf("test error")},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: `{
				"status": false,
				"error": "test error"
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// marshalling payload
			body, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("Unexpected error while marshalling payload: %v", err)
			}

			// setting mocks
			mockService := new(MockEnvService)
			mockClient := &client.Client{}
			mockHandler := NewHandler(mockService, mockClient)
			mockService.On("CreateEnv", mockClient, tc.payload).Return(tc.mockServiceOutput...)

			// setting up gin router
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			// calling the handler
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/createEnv", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")
			mockHandler.RegisterRoutes(router.Group("/"))
			mockHandler.createEnv(c)
			
			// checking the output
			mockService.AssertCalled(t, "CreateEnv", mockClient, tc.payload)
			mockService.AssertExpectations(t)
			
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestGetAllEnvsHandler(t *testing.T) {
	tt := []struct {
		name string
		mockServiceOutput []interface{}
		expectedStatus int
		expectedResponse string
	}{
		{
			name: "happy path get all envs",
			mockServiceOutput: []interface{}{[]types.GetImageResponse{
				{ID: "1", ImageName: "testImage1", AccessLink: "testAccessLink1", CreatedAt: "testCreatedAt1"},
				{ID: "2", ImageName: "testImage2", AccessLink: "testAccessLink2", CreatedAt: "testCreatedAt2"},
			}, nil},
			expectedStatus: http.StatusOK,
			expectedResponse: `{
				"status": true,
				"message": [
					{
						"id": "1",
						"imageName": "testImage1",
						"accessLink": "testAccessLink1",
						"createdAt": "testCreatedAt1"
					},
					{
						"id": "2",
						"imageName": "testImage2",
						"accessLink": "testAccessLink2",
						"createdAt": "testCreatedAt2"
					}
				]
			}`,
		},
		{
			name: "error unhappy path get all envs ",
			expectedResponse: `{
				"status": false,
				"error": "test error"
			}`,
			expectedStatus: http.StatusInternalServerError,
			mockServiceOutput: []interface{}{[]types.GetImageResponse{}, fmt.Errorf("test error")},
		},
		{
			name: "zero unhappy path get all envs ",
			expectedResponse: `{
				"status": false,
				"error": "no envs found"
			}`,
			expectedStatus: http.StatusOK,
			mockServiceOutput: []interface{}{[]types.GetImageResponse{}, nil},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting mocks
			mockService := new(MockEnvService)
			mockClient := &client.Client{}
			mockHandler := NewHandler(mockService, mockClient)
			mockService.On("GetAllEnvs").Return(tc.mockServiceOutput...)

			// setting up gin router
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			// calling the handler
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/getAllEnvs", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			mockHandler.RegisterRoutes(router.Group("/"))
			mockHandler.getAllEnvs(c)
			
			// checking the output
			mockService.AssertCalled(t, "GetAllEnvs")
			mockService.AssertExpectations(t)
			
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestKillEnvHandler(t *testing.T) {
	tt := []struct {
		name string
		id string
		mockServiceOutput []interface{}
		expectedStatus int
		expectedResponse string
	}{
		{
			name: "happy path kill env",
			id: "1",
			mockServiceOutput: []interface{}{nil},
			expectedStatus: http.StatusOK,
			expectedResponse: `{"status": true,"message": "container stopped and removed successfully"}`,
			
		},
		{
			name: "unhappy path kill env",
			id: "99",
			expectedResponse: `{"status": false,"error": "test error"}`,
			expectedStatus: http.StatusInternalServerError,
			mockServiceOutput: []interface{}{fmt.Errorf("test error")},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting mocks
			mockService := new(MockEnvService)
			mockClient := &client.Client{}
			mockHandler := NewHandler(mockService, mockClient)
			mockService.On("KillEnv", mockClient, tc.id).Return(tc.mockServiceOutput...)

			// setting up gin router
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			// calling the handler
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PATCH", "/killEnv/:id", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = []gin.Param{{Key: "id", Value: tc.id}}
			mockHandler.RegisterRoutes(router.Group("/"))
			mockHandler.killEnv(c)
			
			// checking the output
			mockService.AssertCalled(t, "KillEnv", mockClient, tc.id)
			mockService.AssertExpectations(t)
			
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}


// Get Terminal Test Pending

