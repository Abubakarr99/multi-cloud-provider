package multi_cloud_compute

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"multi-cloud-compute/cloud"
)

type VMConfig struct {
	Name            string
	Region          string
	InstanceType    string
	KeyPairName     string
	SubnetID        string // Optional for AWS
	CloudProvider   string
	CredentialPath  string
	AWSAMI          string // Optional for AWS
	GCPImageFamily  string // Optional for GCP
	GCPImageProject string // Optional for GCP
	GCPNetworkName  string // Optional for GCP
	GCPProjectID    string // Optional fot GCP
}

type CloudProvider interface {
	CreateInstance(ctx context.Context, VM *VMConfig) (string, error)
	ProviderName() string
	CreateClient(credential string) (interface{}, error)
	DeleteInstance(ctx context.Context, VM *VMConfig) error
}

func resourceMultiCloudCompute() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateInstance,
		DeleteContext: DeleteInstance,
		Schema:        getVMResourceSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func DeleteInstance(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerName := data.Get("cloud_provider").(string)
	var provider CloudProvider
	var vm *VMConfig
	switch providerName {
	case "aws":
		provider = &cloud.AWSProvider{}
		vm = &VMConfig{
			Name:           data.Get("name").(string),
			Region:         data.Get("region").(string),
			InstanceType:   data.Get("instance_type").(string),
			SubnetID:       data.Get("subnet_id").(string),
			AWSAMI:         data.Get("aws_ami_id").(string),
			CredentialPath: data.Get("credentials").(string),
		}
	case "gcp":
		provider = &cloud.GCProvider{}
		vm = &VMConfig{
			Name:            data.Get("name").(string),
			Region:          data.Get("region").(string),
			InstanceType:    data.Get("instance_type").(string),
			GCPImageFamily:  data.Get("gcp_image_family").(string),
			GCPImageProject: data.Get("gcp_image_project").(string),
			GCPNetworkName:  data.Get("gcp_network_name").(string),
			CredentialPath:  data.Get("credentials").(string),
		}
	default:
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unsupported cloud provider",
			Detail:   fmt.Sprintf("The selected cloud provider '%s' is not supported", providerName),
		})
	}
	err := provider.DeleteInstance(ctx, vm)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func CreateInstance(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerName := data.Get("cloud_provider").(string)
	var provider CloudProvider
	var diags diag.Diagnostics
	var vm *VMConfig
	switch providerName {
	case "aws":
		provider = &cloud.AWSProvider{}
		vm = &VMConfig{
			Name:           data.Get("name").(string),
			Region:         data.Get("region").(string),
			InstanceType:   data.Get("instance_type").(string),
			SubnetID:       data.Get("subnet_id").(string),
			AWSAMI:         data.Get("aws_ami_id").(string),
			CredentialPath: data.Get("credentials").(string),
		}
	case "gcp":
		provider = &cloud.GCProvider{}
		vm = &VMConfig{
			Name:            data.Get("name").(string),
			Region:          data.Get("region").(string),
			InstanceType:    data.Get("instance_type").(string),
			GCPImageFamily:  data.Get("gcp_image_family").(string),
			GCPImageProject: data.Get("gcp_image_project").(string),
			GCPNetworkName:  data.Get("gcp_network_name").(string),
			CredentialPath:  data.Get("credentials").(string),
		}
	default:
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unsupported cloud provider",
			Detail:   "this cloud provider is not supported",
		})
	}
	id, err := provider.CreateInstance(ctx, vm)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(id)
	return diags
}