package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		handlerFn    func(c *gin.Context)
		expectedCode int
		expectedBody Response
	}{
		{
			name: "Success",
			handlerFn: func(c *gin.Context) {
				Success(c, "test data")
			},
			expectedCode: http.StatusOK,
			expectedBody: Response{Code: 0, Message: "success", Data: "test data"},
		},
		{
			name: "Error",
			handlerFn: func(c *gin.Context) {
				Error(c, http.StatusBadRequest, 4001, "bad input")
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: Response{Code: 4001, Message: "bad input", Data: nil},
		},
		{
			name: "InternalError",
			handlerFn: func(c *gin.Context) {
				InternalError(c, errors.New("system failed")) // H-03: raw err is discarded
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: Response{Code: 500, Message: "internal server error", Data: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.handlerFn(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var body Response
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}
