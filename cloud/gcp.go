package cloud

import (
	"context"
	"fmt"
	vmconfig "github.com/Abubakarr99/multi-cloud-compute/vm"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"io/ioutil"
)

type GCPClient struct {
	client *compute.Service
}

type GCProvider struct{}

func (G *GCProvider) DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig) error {
	//TODO implement me
	panic("implement me")
}

func (G *GCProvider) CreateInstance(ctx context.Context, VM *vmconfig.VMConfig) (string, error) {
	gcpClient, err := G.CreateClient(VM.CredentialPath)
	if err != nil {
		return "", nil
	}
	computeService := gcpClient.(*GCPClient).client
	instance := &compute.Instance{
		Name:        VM.Name,
		MachineType: fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", VM.GCPProjectID, VM.Region, VM.InstanceType),
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: fmt.Sprintf("projects/%s/global/images/family/%s", VM.GCPImageProject, VM.GCPImageFamily),
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: fmt.Sprintf("global/networks/%s", VM.GCPNetworkName),
				AccessConfigs: []*compute.AccessConfig{
					{
						Name: "External NAT",
					},
				},
			},
		},
	}
	_, err = computeService.Instances.Insert(VM.GCPProjectID, VM.Region, instance).Context(ctx).Do()
	if err != nil {
		return "", nil
	}
	return VM.Name, nil
}

func (G *GCProvider) ProviderName() string {
	return "gcp"
}

func (G *GCProvider) CreateClient(credential string) (interface{}, error) {
	ctx := context.Background()
	credentialsJSON, err := ioutil.ReadFile(credential)
	if err != nil {
		return nil, err
	}
	cred, err := google.CredentialsFromJSON(ctx, credentialsJSON, compute.ComputeScope)
	if err != nil {
		return nil, err
	}
	sa, err := compute.NewService(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}
	gcpClient := &GCPClient{
		client: sa,
	}
	return gcpClient, nil
}
