package infoblox

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccInfloBoxDataSourceNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInfloBoxDataSourceNetworkConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.inflobox_network.test", "network_view_name", "default"),
				),
			},
		},
	})
}

// testAccInfloBoxDataSourceNetworkConfig: data source network base configuration
const testAccInfloBoxDataSourceNetworkConfig = `
data "infoblox_network" test {
	network_view_name = "default"
	cidr = "10.0.0.0/8"
	tenant_id = "test"
}
`
