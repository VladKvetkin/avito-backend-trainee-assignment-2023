package handler

import (
	"encoding/json"
	"net/http"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/config"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/models"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/storage"
)

type Handler struct {
	storage storage.Storage
	config  config.Config
}

func NewHandler(storage storage.Storage, config config.Config) *Handler {
	return &Handler{
		storage: storage,
		config:  config,
	}
}

func (h *Handler) errorJSONResponse(res http.ResponseWriter, errorMessage string, status int) {
	jsonError, err := json.Marshal(models.ErrorResponse{ErrorMessage: errorMessage})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(status)
	res.Write(jsonError)
}
