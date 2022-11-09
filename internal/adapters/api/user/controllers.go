package adapters

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gophkee.per/internal/adapters"
	"gophkee.per/internal/domain/secret"
)

const (
	secretURL  = "/secrets/:sercet_id"
	secretsURL = "/secrets"
)

type handler struct {
	sercetService secret.Service
}

func NewHandler(service secret.Service) adapters.Handler {
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
