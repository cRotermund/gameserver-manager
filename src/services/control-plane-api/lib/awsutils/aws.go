package awsutils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/logging"
)

type Client struct {
	ec2 *ec2.Client
}

type InstanceState int

const (
	Pending InstanceState = iota
	Running
	Stopping
	Stopped
	ShuttingDown
	Terminated
)

type ServerInstanceFlyweight struct {
	Id         string
	Type       string
	State      string
	LaunchTime time.Time
	Tags       map[string]string
}

type ServerInstanceDetails struct {
	ServerInstanceFlyweight
	// Add more detailed fields here as needed
}

// The default instantiator will read config in from env
func NewFromDefault() *Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic("unable to load AWS SDK config, " + err.Error())
	}

	return NewClient(cfg)
}

// Inject the client from the main context
func NewClient(cfg aws.Config) *Client {
	return &Client{
		ec2: ec2.NewFromConfig(cfg),
	}
}

func (c *Client) ListServerInstances(ctx context.Context, tagKey string, tagValue string) ([]ServerInstanceFlyweight, error) {
	// The filter "name" prefix controls the filter "subject" (e.g. tag:*)
	var tagFilterBuilder strings.Builder
	tagFilterBuilder.WriteString("tag:")
	tagFilterBuilder.WriteString(tagKey)

	//Build the describe query
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(tagFilterBuilder.String()),
				Values: []string{tagValue},
			},
		},
	}

	instPaginator := ec2.NewDescribeInstancesPaginator(c.ec2, input)

	var instances []ServerInstanceFlyweight
	for instPaginator.HasMorePages() {

		page, err := instPaginator.NextPage(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to describe instances on page: %w", err)
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				flyweight := flyweightFromAwsInstance(instance)
				instances = append(instances, flyweight)
			}
		}
	}

	return instances, nil
}

func (c *Client) GetServerInstanceById(ctx context.Context, instanceId string) (*ServerInstanceFlyweight, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}

	output, err := c.ec2.DescribeInstances(ctx, input)
	if err != nil {
		logging.DescribeInstanceError(input, err)
		return nil, err
	}

	if len(output.Reservations) == 0 || len(output.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("Server instance not found")
	}

	instance := output.Reservations[0].Instances[0]
	flyweight := flyweightFromAwsInstance(instance)

	return &flyweight, nil
}

func (c *Client) GetServerInstanceDetails(ctx context.Context, instanceId string) (*ServerInstanceDetails, error) {
	flyweight, err := c.GetServerInstanceById(ctx, instanceId)

	if err != nil {
		logging.GetInstanceDetailsError(instanceId, err)
		return nil, err
	}

	//TODO: Accentuate with additional detail

	detail := &ServerInstanceDetails{
		ServerInstanceFlyweight: *flyweight,
	}

	return detail, nil
}

func (c *Client) StartServerInstance(ctx context.Context, instanceId string) error {
	_, err := c.GetServerInstanceById(ctx, instanceId)

	if err != nil {
		logging.StartInstanceActionError(instanceId, nil, err)
		return err
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceId},
	}

	_, err = c.ec2.StartInstances(ctx, input)

	if err != nil {
		logging.StartInstanceActionError(instanceId, input, err)
		return err
	}

	return nil
}

func (c *Client) StopServerInstance(ctx context.Context, instanceId string) error {
	_, err := c.GetServerInstanceById(ctx, instanceId)

	if err != nil {
		logging.StopInstanceActionError(instanceId, nil, err)
		return err
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceId},
	}

	_, err = c.ec2.StopInstances(ctx, input)

	if err != nil {
		logging.StopInstanceActionError(instanceId, input, err)
		return err
	}

	return nil
}

func (c *Client) RebootServerInstance(ctx context.Context, instanceId string) error {
	_, err := c.GetServerInstanceById(ctx, instanceId)

	if err != nil {
		logging.RebootInstanceActionError(instanceId, nil, err)
		return err
	}

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceId},
	}

	_, err = c.ec2.RebootInstances(ctx, input)

	if err != nil {
		logging.RebootInstanceActionError(instanceId, input, err)
		return err
	}

	return nil
}

func flyweightFromAwsInstance(awsinst types.Instance) ServerInstanceFlyweight {
	flyweight := ServerInstanceFlyweight{
		Id:         aws.ToString(awsinst.InstanceId),
		Type:       string(awsinst.InstanceType),
		State:      string(awsinst.State.Name),
		LaunchTime: aws.ToTime(awsinst.LaunchTime),
		Tags:       mapFromTags(awsinst.Tags),
	}

	return flyweight
}

func mapFromTags(tags []types.Tag) map[string]string {
	tagMap := make(map[string]string)

	for _, tag := range tags {
		tagMap[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tagMap
}
