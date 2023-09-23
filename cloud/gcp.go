package cloud

import (
	"context"
	"errors"
	"fmt"
	vmconfig "github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io/ioutil"
	"strconv"
	"time"
)

type GCPClient struct {
	client *compute.Service
}

type GCProvider struct{}

func (G *GCProvider) DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error {
	computeService := client.(*GCPClient).client
	op, err := computeService.Instances.Delete(VM.GCPProjectID, VM.Region, VM.ID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("error deleting instance %s", err.Error())
	}
	err = G.waitForOperation(ctx, client, VM.GCPProjectID, VM.Region, op.Name)
	if err != nil {
		return err
	}
	return nil
}

func (G *GCProvider) VMtoMap(VM *vmconfig.VMConfig) map[string]interface{} {
	return map[string]interface{}{
		"name":          VM.Name,
		"instance_type": VM.InstanceType,
		"region":        VM.Region,
		"id":            VM.ID,
	}
}

func (G *GCProvider) UpdateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error {
	computeService := client.(*GCPClient).client
	instance := &compute.Instance{
		Name:        VM.Name,
		MachineType: fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", VM.GCPProjectID, VM.Region, VM.InstanceType),
	}
	_, err := computeService.Instances.Update(VM.GCPProjectID, VM.Region, VM.Name, instance).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (G *GCProvider) CreateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) (string, error) {
	computeService := client.(*GCPClient).client
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
	op, err := computeService.Instances.Insert(VM.GCPProjectID, VM.Region, instance).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	err = G.waitForOperation(ctx, client, VM.GCPProjectID, VM.Region, op.Name)
	if err != nil {
		return "", err
	}
	return VM.Name, nil
}

func (G *GCProvider) ProviderName() string {
	return "gcp"
}

func (G *GCProvider) GetInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error {
	computeService := client.(*GCPClient).client
	project := VM.GCPProjectID
	region := VM.Region
	instance, err := computeService.Instances.Get(project, region, VM.Name).Context(ctx).Do()
	if err != nil {
		var gceErr *googleapi.Error
		if errors.As(err, &gceErr) && gceErr.Code == 404 {
			fmt.Printf("instance not found: %s", VM.Name)
			return nil
		}
		return err
	}
	VM.ID = strconv.FormatUint(instance.Id, 10)
	return nil
}

func (G *GCProvider) SetDataFromVM(VM *vmconfig.VMConfig, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	vmMap := G.VMtoMap(VM)
	for k, v := range vmMap {
		if err := data.Set(k, v); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %s", k, err)...)
		}
	}
	return diags
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

func (G *GCProvider) waitForOperation(ctx context.Context, client interface{}, projectID, Zone, operationName string) error {
	computeService := client.(*GCPClient).client
	for {
		operation, err := computeService.ZoneOperations.Get(projectID, Zone, operationName).Do()
		if err != nil {
			return err
		}

		if operation.Status == "DONE" {
			if operation.Error != nil {
				return fmt.Errorf("operation failed: %v", operation.Error.Errors)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Wait for a few seconds before checking the operation status again
			// Adjust the sleep duration as needed
			<-time.After(5 * time.Second)
		}
	}
}
