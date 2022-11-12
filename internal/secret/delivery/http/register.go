package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/secret"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc secret.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/secret")
	{
		authEndpoints.POST("/create",  h.Create)
		authEndpoints.POST("/get",     h.Get)
		authEndpoints.POST("/get-all", h.GetAll)
		authEndpoints.POST("/delete",  h.Delete)
	}
}
