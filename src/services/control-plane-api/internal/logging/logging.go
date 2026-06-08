package logging

import (
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func FailedEnvInit(err error) {
	slog.Error("Error loading .env file", "Error", err)
}

func MissingAwsConfig(err error) {
	slog.Error("Unable to load AWS configuration files", "Error", err)
}

func MissingEnv(varName string) {
	slog.Error("Required environment variable is missing", "VariableName", varName)
}

func MalformedIntEnv(varName string, varValue string) {
	slog.Error(
		"Could not parse required environment variable as integer",
		"VariableName", varName,
		"VariableValue", varValue)
}

func RequestHandlingError(r *http.Request, err error) {
	slog.Error(
		"Error handling http request",
		"Method", r.Method,
		"PathActual", r.URL.EscapedPath(),
		"Route", r.Pattern,
		"Body", r.Body,
		"Error", err)
}

func ServerStarting(port int) {
	slog.Info("control-plane-api starting", "Port", port)
}

func ServerListenLoopError(err error) {
	slog.Info("Error listening and serving http", "Error", err)
}

func DescribeInstanceError(input *ec2.DescribeInstancesInput, err error) {
	slog.Info(
		"Error describing instances",
		"Error", err,
		"Input", input)
}

func GetInstanceDetailsError(instanceId string, err error) {
	slog.Error("Error getting server instance details.", "Error", err)
}

func StartInstanceActionError(instanceId string, input *ec2.StartInstancesInput, err error) {
	slog.Info(
		"Error starting instance",
		"Error", err,
		"Input", input)
}

func StopInstanceActionError(instanceId string, input *ec2.StopInstancesInput, err error) {
	slog.Info(
		"Error stopping instance",
		"Error", err,
		"Input", input)
}

func RebootInstanceActionError(instanceId string, input *ec2.RebootInstancesInput, err error) {
	slog.Info(
		"Error rebooting instance",
		"Error", err,
		"Input", input)
}
