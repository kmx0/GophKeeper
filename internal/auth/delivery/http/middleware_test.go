package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/auth/usecase"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	uc := usecase.NewMockUseCase(ctrl)

	r := gin.Default()

	tests := []struct {
		name        string
		token       string
		err         error
		setHeader   bool
		headerKey   string
		headerValue string
		want        int
	}{
		{
			name:        "internal server error",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "Bearer " + "token",
			token:       "token",
			err:         fmt.Errorf("any error"),
			want:        http.StatusInternalServerError,
		},
		{
			name:        "incorrect token",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "Bearer ",
			err:         auth.ErrInvalidAccessToken,
			want:        http.StatusUnauthorized,
		},
		{
			name: "header authorization not found",
			want: http.StatusUnauthorized,
		},
		{
			name:        "unauthorized",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "",
			want:        http.StatusUnauthorized,
		},
		{
			name:        "fail bearer",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "FailBearer ",
			err:         nil,
			want:        http.StatusUnauthorized,
		},
		{
			name:        "empty token",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "Bearer",
			err:         nil,
			want:        http.StatusUnauthorized,
		},
		{
			name:        "correct token",
			setHeader:   true,
			headerKey:   "Authorization",
			headerValue: "Bearer " + "token",
			token:       "token",
			want:        http.StatusOK,
		},

	}
	var calls []*gomock.Call
	for _, tt := range tests {

		call := uc.EXPECT().ParseToken(ctx, tt.token).Return(&models.User{}, tt.err).AnyTimes()

		calls = append(calls, call)
	}
	gomock.InOrder(calls...)
	r.POST("/api/endpoint", NewAuthMiddleware(uc), func(c *gin.Context) {
		c.Status(http.StatusOK)

	})
	for _, tt := range tests {
		req, _ := http.NewRequest(http.MethodPost, "/api/endpoint", nil)
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if tt.setHeader {
				req.Header.Set(tt.headerKey, tt.headerValue)
			}
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		})
	}

	// w = httptest.NewRecorder()
	// req.Header.Set("Authorization", "")
	// r.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusUnauthorized, w.Code)

	// w = httptest.NewRecorder()
	// uc.EXPECT().ParseToken(ctx, "").Return(&models.User{}, auth.ErrInvalidAccessToken)

	// // uc.On("ParseToken", "").Return(&models.User{}, auth.ErrInvalidAccessToken)
	// req.Header.Set("Authorization", "Bearer ")
	// r.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusUnauthorized, w.Code)

	// w = httptest.NewRecorder()
	// uc.EXPECT().ParseToken(ctx, "token").Return(&models.User{}, nil).AnyTimes()

	// req.Header.Set("Authorization", "Bearer token")

	// r.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusOK, w.Code)

}
