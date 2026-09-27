provider "atcloud" {
  endpoint = "https://api.atcloud/world/v1"
  token    = "YOUR_TOKEN"
}

resource "atcloud_compute" "vm1" {
  name = "vm-example"
  plan = "basic"
}

resource "atcloud_network" "net1" {
  name = "net-example"
  cidr = "10.0.0.0/24"
}

resource "atcloud_volume" "vol1" {
  name    = "data-volume"
  size_gb = 20
}
