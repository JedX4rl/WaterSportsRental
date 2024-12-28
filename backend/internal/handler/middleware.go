package handler

import (
	"context"
	"net/http"
	"strconv"
)

func AdminMiddleware(h *Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tempId := r.Context().Value("user-id")
			if tempId == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			userId, err := strconv.Atoi(tempId.(string))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, err := h.services.Authorization.GetById(r.Context(), userId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if user.Role != "admin" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "admin", true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		})
	}
}
