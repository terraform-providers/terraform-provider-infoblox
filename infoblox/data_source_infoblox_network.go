package infoblox

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"network_view_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of your network block.",
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network block in cidr format.",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},
		},
	}
}

// dataSourceNetworkRead: Read the network information from inflobox based on the network view name and CIDR address.
func dataSourceNetworkRead(d *schema.ResourceData, m interface{}) error {

	networkName := d.Get("network_view_name").(string)
	cidr := d.Get("cidr").(string)
	tenantID := d.Get("tenant_id").(string)

	connector := m.(*ibclient.Connector)
	ea := make(ibclient.EA)
	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetNetwork(networkName, cidr, ea)
	if err != nil {
		return fmt.Errorf("Getting Network block from network view (%s) failed : %s", networkName, err)
	}
	return networkDescriptionAttributes(d, obj)
}

// networkDescriptionAttributes: populate the numberous fields for the returned network.
func networkDescriptionAttributes(d *schema.ResourceData, obj *ibclient.Network) error {
	// simple attribute first
	d.SetId(*&obj.Ref)
	d.Set("network_view_name", obj.NetviewName)
	d.Set("cidr", obj.Cidr)
	d.Set("tenant_id", d.Get("tenant_id").(string))
	return nil
}
