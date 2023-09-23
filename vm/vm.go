package vm

type VMConfig struct {
	ID               string
	Name             string
	Region           string
	InstanceType     string
	KeyPairName      string
	SubnetID         string // Optional for AWS
	AWSSecurityGroup string
	CloudProvider    string
	CredentialPath   string
	AWSAMI           string // Optional for AWS
	GCPImageFamily   string // Optional for GCP
	GCPImageProject  string // Optional for GCP
	GCPNetworkName   string // Optional for GCP
	GCPProjectID     string // Optional fot GCP
}
