package handler

import (
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.signUp(w, r)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := h.services.User.Get(r.Context(), userId)
	if err != nil {
		slog.Info(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.User.GetAll(r.Context())
	if err != nil {
		slog.Info(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateWithRole(w http.ResponseWriter, r *http.Request) {
	var newInfo *entity.User

	if err := json.NewDecoder(r.Body).Decode(&newInfo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := h.services.User.UpdateWithRole(r.Context(), newInfo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GiveRole(w http.ResponseWriter, r *http.Request) {
	var newInfo *entity.User

	if err := json.NewDecoder(r.Body).Decode(&newInfo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := h.services.User.GiveRole(r.Context(), newInfo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) //TODO add middleware
		return
	}
	userId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = h.services.User.Delete(r.Context(), userId); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserInfo(w http.ResponseWriter, r *http.Request) {
	tempUserId := r.Context().Value("user-id")
	if tempUserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(tempUserId.(string))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := h.services.User.Get(r.Context(), userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetAllRentalAgreementsByID(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	agreements, err := h.services.User.GetAllAgreementsById(r.Context(), userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(agreements)

}
