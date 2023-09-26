package cloud

import (
	"context"
	"fmt"
	vmconfig "github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AWSClient struct {
	client *session.Session
}

func (A *AWSProvider) CreateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) (string, error) {
	awsClient, ok := client.(*AWSClient)
	if !ok {
		return "", fmt.Errorf("invalid AWS client")
	}
	ec2Svc := ec2.New(awsClient.client)
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(VM.AWSAMI),
		InstanceType: aws.String(VM.InstanceType),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SubnetId:     aws.String(VM.SubnetID),
	}

	// Create the EC2 instance
	result, err := ec2Svc.RunInstancesWithContext(ctx, runInput)
	if err != nil {
		return "", err
	}
	instanceID := aws.StringValue(result.Instances[0].InstanceId)
	return instanceID, nil
}

func (A *AWSProvider) DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error {
	awsClient, ok := client.(*AWSClient)
	if !ok {
		return fmt.Errorf("invalid AWS client")
	}
	ec2Svc := ec2.New(awsClient.client)
	instanceID := aws.String(VM.ID)
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instanceID},
	}
	_, err := ec2Svc.TerminateInstancesWithContext(ctx, terminateInput)
	if err != nil {
		return err
	}
	return nil
}

func (A *AWSProvider) GetInstance(ctx context.Context, data *schema.ResourceData, client interface{}) (interface{}, error) {
	awsClient, ok := client.(*AWSClient)
	if !ok {
		return nil, fmt.Errorf("invalid aws client")
	}
	ec2Svc := ec2.New(awsClient.client)
	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(data.Id())},
	}
	// Describe the instance
	result, err := ec2Svc.DescribeInstancesWithContext(ctx, describeInput)
	if err != nil {
		return nil, err
	}

	// Check if any instances were found
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	// Extract instance details
	awsInstance := result.Reservations[0].Instances[0]
	return awsInstance, nil
}

type AWSProvider struct {
}

func (A *AWSProvider) VMtoMap(VM *vmconfig.VMConfig) map[string]interface{} {
	return map[string]interface{}{
		"name":           VM.Name,
		"instance_type":  VM.InstanceType,
		"region":         VM.Region,
		"id":             VM.ID,
		"subnet_id":      VM.SubnetID,
		"security_group": VM.AWSSecurityGroup,
	}
}

func (A *AWSProvider) UpdateInstance(ctx context.Context, new interface{}, old interface{}, client interface{}, vmConfig *vmconfig.VMConfig) error {
	//TODO implement me
	panic("implement me")
}

func (A *AWSProvider) ProviderName() string {
	return "aws"
}
func (A *AWSProvider) SetDataFromVM(VM *vmconfig.VMConfig, data *schema.ResourceData) diag.Diagnostics {
	return nil
}

func (A *AWSProvider) CreateClient(credential string) (interface{}, error) {
	creds := credentials.NewSharedCredentials(credential, "default")
	config := aws.NewConfig().WithCredentials(creds)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	awsClient := &AWSClient{
		client: sess,
	}
	return awsClient, nil
}

func (G *AWSProvider) NewInstance(instance interface{}, data *schema.ResourceData) (interface{}, error) {
	return instance.(*AWSProvider), nil
}

func (G *AWSProvider) GetInstanceConfig(instance interface{}, data *schema.ResourceData) *vmconfig.VMConfig {
	return nil
}
