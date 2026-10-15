package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const GsmHostVarName = "GSM_API_HOST"
const GsmPortVarName = "GSM_API_PORT"
const GsmApiTimeoutVarName = "GSM_API_TIMEOUT_MS"
const DiscordTokenVarName = "DISCORD_TOKEN"

const DefaultGsmApiTimeoutMs int = 5000 //This is pretty darn, high, but the aws calls behind the scenes can be slow.

type DiscordBotConfig struct {
	logger       *slog.Logger
	GsmApiConfig GsmApiConfig
	DiscordToken string
}

type GsmApiConfig struct {
	Host      string
	Port      int
	TimeoutMs int
}

func NewBotConfig(logger *slog.Logger) *DiscordBotConfig {
	return &DiscordBotConfig{
		logger: logger,
	}
}

func (cfg *DiscordBotConfig) Load() error {
	err := cfg.initEnv()
	if err != nil {
		return err
	}

	cfg.GsmApiConfig.Host, err = cfg.loadRequiredEnvAsString(GsmHostVarName)
	if err != nil {
		return err
	}

	cfg.GsmApiConfig.Port, err = cfg.loadRequiredEnvAsInt(GsmPortVarName)
	if err != nil {
		return err
	}

	defaultTimeout := DefaultGsmApiTimeoutMs
	parsed, err := cfg.loadOptionalValueAsInt(GsmApiTimeoutVarName, &defaultTimeout)
	if err != nil {
		return err
	}
	cfg.GsmApiConfig.TimeoutMs = *parsed

	cfg.DiscordToken, err = cfg.loadRequiredEnvAsString(DiscordTokenVarName)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *DiscordBotConfig) initEnv() error {
	err := godotenv.Load()

	if err != nil {
		cfg.logger.Warn("No .env file found, using existing environment", "Error", err)
	}

	return nil
}

func (cfg *DiscordBotConfig) loadRequiredEnvAsString(v string) (string, error) {
	output := os.Getenv(v)

	if output == "" {
		cfg.logger.Error("Required environment variable is missing", "VariableName", v)
		return "", errors.New("Required environment variable is not populated")
	}

	return output, nil
}

func (cfg *DiscordBotConfig) loadRequiredEnvAsInt(v string) (int, error) {
	output := os.Getenv(v)

	if output == "" {
		cfg.logger.Error("Required environment variable is missing", "VariableName", v)
		return -1, errors.New("Required environment variable is not populated")
	}

	numValue, err := strconv.Atoi(output)

	if err != nil {
		cfg.logger.Error(
			"Could not parse required environment variable as integer",
			"VariableName", v,
			"VariableValue", output)
		return -1, err
	}

	return numValue, nil
}

func (cfg *DiscordBotConfig) loadOptionalValueAsInt(v string, defaultValue *int) (*int, error) {
	output := os.Getenv(v)

	if output == "" {
		if defaultValue != nil {
			return defaultValue, nil
		}
		return nil, nil
	}

	numValue, err := strconv.Atoi(output)

	if err != nil {
		cfg.logger.Error(
			"Could not parse environment variable as integer",
			"VariableName", v,
			"VariableValue", output)
		return nil, err
	}

	return &numValue, nil
}
