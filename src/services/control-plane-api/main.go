package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/apimodels"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/awsutils"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/handling"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/logging"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/mapping"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const ServerStartError string = "Could not start server."
const EnvPort string = "PORT"
const EnvGsmIsServerTagKey string = "GAMESERVER_ISSERVER_TAG_KEY"
const EnvGsmServerNameTagKey string = "GAMESERVER_NAME_TAG_KEY"
const EnvGsmGameDescriptorTagKey string = "GAMESERVER_DESCRIPTOR_TAG_KEY"

func main() {
	//Initialize the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//Load .env and various configuration
	initEnv()
	port := loadRequiredEnvAsInt(EnvPort)
	gsmIsServerTagKey := loadRequiredEnvAsString(EnvGsmIsServerTagKey)
	gsmServerNameTagKey := loadRequiredEnvAsString(EnvGsmServerNameTagKey)
	gsmGameDescriptorTagKey := loadRequiredEnvAsString(EnvGsmGameDescriptorTagKey)
	awsConfig := loadAwsConfig()

	//Init the server
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	mux.HandleFunc("GET /v1/servers", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		awsclient := awsutils.NewClient(awsConfig)
		instanceArr, err := awsclient.ListServerInstances(ctx, gsmIsServerTagKey, "True")

		if err != nil {
			apiErr := handling.ListServersApiError(err)
			writeError(w, r, apiErr, err)
			return
		}

		var summaryArr []apimodels.ServerSummary
		for _, instance := range instanceArr {
			summaryArr = append(
				summaryArr,
				mapping.InstanceFlyweightToSummary(
					instance,
					gsmServerNameTagKey,
					gsmGameDescriptorTagKey))
		}

		json.NewEncoder(w).Encode(summaryArr)
	})

	mux.HandleFunc("GET /v1/servers/{serverId}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		awsclient := awsutils.NewClient(awsConfig)
		details, err := awsclient.GetServerInstanceDetails(ctx, serverId)

		if err != nil {
			apiErr := handling.ServerDetailError(err)
			writeError(w, r, apiErr, err)
			return
		}

		detailModel := mapping.InstanceDetailsToDetail(
			*details,
			gsmServerNameTagKey,
			gsmGameDescriptorTagKey)

		json.NewEncoder(w).Encode(detailModel)
	})

	mux.HandleFunc("POST /v1/servers/{serverId}/start", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		awsclient := awsutils.NewClient(awsConfig)
		err := awsclient.StartServerInstance(ctx, serverId)

		if err != nil {
			apiErr := handling.ServerStateOperationError("start", err)
			writeError(w, r, apiErr, err)
			return
		}

		operation := apimodels.Operation{
			OperationID: uuid.New().String(),
			Type:        "start",
			ServerID:    serverId,
			Status:      "issued",
			CreatedAt:   time.Now().UTC(),
			CompletedAt: nil,
			Error:       nil,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(operation)
	})

	mux.HandleFunc("POST /v1/servers/{serverId}/stop", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		awsclient := awsutils.NewClient(awsConfig)
		err := awsclient.StopServerInstance(ctx, serverId)

		if err != nil {
			apiErr := handling.ServerStateOperationError("stop", err)
			writeError(w, r, apiErr, err)
			return
		}

		operation := apimodels.Operation{
			OperationID: uuid.New().String(),
			Type:        "stop",
			ServerID:    serverId,
			Status:      "issued",
			CreatedAt:   time.Now().UTC(),
			CompletedAt: nil,
			Error:       nil,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(operation)
	})

	mux.HandleFunc("POST /v1/servers/{serverId}/reboot", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverId := r.PathValue("serverId")

		awsclient := awsutils.NewClient(awsConfig)
		err := awsclient.RebootServerInstance(ctx, serverId)

		if err != nil {
			apiErr := handling.ServerStateOperationError("reboot", err)
			writeError(w, r, apiErr, err)
			return
		}

		operation := apimodels.Operation{
			OperationID: uuid.New().String(),
			Type:        "reboot",
			ServerID:    serverId,
			Status:      "issued",
			CreatedAt:   time.Now().UTC(),
			CompletedAt: nil,
			Error:       nil,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(operation)
	})

	//Start the server and listen...

	logging.ServerStarting(port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		logging.ServerListenLoopError(err)
		log.Fatalf("server error: %v", err)
	}
}

func initEnv() {
	err := godotenv.Load()

	if err != nil {
		logging.FailedEnvInit(err)
		log.Fatal(ServerStartError)
	}
}

func loadAwsConfig() aws.Config {
	awscfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigFiles([]string{"./aws-config.toml"}),
		config.WithSharedCredentialsFiles([]string{"./aws-credentials.toml"}),
		config.WithSharedConfigProfile("control-plane-app"),
	)

	if err != nil {
		logging.MissingAwsConfig(err)
		log.Fatal(ServerStartError)
	}

	return awscfg
}

func loadRequiredEnvAsString(v string) string {
	output := os.Getenv(v)

	if output == "" {
		logging.MissingEnv(v)
		log.Fatal(ServerStartError)
	}

	return output
}

func loadRequiredEnvAsInt(v string) int {
	output := os.Getenv(v)

	if output == "" {
		logging.MissingEnv(v)
		log.Fatal(ServerStartError)
	}

	portNum, err := strconv.Atoi(output)

	if err != nil {
		logging.MalformedIntEnv(v, output)
		log.Fatal(ServerStartError)
	}

	return portNum
}

func writeError(w http.ResponseWriter, r *http.Request, aErr apimodels.APIError, oErr error) {
	//log
	logging.RequestHandlingError(r, oErr)

	//write body and fail
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(aErr)
}
