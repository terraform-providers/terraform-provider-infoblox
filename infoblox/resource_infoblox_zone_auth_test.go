package infoblox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func TestAccResourceZoneAuth(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneAuthDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceZoneAuthCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.zone_auth", "aaa.com", "default", "test"),
				),
			},
			resource.TestStep{
				Config: testAccresourceZoneAuthUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.zone_auth", "aaa.com", "default", "test"),
				),
			},
		},
	})
}

func testAccCheckZoneAuthDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_a_record" {
			continue
		}
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")
		recordName, _ := objMgr.GetZoneAuthByRef(rs.Primary.ID)
		if recordName != nil {
			return fmt.Errorf("record not found")
		}

	}
	return nil
}
func testAccZoneAuthExists(t *testing.T, n string, fqdn string, dns_view string, tenant_id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID i set")
		}
		meta := testAccProvider.Meta()
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")

		recordName, _ := objMgr.GetZoneAuthByRef(rs.Primary.ID)
		if recordName == nil {
			return fmt.Errorf("record not found")
		}

		return nil
	}
}

var testAccresourceZoneAuthCreate = fmt.Sprintf(`
resource "infoblox_zone_auth" "zone_auth"{
	fqdn = "acctest.com"
	dns_view="default"
	tenant_id="test"
	}`)

var testAccresourceZoneAuthUpdate = fmt.Sprintf(`
resource "infoblox_zone_auth" "zone_auth"{
	fqdn = "acctest.com"
	dns_view="default"
	tenant_id="test"
	}`)
