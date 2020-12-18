package infoblox

import (
	"fmt"
	"regexp"
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
				Config: testStep1CreateSingleZone,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.acctest", "aaa.com", "default", "test"),
				),
			},
			resource.TestStep{
				Config: testStep2CreateASubDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.acctest", "aaa.com", "default", "test"),
					testAccZoneAuthExists(t, "infoblox_zone_auth.sub_acctest", "sub.aaa.com", "default", "test"),
				),
			},
			// We expect this step to fail as you can't delete a domain with sub-domains
			resource.TestStep{
				Config:      testStep3DeleteParentZone,
				ExpectError: regexp.MustCompile("Cannot delete an AuthZone that has a sub-domain"),
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.acctest", "aaa.com", "default", "test"),
					testAccZoneAuthExists(t, "infoblox_zone_auth.sub_acctest", "sub.aaa.com", "default", "test"),
				),
			},
			// This final step is to remove the sub-domain so that the state can be cleaned properly
			resource.TestStep{
				Config: testStep4DeleteSubDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneAuthExists(t, "infoblox_zone_auth.acctest", "aaa.com", "default", "test"),
				),
			},
		},
	})
}

var testStep1CreateSingleZone = fmt.Sprintf(`
	resource "infoblox_zone_auth" "acctest" {
		fqdn = "acctest.com"
		dns_view="default"
		tenant_id="test"
	}
`)

var testStep2CreateASubDomain = fmt.Sprintf(`
	resource "infoblox_zone_auth" "acctest" {
		fqdn = "acctest.com"
		dns_view="default"
		tenant_id="test"
	}

	resource "infoblox_zone_auth" "sub_acctest" {
		fqdn = "sub.acctest.com"
		dns_view="default"
		tenant_id="test"
	}
`)

var testStep3DeleteParentZone = fmt.Sprintf(`
	resource "infoblox_zone_auth" "sub_acctest" {
		fqdn = "sub.acctest.com"
		dns_view="default"
		tenant_id="test"
	}
`)

var testStep4DeleteSubDomain = fmt.Sprintf(`
	resource "infoblox_zone_auth" "acctest" {
		fqdn = "acctest.com"
		dns_view="default"
		tenant_id="test"
	}
`)

func testAccCheckZoneAuthDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_a_record" {
			continue
		}
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")
		_, err := objMgr.GetZoneAuthByRef(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error:%s - record not found", err)
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

		_, err := objMgr.GetZoneAuthByRef(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error:%s - record not found", err)
		}

		return nil
	}
}

func TestHasSubdomain(t *testing.T) {
	main := ibclient.ZoneAuth{Fqdn: "aaa.com"}
	subdomain := ibclient.ZoneAuth{Fqdn: "test.aaa.com"}
	other := ibclient.ZoneAuth{Fqdn: "foo.com"}

	list := []ibclient.ZoneAuth{main, subdomain, other}

	if hasSubdomain(main, list) == false {
		fmt.Printf("'%s' has not been identified as having a subdomain", main.Fqdn)
		t.Fail()
	}

	if hasSubdomain(other, list) == true {
		fmt.Printf("'%s' has been identified incorrectly as having a subdomain", other.Fqdn)
		t.Fail()
	}
}
