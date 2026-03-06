package wkhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
			name:           "returns 404 status",
			err:            errors.New("not found"),
			status:         http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 500 status",
			err:            errors.New("internal error"),
			status:         http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "returns 403 status",
			err:            errors.New("forbidden"),
			status:         http.StatusForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns 400 status",
			err:            errors.New("bad request"),
			status:         http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ctx := &Context{Context: c}
			ctx.ResponseErrorWithStatus(tt.err, tt.status)

			assert.Equal(t, tt.expectedStatus, w.Code, "HTTP status code should match the provided status")
		})
	}
}

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseError(errors.New("test error"))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResponseOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseOK()

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.Response(gin.H{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResponseWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := &Context{Context: c}
	ctx.ResponseWithStatus(http.StatusCreated, gin.H{"created": true})

	assert.Equal(t, http.StatusCreated, w.Code)
}
