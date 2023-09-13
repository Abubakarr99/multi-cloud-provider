package multi_cloud_compute

import (
	"context"
	"fmt"
	"github.com/Abubakarr99/multi-cloud-compute/cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cloud provider to use)",
			},
			"credentials": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Credentials for authenticating to the cloud provider",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"multi_cloud_compute": resourceMultiCloudCompute(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"multi_cloud_compute": dataSourceCompute(),
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

	providerClient, err := provider.CreateClient(credentials)
	if err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create cloud provider client",
			Detail:   "Unable to connect to the cloud provider client",
		})
	}
	return providerClient, diags
}
