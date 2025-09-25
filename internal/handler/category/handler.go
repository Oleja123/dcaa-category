package categoryhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Oleja123/dcaa-category/internal/domain/category"
	categorydto "github.com/Oleja123/dcaa-property/pkg/dto/category"
	myErrors "github.com/Oleja123/dcaa-property/pkg/errors"
)

type CategoryHandler struct {
	service category.Service
}

func NewHandler(s category.Service) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.FindAll(rw, r)
	case http.MethodPost:
		h.Create(rw, r)
	case http.MethodPut:
		h.Update(rw, r)
	default:
		http.Error(rw, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleWithId(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.FindOne(rw, r)
	case http.MethodDelete:
		h.Delete(rw, r)
	default:
		http.Error(rw, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var dto categorydto.CategoryDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(rw, "неправильное тело запроса", http.StatusBadRequest)
		return
	}

	if !dto.Validate(false) {
		http.Error(rw, "неправильное тело запроса", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrInternalError):
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		case errors.Is(err, myErrors.ErrNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
		}
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]int{"id": id})
}

func (h *CategoryHandler) FindAll(rw http.ResponseWriter, r *http.Request) {
	list, err := h.service.FindAll(r.Context())

	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrInternalError):
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		case errors.Is(err, myErrors.ErrNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
		}
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(list)
}

func (h *CategoryHandler) FindOne(rw http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	dto, err := h.service.FindOne(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrInternalError):
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		case errors.Is(err, myErrors.ErrNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
		}
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(dto)
}

func (h *CategoryHandler) Update(rw http.ResponseWriter, r *http.Request) {
	var dto categorydto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(rw, "неправильное тело запроса", http.StatusBadRequest)
		return
	}

	if !dto.Validate(true) {
		http.Error(rw, "неправильное тело запроса", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), dto); err != nil {
		switch {
		case errors.Is(err, myErrors.ErrInternalError):
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		case errors.Is(err, myErrors.ErrNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
		}
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) Delete(rw http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, myErrors.ErrInternalError):
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		case errors.Is(err, myErrors.ErrNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
		}
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
