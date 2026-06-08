package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/logging"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
)

func SendServerError(w http.ResponseWriter, r *http.Request, aErr models.APIError, oErr error) {
	//log
	logging.RequestHandlingError(r, oErr)

	//write body and fail
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(aErr)
}
