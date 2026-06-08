package servermanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/config"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
)

type Service interface {
	ListServers(ctx context.Context) ([]models.ServerSummary, error)
	GetServer(ctx context.Context, serverId string) (*models.ServerDetail, error)
	StartServer(ctx context.Context, serverId string) (*models.Operation, error)
	StopServer(ctx context.Context, serverId string) (*models.Operation, error)
	RebootServer(ctx context.Context, serverId string) (*models.Operation, error)
	GetOperation(ctx context.Context, operationId string) (*models.Operation, error)
}

type service struct {
	ec2 *ec2.Client
	cfg *config.ControlPlaneConfig
}

func New(ec2Client *ec2.Client, cfg *config.ControlPlaneConfig) Service {
	return &service{
		ec2: ec2Client,
		cfg: cfg,
	}
}
