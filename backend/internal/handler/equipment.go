package handler

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/filter"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) getAllEquipments(w http.ResponseWriter, r *http.Request) {

	tempQuery := r.URL.Query()
	id := tempQuery.Get("id")
	if id != "" && len(tempQuery) == 1 {
		h.getItem(w, r)
		return
	}

	tempOptions := r.Context().Value(filter.OptionsContextKey).(filter.OptionsMap)

	locations := r.URL.Query()["location"]
	for _, location := range locations {
		if location != "" {
			err := tempOptions.AddField("locationId", filter.OperatorEq, location, filter.DataTypeInt)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	eTypes := r.URL.Query()["type"]
	for _, eType := range eTypes {
		if eType != "" {
			err := tempOptions.AddField("type", filter.OperatorEq, eType, filter.DataTypeStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	brands := r.URL.Query()["brand"]
	for _, brand := range brands {
		if brand != "" {
			err := tempOptions.AddField("brand", filter.OperatorEq, brand, filter.DataTypeStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	models := r.URL.Query()["model"]

	for _, model := range models {
		if model != "" {
			err := tempOptions.AddField("model", filter.OperatorEq, model, filter.DataTypeStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	years := r.URL.Query()["year"]
	for _, year := range years {
		if year != "" {
			err := tempOptions.AddField("year", filter.OperatorEq, year, filter.DataTypeInt)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	prices := r.URL.Query()["price"]
	for _, price := range prices {
		if price != "" {
			operator := filter.OperatorEq
			if strings.Index(price, ":") != -1 {
				operator = filter.OperatorBetween
			}
			err := tempOptions.AddField("RentalPricePerDay", operator, price, filter.DataTypeFloat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	all, err := h.services.Equipment.GetAll(r.Context(), tempOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Info(err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var item entity.Equipment
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err := h.services.Equipment.Create(r.Context(), &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {

	var newInfo entity.Equipment
	if err := json.NewDecoder(r.Body).Decode(&newInfo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := h.services.Equipment.Update(r.Context(), &newInfo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) //TODO add middleware
		return
	}
	itemId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = h.services.Equipment.Delete(r.Context(), itemId); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) {

	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	itemId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := h.services.Equipment.Get(r.Context(), itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAvailableProducts(w http.ResponseWriter, r *http.Request) {

	tempId := r.URL.Query().Get("id")
	if tempId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	itemId, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	availableProducts, err := h.services.Equipment.GetAvailable(r.Context(), itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(availableProducts)
}

func (h *Handler) addProduct(w http.ResponseWriter, r *http.Request) {
	var availableProduct entity.CreateAvailableProductsInput
	if err := json.NewDecoder(r.Body).Decode(&availableProduct); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		slog.Info(err.Error())
		return
	}
	err := h.services.Equipment.AddProductToLocation(r.Context(), availableProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.services.Equipment.GetLogs(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(logs)
}
