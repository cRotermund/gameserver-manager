package servermanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/logging"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
	"github.com/google/uuid"
)

type InstanceState int

const (
	Pending InstanceState = iota
	Running
	Stopping
	Stopped
	ShuttingDown
	Terminated
)

func (s *service) ListServers(ctx context.Context) ([]models.ServerSummary, error) {
	// The filter "name" prefix controls the filter "subject" (e.g. tag:*)
	var tagFilterBuilder strings.Builder
	tagFilterBuilder.WriteString("tag:")
	tagFilterBuilder.WriteString(s.cfg.IsServerAwsTagName)

	//Build the describe query
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(tagFilterBuilder.String()),
				Values: []string{s.cfg.IsServerAwsTagValue},
			},
		},
	}

	instPaginator := ec2.NewDescribeInstancesPaginator(s.ec2, input)

	var instances []models.ServerSummary
	for instPaginator.HasMorePages() {

		page, err := instPaginator.NextPage(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to describe instances on page: %w", err)
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				summary := summaryFromInstance(instance, s.cfg.ServerNameAwsTagName, s.cfg.GameDescriptorAwsTagName)
				instances = append(instances, summary)
			}
		}
	}

	return instances, nil
}

func (s *service) GetServer(ctx context.Context, serverId string) (*models.ServerDetail, error) {
	summary, err := s.getInstanceById(ctx, serverId)

	if err != nil {
		logging.GetInstanceDetailsError(serverId, err)
		return nil, err
	}

	//TODO: Accentuate with additional detail

	detail := &models.ServerDetail{
		ServerSummary: *summary,
	}

	return detail, nil
}

func (s *service) StartServer(ctx context.Context, serverId string) (*models.Operation, error) {
	_, err := s.getInstanceById(ctx, serverId)

	if err != nil {
		logging.StartInstanceActionError(serverId, nil, err)
		return nil, err
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: []string{serverId},
	}

	_, err = s.ec2.StartInstances(ctx, input)

	if err != nil {
		logging.StartInstanceActionError(serverId, input, err)
		return nil, err
	}

	op := operationFromAction(serverId, models.OperationType("stop"))

	return &op, nil
}

func (s *service) StopServer(ctx context.Context, serverId string) (*models.Operation, error) {
	_, err := s.getInstanceById(ctx, serverId)

	if err != nil {
		logging.StopInstanceActionError(serverId, nil, err)
		return nil, err
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{serverId},
	}

	_, err = s.ec2.StopInstances(ctx, input)

	if err != nil {
		logging.StopInstanceActionError(serverId, input, err)
		return nil, err
	}

	op := operationFromAction(serverId, models.OperationType("stop"))

	return &op, nil
}

func (s *service) RebootServer(ctx context.Context, serverId string) (*models.Operation, error) {
	_, err := s.getInstanceById(ctx, serverId)

	if err != nil {
		logging.RebootInstanceActionError(serverId, nil, err)
		return nil, err
	}

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{serverId},
	}

	_, err = s.ec2.RebootInstances(ctx, input)

	if err != nil {
		logging.RebootInstanceActionError(serverId, input, err)
		return nil, err
	}

	op := operationFromAction(serverId, models.OperationType("reboot"))

	return &op, nil
}

func (s *service) GetOperation(ctx context.Context, operationId string) (*models.Operation, error) {

	//TODO implement this for real
	fakeOp := operationFromAction(uuid.NewString(), models.OperationType("start"))
	return &fakeOp, nil
}

func (s *service) getInstanceById(ctx context.Context, instanceId string) (*models.ServerSummary, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}

	output, err := s.ec2.DescribeInstances(ctx, input)
	if err != nil {
		logging.DescribeInstanceError(input, err)
		return nil, err
	}

	if len(output.Reservations) == 0 || len(output.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("Server instance not found")
	}

	instance := output.Reservations[0].Instances[0]
	summary := summaryFromInstance(instance, s.cfg.ServerNameAwsTagName, s.cfg.GameDescriptorAwsTagName)

	return &summary, nil
}

func summaryFromInstance(instance types.Instance, nameTag string, gameTag string) models.ServerSummary {
	tagsMap := mapFromTags(instance.Tags)

	summary := models.ServerSummary{
		ServerId: aws.ToString(instance.InstanceId),
		Name:     getTagValue(tagsMap, nameTag),
		Game:     getTagValue(tagsMap, gameTag),
		Specs:    models.HardwareSpecs{}, //TODO
		Status:   models.ServerStatus(string(instance.State.Name)),
	}

	return summary
}

func mapFromTags(tags []types.Tag) map[string]string {
	tagMap := make(map[string]string)

	for _, tag := range tags {
		tagMap[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tagMap
}

func getTagValue(tags map[string]string, tagkey string) string {
	value, ok := tags[tagkey]

	if ok {
		return value
	} else {
		return ""
	}
}

func operationFromAction(serverId string, oType models.OperationType) models.Operation {
	return models.Operation{
		OperationId: uuid.NewString(),
		Type:        oType,
		ServerId:    serverId,
		Status:      "issued",
		CreatedAt:   time.Now().UTC(),
		CompletedAt: nil,
		Error:       nil,
	}
}
