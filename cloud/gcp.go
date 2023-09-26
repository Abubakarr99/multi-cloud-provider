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
type GCPInstance struct {
	Instance *compute.Instance
}

type GCProvider struct{}

func (G *GCProvider) DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error {
	computeService := client.(*GCPClient).client
	op, err := computeService.Instances.Delete(VM.GCPProjectID, VM.Region, VM.Name).Context(ctx).Do()
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

func (G *GCProvider) UpdateInstance(ctx context.Context, new interface{}, old interface{}, client interface{}, vmConfig *vmconfig.VMConfig) error {
	computeService := client.(*GCPClient).client
	newInstance := new.(*GCPInstance).Instance
	oldInstance := old.(*GCPInstance).Instance
	_, err := computeService.Instances.Update(vmConfig.GCPProjectID, vmConfig.Region, oldInstance.Name, newInstance).Context(ctx).Do()
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
	createInstance, err := computeService.Instances.Get(VM.GCPProjectID, VM.Region, VM.Name).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	instanceID := strconv.FormatUint(createInstance.Id, 10)
	return instanceID, nil
}

func (G *GCProvider) ProviderName() string {
	return "gcp"
}

func (G *GCProvider) GetInstance(ctx context.Context, data *schema.ResourceData, client interface{}) (interface{}, error) {
	computeService := client.(*GCPClient).client
	project := data.Get("gcp_project").(string)
	region := data.Get("region").(string)
	name := data.Get("name").(string)
	instance, err := computeService.Instances.Get(project, region, data.Id()).Context(ctx).Do()
	if err != nil {
		var gceErr *googleapi.Error
		if errors.As(err, &gceErr) && gceErr.Code == 404 {
			fmt.Printf("instance not found: %s", name)
			return nil, nil
		}
		return nil, err
	}
	return &GCPInstance{Instance: instance}, nil
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

func (G *GCProvider) NewInstance(instance interface{}, data *schema.ResourceData) (interface{}, error) {
	newInstance, ok := instance.(*GCPInstance)
	if !ok {
		return nil, fmt.Errorf("expected GCPInstance got %T", instance)
	}
	project := data.Get("gcp_project").(string)
	zone := data.Get("region").(string)
	machineType := data.Get("instance_type").(string)
	newInstance.Instance.MachineType = fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", project, zone, machineType)
	return newInstance, nil
}

func (G *GCProvider) GetInstanceConfig(instance interface{}, data *schema.ResourceData) *vmconfig.VMConfig {
	newInstance := instance.(*GCPInstance).Instance
	return &vmconfig.VMConfig{
		Name:         newInstance.Name,
		InstanceType: newInstance.MachineType,
		GCPProjectID: data.Get("gcp_project").(string),
		Region:       data.Get("region").(string),
	}
}
