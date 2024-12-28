package handler

import (
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	tempUserId := r.Context().Value("user-id")
	if tempUserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(tempUserId.(string)) //TODO: add middleware
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var newInfo *entity.User

	err = json.NewDecoder(r.Body).Decode(&newInfo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newInfo.Id = userId

	err = h.services.Profile.Update(r.Context(), newInfo) //TODO: think about ctx
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
