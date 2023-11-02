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
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5" // Import the correct JWT package

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserLogin(t *testing.T) {

	nl := models.Login{
		Email:    "jeevan@gmail.com",
		Password: "password",
	}

	mocklogin := models.User{
		Model:        gorm.Model{ID: 1},
		Name:         "jeevan",
		Email:        "jeevan@gmail.com",
		PasswordHash: "jskskslsms",
	}
	var u models.User
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"students"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tt := [...]struct {
		name             string
		body             any
		expectedStatus   int
		expectedResponse string
		expectedUser     models.User
		mockUserService  func(m *mockmodels.MockService)
	}{
		{
			name:             "login success",
			body:             nl,
			expectedStatus:   http.StatusOK,
			expectedResponse: ``,
			expectedUser:     mocklogin,
			//set expection inside this field
			mockUserService: func(m *mockmodels.MockService) {
				m.EXPECT().Authenticate(gomock.Any(), gomock.Eq(nl.Email), gomock.Eq(nl.Password)).
					AnyTimes().Return(c, nil)
			},
		},
		// {

		// 	name: "Fail_NoEmail",
		// 	body: models.NewUser{
		// 		Name:     "testuser",
		// 		Password: "password",
		// 	},
		// 	expectedStatus:   http.StatusBadRequest,
		// 	expectedResponse: `{"msg":"please provide Name, Email and Password"}`,
		// 	mockUserService: func(m *mockmodels.MockService) {
		// 		m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
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
				A: nil,
			}
			ctx := context.Background()
			traceID := "fake-trace-id"
			ctx = context.WithValue(ctx, middlewares.TraceIdKey, traceID)

			//register endpoints
			router.POST("/userLogin", h.userLogin)
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/userLogin", bytes.NewReader(body))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)
			require.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
