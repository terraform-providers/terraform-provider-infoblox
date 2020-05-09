provider "infoblox" {
}

data "infoblox_network" "network" {
  network_view_name = "default"
  cidr = "10.0.0.0/8"
  tenant_id = ""
}
