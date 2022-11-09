package secret

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	secretURL  = "/secrets/:sercet_id"
	secretsURL = "/secrets"
)

type handler struct {
	sercetService Service
}

func NewHandler(service Service) Handler {
	return &handler{sercetService: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(secretsURL, h.GetAllSecrets)
}

func (h *handler) GetAllSecrets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_ = h.sercetService.GetAllSecrets(context.Background(), time.Hour)
	// json.Unmarshal
	w.Write([]byte("secrets"))
	w.WriteHeader(http.StatusOK)
}
