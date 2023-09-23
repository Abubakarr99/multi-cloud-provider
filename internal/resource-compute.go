package multi_cloud_compute

import (
	"context"
	vmconfig "github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProvider interface {
	CreateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) (string, error)
	ProviderName() string
	CreateClient(credential string) (interface{}, error)
	DeleteInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error
	GetInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error
	SetDataFromVM(VM *vmconfig.VMConfig, data *schema.ResourceData) diag.Diagnostics
	VMtoMap(VM *vmconfig.VMConfig) map[string]interface{}
	UpdateInstance(ctx context.Context, VM *vmconfig.VMConfig, client interface{}) error
}

func resourceMultiCloudCompute() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateInstance,
		DeleteContext: DeleteInstance,
		UpdateContext: UpdateInstance,
		ReadContext:   ReadInstance,
		Schema:        getVMResourceSchema(),
	}
}

func DeleteInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig, ok := m.(*ProviderConfig)
	if !ok {
		return diag.Errorf("meta is not of type CloudProvider")
	}
	providerName := data.Get("cloud_provider").(string)
	vm, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	provider := providerConfig.Provider
	client := providerConfig.Client
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
	providerName, ok := data.Get("cloud_provider").(string)
	if !ok {
		return diag.Errorf("cloud provider not a string")
	}
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
	providerName, ok := data.Get("cloud_provider").(string)
	if !ok {
		return diag.Errorf("cloud provider not a string")
	}
	vm, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	err := provider.UpdateInstance(ctx, vm, client)
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
	providerName, ok := data.Get("cloud_provider").(string)
	if !ok {
		return diag.Errorf("cloud provider not a string")
	}
	config, diags := createVMConfig(providerName, data)
	if diags.HasError() {
		return diags
	}
	provider := providerConfig.Provider
	client := providerConfig.Client
	err := provider.GetInstance(ctx, config, client)
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
			Name:           data.Get("name").(string),
			Region:         data.Get("region").(string),
			InstanceType:   data.Get("instance_type").(string),
			SubnetID:       data.Get("subnet_id").(string),
			AWSAMI:         data.Get("aws_ami_id").(string),
			CredentialPath: data.Get("credentials").(string),
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
			CredentialPath:  data.Get("credentials").(string),
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
