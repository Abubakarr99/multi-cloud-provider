package main

import (
	"github.com/Abubakarr99/multi-cloud-compute/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return multi_cloud_compute.Provider()
		},
	})
}
