package handler

import (
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) postReview(w http.ResponseWriter, r *http.Request) {

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

	var review entity.Review

	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	review.UserId = userId

	user, err := h.services.Profile.Get(r.Context(), userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	review.Name = user.FirstName
	review.ReviewDate = time.Now()

	err = h.services.Review.Create(r.Context(), &review)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllReviews(w http.ResponseWriter, r *http.Request) {

	reviews, err := h.services.Review.GetAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) getReviewsByItemId(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	reviewId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	reviews, err := h.services.Review.GetByItemID(r.Context(), reviewId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)

}

func (h *Handler) getReviewsByUserId(w http.ResponseWriter, r *http.Request) {

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

	reviews, err := h.services.Review.GetByUserID(r.Context(), userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)

}

func (h *Handler) DeleteReviewById(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	reviewId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = h.services.Review.DeleteByID(r.Context(), reviewId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
