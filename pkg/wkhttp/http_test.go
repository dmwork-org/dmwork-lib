package wkhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponseErrorWithStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		status         int
		expectedStatus int
	}{
		{
			name:           "returns 404 when status is 404",
			err:            errors.New("not found"),
			status:         http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 500 when status is 500",
			err:            errors.New("internal error"),
			status:         http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "returns 401 when status is 401",
			err:            errors.New("unauthorized"),
			status:         http.StatusUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "returns 403 when status is 403",
			err:            errors.New("forbidden"),
			status:         http.StatusForbidden,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ctx := &Context{Context: c}
			ctx.ResponseErrorWithStatus(tt.err, tt.status)

			if w.Code != tt.expectedStatus {
				t.Errorf("ResponseErrorWithStatus() HTTP status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseError(errors.New("test error"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("ResponseError() HTTP status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestResponseOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseOK()

	if w.Code != http.StatusOK {
		t.Errorf("ResponseOK() HTTP status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestResponseWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseWithStatus(http.StatusCreated, gin.H{"id": 1})

	if w.Code != http.StatusCreated {
		t.Errorf("ResponseWithStatus() HTTP status = %d, want %d", w.Code, http.StatusCreated)
	}
}
