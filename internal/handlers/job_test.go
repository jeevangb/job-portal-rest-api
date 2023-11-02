package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	middlewares "jobportalapi/internal/middlerwares"
	"jobportalapi/internal/models"
	"jobportalapi/internal/services"
	"jobportalapi/internal/services/mockmodels"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestCreateJob(t *testing.T) {

	nj := models.NewJob{
		Title: "Developer",
		Desc:  "golang developer",
	}

	fnj := models.NewJob{
		Title: "",
		Desc:  "",
	}
	mockjob := models.Job{
		Model: gorm.Model{ID: 1},
		Cid:   1,
		Title: "Developer",
		Desc:  "golang developer",
	}

	tt := [...]struct {
		name             string
		body             any
		expectedStatus   int
		expectedResponse string
		expectedUser     models.Job
		mockUserService  func(m *mockmodels.MockService)
	}{
		{
			name:             "success case",
			body:             nj,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"cid":1,"title":"Developer","desc":"golang developer"}`,
			expectedUser:     mockjob,
			//set expection inside this field
			mockUserService: func(m *mockmodels.MockService) {
				m.EXPECT().StoreJob(gomock.Any(), gomock.Eq(nj), uint64(1)).
					Times(1).Return(mockjob, nil)
			},
		},
		{

			name:             "Failure no job creation",
			body:             fnj,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"msg":"please provide title and desc"}`,
			mockUserService: func(m *mockmodels.MockService) {
				m.EXPECT().StoreJob(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new Gomock controller.
			ctrl := gomock.NewController(t)
			//this func give us the mock implementation of the interface
			mockService := mockmodels.NewMockService(ctrl)
			s := services.NewStore(mockService)

			// Apply the mock to the user service.
			tc.mockUserService(mockService)

			router := gin.New()
			h := handler{
				S: s,
			}
			ctx := context.Background()
			traceID := "fake-trace-id"
			ctx = context.WithValue(ctx, middlewares.TraceIdKey, traceID)

			//register endpoints
			router.POST("/createJob/:cmpny_id", h.createJob)
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/createJob/1", bytes.NewReader(body))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)
			require.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
