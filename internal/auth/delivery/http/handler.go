package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/auth"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type signInput struct {
	Login    string
	Password string
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := h.useCase.SignUp(c.Request.Context(), inp.Login, inp.Password); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)

}
