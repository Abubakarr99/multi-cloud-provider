
provider "cloudfusion" {
  cloud_provider = "gcp"
  credentials = "/Users/abubakarrkamara/dantata-b059eea46359.json"
}

resource "cloudfusion_server" "toto" {
  name          = "example-vm"
  region        = "europe-west1-b"
  instance_type = "e2-small"
  gcp_image_family = "ubuntu-2004-lts"
  gcp_image_project = "ubuntu-os-cloud"
  gcp_project = "dantata"
}

