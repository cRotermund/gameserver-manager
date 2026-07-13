package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/errors"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/middleware"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
)

func SendServerError(w http.ResponseWriter, r *http.Request, e error) {
	logger := middleware.LoggerFromContext(r.Context())

	logger.Error("Error handling http request", "error", e)

	apiErr := errors.From(e)

	errResp := models.APIError{
		Code:  apiErr.Code,
		Error: apiErr.Message,
	}

	w.WriteHeader(apiErr.HTTPStatus)
	json.NewEncoder(w).Encode(errResp)
}
