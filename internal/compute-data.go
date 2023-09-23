package multi_cloud_compute

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceCompute() *schema.Resource {
	return &schema.Resource{
		Schema: getVMResourceSchema(),
	}
}
