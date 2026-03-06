package wkhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseErrorWithStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		status         int
		err            error
		expectedStatus int
	}{
		{
			name:           "returns custom status 500",
			status:         http.StatusInternalServerError,
			err:            errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "returns custom status 404",
			status:         http.StatusNotFound,
			err:            errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns custom status 403",
			status:         http.StatusForbidden,
			err:            errors.New("forbidden"),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns custom status 422",
			status:         http.StatusUnprocessableEntity,
			err:            errors.New("validation error"),
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ctx := &Context{Context: c}
			ctx.ResponseErrorWithStatus(tt.err, tt.status)

			if w.Code != tt.expectedStatus {
				t.Errorf("ResponseErrorWithStatus() returned status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseError(errors.New("bad request"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("ResponseError() returned status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
