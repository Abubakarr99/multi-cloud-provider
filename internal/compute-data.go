package multi_cloud_compute

import (
	schema2 "github.com/Abubakarr99/multi-cloud-compute/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCompute() *schema.Resource {
	return &schema.Resource{
		Schema: schema2.GetVMResourceSchema(),
	}
}
