package config

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/logging"
	"github.com/joho/godotenv"
)

const EnvPort string = "PORT"
const EnvGsmIsServerTagKey string = "GAMESERVER_ISSERVER_TAG_KEY"
const EnvGsmServerNameTagKey string = "GAMESERVER_NAME_TAG_KEY"
const EnvGsmGameDescriptorTagKey string = "GAMESERVER_DESCRIPTOR_TAG_KEY"

type ControlPlaneConfig struct {
	Port                     int
	IsServerAwsTagName       string
	ServerNameAwsTagName     string
	GameDescriptorAwsTagName string
	IsServerAwsTagValue      string
	AWS                      aws.Config
}

func Load() (*ControlPlaneConfig, error) {
	var cfg ControlPlaneConfig

	err := initEnv()
	if err != nil {
		return nil, err
	}

	p, err := loadRequiredEnvAsInt(EnvPort)
	if err != nil {
		return nil, err
	}
	cfg.Port = *p

	isServerTag, err := loadRequiredEnvAsString(EnvGsmIsServerTagKey)
	if err != nil {
		return nil, err
	}
	cfg.IsServerAwsTagName = *isServerTag

	serverNameTag, err := loadRequiredEnvAsString(EnvGsmServerNameTagKey)
	if err != nil {
		return nil, err
	}
	cfg.ServerNameAwsTagName = *serverNameTag

	descriptorTagName, err := loadRequiredEnvAsString(EnvGsmGameDescriptorTagKey)
	if err != nil {
		return nil, err
	}
	cfg.GameDescriptorAwsTagName = *descriptorTagName

	awsCfg, err := loadAwsConfig()
	if err != nil {
		return nil, err
	}
	cfg.AWS = *awsCfg

	cfg.IsServerAwsTagValue = "True" //TODO roll into config, for real. Don't hardcode

	return &cfg, nil
}

func initEnv() error {
	err := godotenv.Load()

	if err != nil {
		logging.FailedEnvInit(err)
		return err
	}

	return nil
}

func loadRequiredEnvAsString(v string) (*string, error) {
	output := os.Getenv(v)

	if output == "" {
		logging.MissingEnv(v)
		return nil, errors.New("Required environment variable is not populated")
	}

	return &output, nil
}

func loadRequiredEnvAsInt(v string) (*int, error) {
	output := os.Getenv(v)

	if output == "" {
		logging.MissingEnv(v)
		return nil, errors.New("Required environment variable is not populated")
	}

	portNum, err := strconv.Atoi(output)

	if err != nil {
		logging.MalformedIntEnv(v, output)
		return nil, err
	}

	return &portNum, nil
}

func loadAwsConfig() (*aws.Config, error) {

	awscfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigFiles([]string{"./aws-config.toml"}),
		config.WithSharedCredentialsFiles([]string{"./aws-credentials.toml"}),
		config.WithSharedConfigProfile("control-plane-app"),
	)

	if err != nil {
		logging.MissingAwsConfig(err)
		return nil, err
	}

	return &awscfg, nil
}
