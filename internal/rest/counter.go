package rest

import (
	"net/http"
)

// CounterService
type CounterService interface {
	Create() error
	Find() (string, error)
}

// CounterHandler
type CounterHandler struct {
	svc CounterService
}

// NewCounterHandler
func NewCounterHandler(svc CounterService) *CounterHandler {
	return &CounterHandler{
		svc: svc,
	}
}

// Register connects the handlers to the router
func (c *CounterHandler) Register(r *http.ServeMux) {
	r.HandleFunc("/count", c.find)
	r.HandleFunc("/", c.create)
}

func (c *CounterHandler) find(w http.ResponseWriter, r *http.Request) {
	val, err := c.svc.Find()
	if err != nil {
		renderErrorResponse(r.Context(), w, "create failed", err)
		return
	}
	renderResponse(w, val, http.StatusOK)
}

func (c *CounterHandler) create(w http.ResponseWriter, r *http.Request) {
	err := c.svc.Create()
	if err != nil {
		renderErrorResponse(r.Context(), w, "create failed", err)
		return
	}
	renderResponse(w, struct{}{}, http.StatusOK)
}
