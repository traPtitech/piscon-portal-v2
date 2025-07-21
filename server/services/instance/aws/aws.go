package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance"
)

var _ instance.Manager = (*Client)(nil)

type Client struct {
	client *ec2.Client
	cfg    Config
}

type Config struct {
	ImageID      string
	InstanceType string
	Region       string

	AccessKey string
	SecretKey string

	SubnetID        string
	SecurityGroupID string
	KeyPairName     string
}

func NewClient(cfg Config) (*Client, error) {
	awsConfigFuncs := make([]func(*config.LoadOptions) error, 0, 2)

	awsConfigFuncs = append(awsConfigFuncs, config.WithRegion(cfg.Region))
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		awsConfigFuncs = append(awsConfigFuncs, config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     cfg.AccessKey,
					SecretAccessKey: cfg.SecretKey,
				},
			},
		))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), awsConfigFuncs...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	return &Client{
		client: ec2.NewFromConfig(awsCfg),
		cfg:    cfg,
	}, nil
}

var pisconTagFilter = types.Filter{
	Name:   lo.ToPtr("tag:piscon"),
	Values: []string{"true"},
}

func (a *Client) Create(ctx context.Context, name string, sshPubKeys []string) (string, error) {
	tagSpec := types.TagSpecification{
		ResourceType: types.ResourceTypeInstance,
		Tags: []types.Tag{
			{
				Key:   lo.ToPtr("Name"),
				Value: &name,
			},
			{
				Key:   lo.ToPtr("piscon"),
				Value: lo.ToPtr("true"),
			},
		},
	}

	cloudConfig := CloudConfig{
		Users: []User{
			{
				Name:              "isucon",
				Sudo:              "ALL=(ALL) NOPASSWD:ALL",
				Groups:            []string{"sudo"},
				SSHAuthorizedKeys: sshPubKeys,
			},
		},
	}
	userdata, err := cloudConfig.ConvertToUserData()
	if err != nil {
		return "", fmt.Errorf("generate user data: %w", err)
	}

	nispec := types.InstanceNetworkInterfaceSpecification{
		AssociatePublicIpAddress: lo.ToPtr(true),
		DeleteOnTermination:      lo.ToPtr(true),
		DeviceIndex:              lo.ToPtr[int32](0),
		SubnetId:                 &a.cfg.SubnetID,
		Groups:                   []string{a.cfg.SecurityGroupID},
	}
	instanceInput := &ec2.RunInstancesInput{
		ImageId:           &a.cfg.ImageID,
		InstanceType:      types.InstanceType(a.cfg.InstanceType),
		MinCount:          lo.ToPtr[int32](1),
		MaxCount:          lo.ToPtr[int32](1),
		TagSpecifications: []types.TagSpecification{tagSpec},
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{nispec},
		KeyName:           &a.cfg.KeyPairName,
		UserData:          lo.ToPtr(userdata),
	}

	res, err := a.client.RunInstances(ctx, instanceInput)
	if err != nil {
		return "", fmt.Errorf("run instances: %w", err)
	}

	instance := res.Instances[0]
	return *instance.InstanceId, nil
}

func (a *Client) Delete(ctx context.Context, instance domain.InfraInstance) error {
	_, err := a.client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{instance.ProviderInstanceID},
	})
	if err != nil {
		return fmt.Errorf("terminate instance: %w", err)
	}

	return nil
}

func (a *Client) Get(ctx context.Context, id string) (domain.InfraInstance, error) {
	res, err := a.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return domain.InfraInstance{}, fmt.Errorf("describe instance: %w", err)
	}

	if len(res.Reservations) == 0 || len(res.Reservations[0].Instances) == 0 {
		return domain.InfraInstance{}, fmt.Errorf("instance not found: %s", id)
	}

	instance := res.Reservations[0].Instances[0]

	return domain.InfraInstance{
		ProviderInstanceID: *instance.InstanceId,
		PrivateIP:          instance.PrivateIpAddress,
		PublicIP:           instance.PublicIpAddress,
		Status:             convertInstanceState(instance.State.Name),
	}, nil
}

func (a *Client) GetByIDs(ctx context.Context, ids []string) ([]domain.InfraInstance, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	res, err := a.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: ids,
		Filters:     []types.Filter{pisconTagFilter},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instances: %w", err)
	}

	instances := make([]domain.InfraInstance, 0, len(res.Reservations))
	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			infraInstance := domain.InfraInstance{
				ProviderInstanceID: *instance.InstanceId,
				PrivateIP:          instance.PrivateIpAddress,
				PublicIP:           instance.PublicIpAddress,
				Status:             convertInstanceState(instance.State.Name),
			}

			instances = append(instances, infraInstance)
		}
	}

	return instances, nil
}

func (a *Client) GetAll(ctx context.Context) ([]domain.InfraInstance, error) {
	res, err := a.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{pisconTagFilter},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instances: %w", err)
	}

	var instances []domain.InfraInstance
	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, domain.InfraInstance{
				ProviderInstanceID: *instance.InstanceId,
				PrivateIP:          instance.PrivateIpAddress,
				PublicIP:           instance.PublicIpAddress,
				Status:             convertInstanceState(instance.State.Name),
			})
		}
	}

	return instances, nil
}

func (a *Client) Stop(ctx context.Context, instance domain.InfraInstance) error {
	_, err := a.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []string{instance.ProviderInstanceID},
	})
	if err != nil {
		return fmt.Errorf("stop instance: %w", err)
	}

	return nil
}

func (a *Client) Start(ctx context.Context, instance domain.InfraInstance) error {
	_, err := a.client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{instance.ProviderInstanceID},
	})
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}

	return nil
}

func convertInstanceState(state types.InstanceStateName) domain.InstanceStatus {
	switch state {
	case types.InstanceStateNamePending:
		return domain.InstanceStatusPending
	case types.InstanceStateNameRunning:
		return domain.InstanceStatusRunning
	case types.InstanceStateNameStopping:
		return domain.InstanceStatusStopping
	case types.InstanceStateNameStopped:
		return domain.InstanceStatusStopped
	case types.InstanceStateNameShuttingDown:
		return domain.InstanceStatusDeleting
	case types.InstanceStateNameTerminated:
		return domain.InstanceStatusDeleted
	default:
		return domain.InstanceStatusUnknown
	}
}
