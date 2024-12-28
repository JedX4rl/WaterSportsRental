package handler

import (
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

func (h *Handler) CreateLocation(w http.ResponseWriter, r *http.Request) {

	var location entity.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(location); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.services.Location.Create(r.Context(), &location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAllLocations(w http.ResponseWriter, r *http.Request) {

	locations, err := h.services.Location.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateLocationInfo(w http.ResponseWriter, r *http.Request) {
	var location entity.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} //TODO add checks?
	err := h.services.Location.Update(r.Context(), &location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteLocation(w http.ResponseWriter, r *http.Request) {

	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	locationId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = h.services.Location.Delete(r.Context(), locationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteAllLocations(w http.ResponseWriter, r *http.Request) {
	err := h.services.Location.DeleteAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
