package handler

import (
	"WaterSportsRental/internal/entity"
	Errors "WaterSportsRental/internal/errors"
	"WaterSportsRental/internal/parser"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var input entity.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, Errors.JsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if err = parser.IsUserValid(input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, Errors.JsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	input.Password = string(encryptedPassword)

	user := entity.User{
		Email:            input.Email,
		Password:         input.Password,
		RegistrationDate: time.Now(),
		Role:             "user",
	}

	err = h.services.Authorization.Create(r.Context(), &user)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr) && postgresErr.Code.Name() == "unique_violation" {
			http.Error(w, Errors.JsonError("User with the given email already exists"), http.StatusConflict)
			return
		}
		http.Error(w, Errors.JsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User created successfully"}`))
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input entity.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.services.Authorization.GetByEmail(r.Context(), input.Email)
	if err != nil {
		http.Error(w, Errors.JsonError("User does not exist"), http.StatusNotFound)
		return
	} //TODO: think about validation
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, Errors.JsonError("Invalid email or password"), http.StatusUnauthorized)
		return
	}

	token := os.Getenv("ACCESS_TOKEN_SECRET")
	if token == "" {
		log.Fatalf("ACCESS_TOKEN_SECRET environment variable not set")
	}
	expiry, _ := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"), 10, 64) //TODO: move
	if expiry == 0 {
		log.Fatalf("ACCESS_TOKEN_EXPIRY_HOUR environment variable not set")
	}

	accessToken, err := h.services.CreateAccessToken(&user, token, int(expiry))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		return
	}
}

func (h *Handler) hasPermission(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
