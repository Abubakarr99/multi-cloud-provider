package cloud

import (
	"context"
	multicloudcompute "github.com/Abubakarr99/multi-cloud-compute/schema"
	"github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	user2 "os/user"
	"path/filepath"
	"testing"
)

func TestGCProvider_ProviderName(t *testing.T) {
	provider := GCProvider{}
	providerName := provider.ProviderName()
	expectedProviderName := "gcp"
	if providerName != expectedProviderName {
		t.Errorf("expected provider name %s, but got %s", expectedProviderName, providerName)
	}
}

func getCredentialFilePath(file string) string {
	user, err := user2.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(user.HomeDir, file)
}

func TestGCProvider_CreateClient(t *testing.T) {
	provider := &GCProvider{}
	credentialPath := getCredentialFilePath("dantata-b059eea46359.json")
	client, err := provider.CreateClient(credentialPath)
	assert.NoError(t, err, "Createclient should not return an error")
	assert.NotNilf(t, client, "Createclient should return a non-nil client")
	_, ok := client.(*GCPClient)
	assert.True(t, ok, "client should be of type gcpclient")
}

func TestGCProvider_CreateInstance(t *testing.T) {
	provider := &GCProvider{}
	client, _ := provider.CreateClient(getCredentialFilePath("dantata-b059eea46359.json"))
	vmConfig := &vm.VMConfig{
		Name:            "toto",
		Region:          "europe-west1-b",
		InstanceType:    "e2-small",
		GCPProjectID:    "dantata",
		GCPImageFamily:  "ubuntu-2004-lts",
		GCPImageProject: "ubuntu-os-cloud",
		GCPNetworkName:  "default",
	}
	id, err := provider.CreateInstance(context.Background(), vmConfig, client)
	assert.NoError(t, err, "CreateInstance should not return an error")
	assert.NotEmpty(t, id, "CreateInstance should return a non-empty ID")
}

func TestGCProvider_DeleteInstance(t *testing.T) {
	provider := &GCProvider{}
	client, _ := provider.CreateClient(getCredentialFilePath("dantata-b059eea46359.json"))
	resourceSchema := multicloudcompute.GetVMResourceSchema()
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"gcp_project": "dantata",
		"name":        "example-vm-1",
		"region":      "europe-west1-b",
	})
	data.SetId("3007912269376857942")
	gcpInstance, err := provider.GetInstance(context.Background(), data, client)
	expectedVM := &vm.VMConfig{
		Name: "example-vm-1",
	}
	vmConfig := provider.GetInstanceConfig(gcpInstance, data)
	instance := gcpInstance.(*GCPInstance)
	assert.Equal(t, expectedVM.Name, instance.Instance.Name)
	assert.NoError(t, err, "get instance should not return an error")
	err = provider.DeleteInstance(context.Background(), vmConfig, client)
	assert.NoError(t, err, "deletion should not return an error")
}

func TestGCProvider_UpdateInstance(t *testing.T) {
	provider := &GCProvider{}
	client, _ := provider.CreateClient(getCredentialFilePath("dantata-b059eea46359.json"))
	resourceSchema := multicloudcompute.GetVMResourceSchema()
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"gcp_project":   "dantata",
		"region":        "europe-west1-b",
		"instance_type": "e2-medium",
	})
	data.SetId("1009513919837837499")
	oldInstance, err := provider.GetInstance(context.Background(), data, client)
	instance := oldInstance.(*GCPInstance)
	expectedVM := &vm.VMConfig{
		Name: "toto",
	}
	assert.Equal(t, expectedVM.Name, instance.Instance.Name)
	assert.NoError(t, err, "get instance should not return an error")
	vmConfig := provider.GetInstanceConfig(oldInstance, data)
	newinstance, err := provider.NewInstance(oldInstance, data)
	assert.NoError(t, err, "new instance should not return err")
	err = provider.UpdateInstance(context.Background(), newinstance, oldInstance, client, vmConfig)
	assert.NoError(t, err, "update instance should not return an error")
}

func TestAWSProvider_CreateClient(t *testing.T) {
	provider := AWSProvider{}
	client, err := provider.CreateClient(getCredentialFilePath("~/credentials"))
	assert.NoError(t, err, "credentials file should not be present")
	assert.NotNilf(t, client, "Createclient should return a non-nil client")
	_, ok := client.(*AWSClient)
	assert.True(t, ok, "client should be of type awsclient")
}
