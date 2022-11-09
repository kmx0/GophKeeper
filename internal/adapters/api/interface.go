package adapters

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Register(register *httprouter.Router)
}
