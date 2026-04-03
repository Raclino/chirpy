package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

type CreateUserReq struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) HandlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("couldn't decode the body")
		return
	}

	timeNow := time.Now()
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Email:     req.Email,
	}

	createdUser, err := cfg.Db.CreateUser(r.Context(), newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
