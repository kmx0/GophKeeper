package user

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kmx0/GophKeeper/internal/adapters/api"
)

const (
	userURL  = "/users/:sercet_id"
	usersURL = "/users"
)

type handler struct {
	userService Service
}

func NewHandler(service Service) api.Handler {
	return &handler{userService: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetAllUsers)
}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_ = h.userService.GetAll(context.Background(), time.Hour)
	// json.Unmarshal
	w.Write([]byte("users"))
	w.WriteHeader(http.StatusOK)
}
