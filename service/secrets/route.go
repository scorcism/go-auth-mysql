package secrets

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/scorcism/go-auth/service/auth"
	"github.com/scorcism/go-auth/types"
	"github.com/scorcism/go-auth/utils"
)

type Handler struct {
	store     types.SecretsStore
	userStore types.UserStore
}

func NewHandler(store types.SecretsStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get-secrets", h.handleGetSecrets).Methods(http.MethodGet)
	router.HandleFunc("/add-secret", auth.WithJWTAuth(h.handleAddNewSecret, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleGetSecrets(w http.ResponseWriter, r *http.Request) {
	secrets, err := h.store.GetSecrets()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, secrets)
}

func (h *Handler) handleAddNewSecret(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())

	if !ok {
		fmt.Println("User not found in context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error while context"))
		return
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	userId := 1 

	fmt.Printf("userId: %v", userId)

	// get the payload
	var payload types.AddSecretPayload

	// parse the payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := h.store.AddSecret(types.Secret{
		SecretKey: payload.SecretKey,
		Label:     payload.Label,
		UserId:    userId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "New secret add success",
		"success": true,
	})
}
