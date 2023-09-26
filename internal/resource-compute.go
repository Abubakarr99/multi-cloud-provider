package multi_cloud_compute

import (
	"context"
	schema2 "github.com/Abubakarr99/multi-cloud-compute/schema"
	vmconfig "github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProvider interface {
	CreateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) (string, error)
	NewInstance(instance interface{}, data *schema.ResourceData) (interface{}, error)
	GetInstanceConfig(instance interface{}, data *schema.ResourceData) *vmconfig.VMConfig
	ProviderName() string
	CreateClient(credential string) (interface{}, error)
	DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error
	GetInstance(ctx context.Context, data *schema.ResourceData, client interface{}) (interface{}, error)
	SetDataFromVM(VM *vmconfig.VMConfig, data *schema.ResourceData) diag.Diagnostics
	VMtoMap(VM *vmconfig.VMConfig) map[string]interface{}
	UpdateInstance(ctx context.Context, new interface{}, old interface{}, client interface{}, vmConfig *vmconfig.VMConfig) error
}

func resourceMultiCloudCompute() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateInstance,
		DeleteContext: DeleteInstance,
		UpdateContext: UpdateInstance,
		ReadContext:   ReadInstance,
		Schema:        schema2.GetVMResourceSchema(),
	}
}

func DeleteInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig, ok := m.(*ProviderConfig)
	provider := providerConfig.Provider
	providerName := provider.ProviderName()
	client := providerConfig.Client
	if !ok {
		return diag.Errorf("meta is not of type CloudProvider")
	}
	vm, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	err := provider.DeleteInstance(ctx, vm, client)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func CreateInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig, ok := m.(*ProviderConfig)
	if !ok {
		return diag.Errorf("meta is not of type CloudProvider")
	}
	provider := providerConfig.Provider
	client := providerConfig.Client
	providerName := provider.ProviderName()
	vm, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	id, err := provider.CreateInstance(ctx, vm, client)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(id)
	return nil
}

func UpdateInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig, ok := m.(*ProviderConfig)
	if !ok {
		return diag.Errorf("meta is not of type CloudProvider")
	}
	provider := providerConfig.Provider
	client := providerConfig.Client
	oldInstance, err := provider.GetInstance(ctx, data, client)
	newInstance, err := provider.NewInstance(oldInstance, data)
	providerName := provider.ProviderName()
	vm, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	err = provider.UpdateInstance(ctx, newInstance, oldInstance, client, vm)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ReadInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig, ok := m.(*ProviderConfig)
	if !ok {
		return diag.Errorf("meta is not of type CloudProvider")
	}
	provider := providerConfig.Provider
	client := providerConfig.Client
	instanceResource, err := provider.GetInstance(ctx, data, client)
	config := provider.GetInstanceConfig(instanceResource, data)
	if err != nil {
		return diag.FromErr(err)
	}
	return provider.SetDataFromVM(config, data)
}

func createVMConfig(providerName string, data *schema.ResourceData) (*vmconfig.VMConfig, diag.Diagnostics) {
	var vm *vmconfig.VMConfig
	var diags diag.Diagnostics

	switch providerName {
	case "aws":
		vm = &vmconfig.VMConfig{
			Name:         data.Get("name").(string),
			Region:       data.Get("region").(string),
			InstanceType: data.Get("instance_type").(string),
			SubnetID:     data.Get("subnet_id").(string),
			AWSAMI:       data.Get("aws_ami_id").(string),
		}
	case "gcp":
		vm = &vmconfig.VMConfig{
			Name:            data.Get("name").(string),
			Region:          data.Get("region").(string),
			InstanceType:    data.Get("instance_type").(string),
			GCPImageFamily:  data.Get("gcp_image_family").(string),
			GCPImageProject: data.Get("gcp_image_project").(string),
			GCPNetworkName:  data.Get("gcp_network_name").(string),
			GCPProjectID:    data.Get("gcp_project").(string),
		}
	default:
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unsupported cloud provider",
			Detail:   "this cloud provider is not supported",
		})
	}

	return vm, diags
}
