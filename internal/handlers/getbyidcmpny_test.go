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

func TestGetCmpnyById(t *testing.T) {

	nc := models.NewCompany{
		Name:     "teksystem",
		Location: "banglore",
	}
	// fnc := models.NewCompany{
	// 	Name:     "",
	// 	Location: "",
	// }

	mockloc := models.Company{
		Model:    gorm.Model{ID: 1},
		Name:     "teksystem",
		Location: "banglore",
	}

	tt := [...]struct {
		name             string
		body             any
		expectedStatus   int
		expectedResponse string
		expectedUser     models.Company
		mockUserService  func(m *mockmodels.MockService)
	}{
		{
			name:             "ok",
			body:             nc,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"Name":"teksystem","location":"banglore"}`,
			expectedUser:     mockloc,
			//set expection inside this field
			mockUserService: func(m *mockmodels.MockService) {
				m.EXPECT().GetCompanyData(gomock.Any()).
					Times(1).Return(mockloc, nil)
			},
		},
		// {

		// 	name:             "Failure case",
		// 	body:             0,
		// 	expectedStatus:   http.StatusBadRequest,
		// 	expectedResponse: `{"msg":"please provide Name and location"}`,
		// 	expectedUser:     ,
		// 	mockUserService: func(m *mockmodels.MockService) {
		// 		m.EXPECT().GetCompanyData(gomock.Any()).Times(1).Return(nil, errors.New("yjfjyg"))
		// 	},
		// },
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
			router.GET("/getCompanyById/:cmpny_id", h.getCompanyById)
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/getCompanyById/1", bytes.NewReader(body))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)
			require.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
