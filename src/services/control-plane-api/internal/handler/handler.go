package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/middleware"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
)

func SendServerError(w http.ResponseWriter, r *http.Request, aErr models.APIError, oErr error) {
	logger := middleware.LoggerFromContext(r.Context())

	logger.Error("Error handling http request", "error", oErr)

	//write body and fail
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(aErr)
}
