package schema

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetVMResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the virtual machine.",
		},
		"gcp_project": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "the gcp project",
			DefaultFunc: schema.EnvDefaultFunc("GCLOUD_PROJECT", nil),
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The region where the virtual machine should be deployed.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The instance type or size of the virtual machine.",
		},
		"key_pair_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the SSH key pair for authentication.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the subnet where the virtual machine should be placed.",
		},
		"aws_security_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "the security group of the aws instance",
		},
		"aws_ami_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the AWS AMI to use for the virtual machine (AWS-specific).",
		},
		"gcp_image_family": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The image family of the GCP image to use for the virtual machine (GCP-specific).",
		},
		"gcp_image_project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The project ID of the GCP image to use for the virtual machine (GCP-specific).",
		},
		"gcp_network_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the network to attach the virtual machine to (GCP-specific).",
			DefaultFunc: schema.EnvDefaultFunc("GCLOUD_NETWORK", "default"),
		},
	}
}
