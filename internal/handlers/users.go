package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

type CreateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// Password  string    `json:"password"`
}

func (cfg *ApiConfig) HandlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	timeNow := time.Now()
	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	newUser := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		Email:          req.Email,
		HashedPassword: hashedPwd,
	}

	createdUser, err := cfg.Db.CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	res, err := json.Marshal(user)
	if err != nil {
		fmt.Println("couldn't Marshal the user: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (cfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := cfg.Db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	isPwdValid, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil || !isPwdValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	userResponse := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}
