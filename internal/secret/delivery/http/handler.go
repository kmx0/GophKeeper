package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

type Secret struct {
	//int
	ID     string
	UserID string
	Value  string
}

type Handler struct {
	useCase secret.UseCase
}

func NewHandler(useCase secret.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Value string `json:"value"`
}

func (h *Handler) Create(c *gin.Context) {
	input := new(createInput)
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.CreateSecret(c, user, input.Value); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

type getResponseSingle struct {
	Secret *Secret `json:"secret"`
}

type getInput struct {
	ID string `json:"id"`
}

func (h *Handler) Get(c *gin.Context) {
	input := new(getInput)
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)
	sc, err := h.useCase.GetSecret(c, user, input.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &getResponseSingle{
		Secret: toSecret(sc),
	})
}

type getResponse struct {
	Secrets []*Secret `json:"secrets"`
}

func (h *Handler) List(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)
	scs, err := h.useCase.GetSecrets(c, user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &getResponse{
		Secrets: toSecrets(scs),
	})

}

type deleteInput struct {
	ID string `json:"id"`
}

func (h *Handler) Delete(c *gin.Context) {
	input := new(deleteInput)
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := c.MustGet(auth.CtxUserKey).(*models.User)
	err := h.useCase.DeleteSecret(c, user, input.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func toSecrets(s []*models.Secret) []*Secret {
	scs := make([]*Secret, len(s))
	for i, v := range s {
		scs[i] = toSecret(v)
	}
	return scs
}

// ?
// преобразуем к типу *Secret
func toSecret(s *models.Secret) *Secret {
	return &Secret{
		ID:     s.ID,
		UserID: s.UserID,
		Value:  s.Value,
	}
}
