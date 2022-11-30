package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/auth"
)

type AuthMiddleware struct {
	usecase auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	switch {

	case err == nil:
		c.Set(auth.CtxUserKey, user)
	case errors.Is(errors.Unwrap(err), auth.ErrInvalidAccessToken):
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

}
