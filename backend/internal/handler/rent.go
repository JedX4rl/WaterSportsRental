package handler

import (
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) rent(w http.ResponseWriter, r *http.Request) {

	tempUserId := r.Context().Value("user-id") //useless
	if tempUserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(tempUserId.(string))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var equipmentReq entity.EquipmentRequest

	err = json.NewDecoder(r.Body).Decode(&equipmentReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if equipmentReq.EndDate.Compare(equipmentReq.StartDate) < 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	agreement, err := h.services.Agreement.Create(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.services.Equipment.Rent(r.Context(), agreement.Id, equipmentReq) //TODO: add amount of available
	if err != nil {
		_ = h.services.Agreement.Delete(r.Context(), agreement.Id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("da")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getRentedItems(w http.ResponseWriter, r *http.Request) {

	tempUserId := r.Context().Value("user-id")
	if tempUserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(tempUserId.(string))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := h.services.Equipment.GetRentedItems(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(items)
}

func (h *Handler) getAvailableDates(w http.ResponseWriter, r *http.Request) {
	tempUserId := r.Context().Value("user-id")
	if tempUserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusInternalServerError)
		return
	}

	var datesRequest entity.AvailableDatesRequest

	tempLocationId := r.URL.Query().Get("location_id")
	if tempLocationId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	locationId, err := strconv.Atoi(tempLocationId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tempItemId := r.URL.Query().Get("item_id")
	if tempItemId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	itemId, err := strconv.Atoi(tempItemId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	datesRequest.LocationId = locationId
	datesRequest.ItemId = itemId

	availableDates, err := h.services.Equipment.GetAvailableDates(r.Context(), datesRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(availableDates)
}
