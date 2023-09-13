package multi_cloud_compute

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getVMResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the virtual machine.",
		},
		"credentials": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The path to the credentials file.",
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
		"cloud_provider": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cloud provider where the virtual machine should be created (e.g., 'aws' or 'gcp').",
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
		},
	}
}
