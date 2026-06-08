package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/config"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/handler"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/logging"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/service/servermanager"
	"github.com/go-chi/chi/v5"
)

const ServerStartError string = "Could not start server."

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(ServerStartError)
	}

	ec2 := ec2.NewFromConfig(cfg.AWS)
	svc := servermanager.New(ec2, cfg)

	r := chi.NewRouter()

	r.Get("/health", handler.Health())

	r.Route("/v1", func(r chi.Router) {
		r.Get("/servers", handler.ListServers(svc))
		r.Get("/servers/{serverId}", handler.ServerDetails(svc))
		r.Post("/servers/{serverId}/start", handler.StartServer(svc))
		r.Post("/servers/{serverId}/stop", handler.StopServer(svc))
		r.Post("/servers/{serverId}/reboot", handler.RebootServer(svc))
	})

	//Start the server and listen...

	logging.ServerStarting(cfg.Port)

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), r); err != nil {
		logging.ServerListenLoopError(err)
		log.Fatalf("server error: %v", err)
	}
}
