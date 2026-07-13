package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/joho/godotenv"
)

const EnvPort string = "PORT"
const EnvGsmIsServerTagKey string = "GAMESERVER_ISSERVER_TAG_KEY"
const EnvGsmServerNameTagKey string = "GAMESERVER_NAME_TAG_KEY"
const EnvGsmGameDescriptorTagKey string = "GAMESERVER_DESCRIPTOR_TAG_KEY"
const EnvLogLevel string = "LOG_LEVEL"

type ControlPlaneConfig struct {
	Port                     int
	IsServerAwsTagName       string
	ServerNameAwsTagName     string
	GameDescriptorAwsTagName string
	IsServerAwsTagValue      string
	AWS                      aws.Config
	LogLevel                 slog.Level
}

func Load(logger *slog.Logger) (*ControlPlaneConfig, error) {
	var cfg ControlPlaneConfig

	err := initEnv(logger)
	if err != nil {
		return nil, err
	}

	p, err := loadRequiredEnvAsInt(EnvPort, logger)
	if err != nil {
		return nil, err
	}
	cfg.Port = *p

	isServerTag, err := loadRequiredEnvAsString(EnvGsmIsServerTagKey, logger)
	if err != nil {
		return nil, err
	}
	cfg.IsServerAwsTagName = *isServerTag

	serverNameTag, err := loadRequiredEnvAsString(EnvGsmServerNameTagKey, logger)
	if err != nil {
		return nil, err
	}
	cfg.ServerNameAwsTagName = *serverNameTag

	descriptorTagName, err := loadRequiredEnvAsString(EnvGsmGameDescriptorTagKey, logger)
	if err != nil {
		return nil, err
	}
	cfg.GameDescriptorAwsTagName = *descriptorTagName

	awsCfg, err := loadAwsConfig(logger)
	if err != nil {
		return nil, err
	}
	cfg.AWS = *awsCfg

	cfg.IsServerAwsTagValue = "True" //TODO roll into config, for real. Don't hardcode

	return &cfg, nil
}

func LogLevelFromEnv() slog.Level {
	switch strings.ToUpper(os.Getenv(EnvLogLevel)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initEnv(logger *slog.Logger) error {
	err := godotenv.Load()

	if err != nil {
		logger.Warn("No .env file found, using existing environment", "Error", err)
	}

	return nil
}

func loadRequiredEnvAsString(v string, logger *slog.Logger) (*string, error) {
	output := os.Getenv(v)

	if output == "" {
		logger.Error("Required environment variable is missing", "VariableName", v)
		return nil, errors.New("Required environment variable is not populated")
	}

	return &output, nil
}

func loadRequiredEnvAsInt(v string, logger *slog.Logger) (*int, error) {
	output := os.Getenv(v)

	if output == "" {
		logger.Error("Required environment variable is missing", "VariableName", v)
		return nil, errors.New("Required environment variable is not populated")
	}

	portNum, err := strconv.Atoi(output)

	if err != nil {
		logger.Error(
			"Could not parse required environment variable as integer",
			"VariableName", v,
			"VariableValue", output)
		return nil, err
	}

	return &portNum, nil
}

func loadAwsConfig(logger *slog.Logger) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Error("Unable to load AWS configuration", "Error", err)
		return nil, err
	}

	if roleARN := os.Getenv("AWS_ROLE_ARN"); roleARN != "" {
		stsClient := sts.NewFromConfig(cfg)
		cfg.Credentials = aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(stsClient, roleARN),
		)
	}

	return &cfg, nil
}
