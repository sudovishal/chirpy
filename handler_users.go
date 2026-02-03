package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sudovishal/chirpy/internal/auth"
	"github.com/sudovishal/chirpy/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	// fmt.Println(params.Email)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.CreateUser(
		r.Context(), database.CreateUserParams{
			Email:          sql.NullString{String: params.Email, Valid: true},
			HashedPassword: hashedPwd,
		})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	apiUser := User{
		ID:             user.ID,
		CreatedAt:      user.UpdatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Email:          user.Email.String,
		HashedPassword: user.HashedPassword,
	}
	// fmt.Println(apiUser.ID)

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(apiUser)
}

func (cfg *apiConfig) deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(403)
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		w.WriteHeader(500)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "reset successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	// r.Header.Add("Content-Type", "application/json")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	// fmt.Println(params.Email)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	user, err := cfg.db.GetUserByEmail(
		r.Context(),
		sql.NullString{
			String: params.Email,
			Valid:  true,
		})

	apiUser := User{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Email:          user.Email.String,
		HashedPassword: user.HashedPassword,
	}

	fmt.Println(params.Password, apiUser.HashedPassword)
	passwordVerify, err := auth.CheckPasswordHash(params.Password, apiUser.HashedPassword)
	if err != nil {
		log.Printf("Error comparing passwords: %s", err)
		w.WriteHeader(401)
		return
	}

	if !passwordVerify {
		w.WriteHeader(401)
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiUser)
	}

}
