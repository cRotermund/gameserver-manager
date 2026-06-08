package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/service/servermanager"
)

func ListServers(svc servermanager.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		serversArr, err := svc.ListServers(ctx)

		if err != nil {
			apiErr := ListServersApiError(err)
			SendServerError(w, r, apiErr, err)
			return
		}

		json.NewEncoder(w).Encode(serversArr)
	}
}

func ServerDetails(svc servermanager.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		detail, err := svc.GetServer(ctx, serverId)

		if err != nil {
			apiErr := ServerDetailError(err)
			SendServerError(w, r, apiErr, err)
			return
		}

		json.NewEncoder(w).Encode(detail)
	}
}
