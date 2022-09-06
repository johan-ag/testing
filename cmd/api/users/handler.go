package users

import (
	"net/http"

	"github.com/johan-ag/testing/internal/users"
	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type handler struct {
	service users.Service
}

func NewHandler(service users.Service) *handler {
	return &handler{
		service,
	}
}

func (h *handler) Save(w http.ResponseWriter, r *http.Request) error {
	var user users.User
	err := web.DecodeJSON(r, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return web.NewError(http.StatusBadGateway, "error to read body")
		// return web.EncodeJSON(w, err, http.StatusBadRequest)
	}

	id, err := h.service.Save(r.Context(), user.Name, user.Age)
	if err != nil {
		err = web.NewError(http.StatusInternalServerError, "error to save user")
		return web.EncodeJSON(w, err, http.StatusInternalServerError)
	}

	return web.EncodeJSON(w, id, http.StatusCreated)
}

func (h *handler) Find(w http.ResponseWriter, r *http.Request) error {
	id, err := web.Params(r).Uint("id")
	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.Find(r.Context(), id)
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err.Error())
	}

	return web.EncodeJSON(w, user, http.StatusCreated)
}
