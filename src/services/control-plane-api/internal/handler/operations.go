package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/service/servermanager"
)

func StartServer(svc servermanager.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		op, err := svc.StartServer(ctx, serverId)

		if err != nil {
			apiErr := ServerStateOperationError("start", err)
			SendServerError(w, r, apiErr, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(op)
	}
}

func StopServer(svc servermanager.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		op, err := svc.StopServer(ctx, serverId)

		if err != nil {
			apiErr := ServerStateOperationError("stop", err)
			SendServerError(w, r, apiErr, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(op)
	}
}

func RebootServer(svc servermanager.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		op, err := svc.RebootServer(ctx, serverId)

		if err != nil {
			apiErr := ServerStateOperationError("reboot", err)
			SendServerError(w, r, apiErr, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(op)
	}
}
