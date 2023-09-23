package multi_cloud_compute

import (
	"context"
	"fmt"
	"github.com/Abubakarr99/multi-cloud-compute/cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	Provider CloudProvider
	Client   interface{}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cloud provider to use)",
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_PROVIDER", nil),
			},
			"credentials": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Credentials for authenticating to the cloud provider",
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_CREDS", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudfusion_server": resourceMultiCloudCompute(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	cloudProvider := data.Get("cloud_provider").(string)
	credentials := data.Get("credentials").(string)
	var provider CloudProvider
	switch cloudProvider {
	case "aws":
		provider = &cloud.AWSProvider{}
	case "gcp":
		provider = &cloud.GCProvider{}
	default:
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unsupported cloud provider",
			Detail:   fmt.Sprintf("The selected cloud provider '%s' is not supported", cloudProvider),
		})
	}
	client, err := provider.CreateClient(credentials)
	if err != nil {
		diag.FromErr(err)
	}
	providerConfig := &ProviderConfig{
		Provider: provider,
		Client:   client,
	}
	return providerConfig, diags
}
