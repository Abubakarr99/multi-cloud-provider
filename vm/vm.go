package vm

type VMConfig struct {
	Name            string
	Region          string
	InstanceType    string
	KeyPairName     string
	SubnetID        string // Optional for AWS
	CloudProvider   string
	CredentialPath  string
	AWSAMI          string // Optional for AWS
	GCPImageFamily  string // Optional for GCP
	GCPImageProject string // Optional for GCP
	GCPNetworkName  string // Optional for GCP
	GCPProjectID    string // Optional fot GCP
}
