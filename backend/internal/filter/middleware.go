package filter

import (
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

const (
	OptionsContextKey = "filter_options"
)

func Middleware(h http.HandlerFunc, defaultLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitFromQuery := r.URL.Query().Get("limit")

		limit := defaultLimit

		var limitParseErr error

		if limitFromQuery != "" {
			if limit, limitParseErr = strconv.Atoi(limitFromQuery); limitParseErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "limit must be an integer"}`))
				return
			}
		}

		_ = limit //TODO: remove

		optionsI := NewOptionsMap()
		ctx := context.WithValue(r.Context(), OptionsContextKey, optionsI)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
