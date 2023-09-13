package cloud

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	vmconfig "multi-cloud-compute/vm"
)

type AWSClient struct {
	client *session.Session
}

func (A *AWSProvider) CreateInstance(ctx context.Context, VM *vmconfig.VMConfig) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (A *AWSProvider) DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig) error {
	panic("toto")
}

type AWSProvider struct {
}

func (A *AWSProvider) ProviderName() string {
	return "aws"
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
