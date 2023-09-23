package cloud

import (
	"context"
	"github.com/Abubakarr99/multi-cloud-compute/vm"
	"github.com/stretchr/testify/assert"
	user2 "os/user"
	"path/filepath"
	"testing"
)

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
	vmConfig := &vm.VMConfig{
		Name:         "toto",
		Region:       "europe-west1-b",
		InstanceType: "e2-small",
		GCPProjectID: "dantata",
	}
	err := provider.GetInstance(context.Background(), vmConfig, client)
	assert.NoError(t, err, "get instance should not return an error")
	err = provider.DeleteInstance(context.Background(), vmConfig, client)
	assert.NoError(t, err, "deletion should not return an error")
}
