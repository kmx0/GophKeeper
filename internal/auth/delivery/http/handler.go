package http

import (
	"errors"
	"fmt"
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
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)
	if err := c.BindJSON(inp); err != nil {
		fmt.Printf("error on SignUp: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := h.useCase.SignUp(c.Request.Context(), inp.Login, inp.Password)
	if err != nil && !errors.Is(errors.Unwrap(err), auth.ErrLoginBusy) {
		fmt.Printf("error on SignUp--: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if errors.Is(errors.Unwrap(err), auth.ErrLoginBusy) {
		fmt.Printf("error on SignUp: %s", err.Error())
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	c.Status(http.StatusOK)

}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		fmt.Printf("error on SignIn: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.Login, inp.Password)
	if err != nil && !errors.Is(err, auth.ErrUserNotFound) {
		fmt.Printf("error on SignIn: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if errors.Is(errors.Unwrap(err), auth.ErrUserNotFound) {
		fmt.Printf("error on SignIn: %s", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}
