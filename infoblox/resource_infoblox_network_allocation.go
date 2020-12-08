package infoblox

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func resourceNetworkAllocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkAllocationCreate,
		Read:   resourceNetworkAllocationRead,
		Update: resourceNetworkAllocationUpdate,
		Delete: resourceNetworkAllocationDelete,

		Schema: map[string]*schema.Schema{
			"network_view_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Network view name available in NIOS Server.",
			},
			"network_name": &schema.Schema{
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
			"reserve_ip": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The no of IP's you want to reserve.",
			},
			"gateway": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "gateway ip address of your network block.By default first IPv4 address is set as gateway address.",
				Computed:    true,
			},
			"prefix_length": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Prefix length of the network to be allocated.",
			},
		},
	}
}

func resourceNetworkAllocationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning network block Creation", resourceNetworkAllocationIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	cidr := d.Get("cidr").(string)
	networkName := d.Get("network_name").(string)
	reserveIP := d.Get("reserve_ip").(int)
	gateway := d.Get("gateway").(string)
	prefixLength := d.Get("prefix_length").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	ZeroMacAddr := "00:00:00:00:00:00"
	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)
	ea := make(ibclient.EA)

	prefixUint, err := strconv.ParseUint(prefixLength, 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting prefix to uint : %s", err)
	}

	nwname, err := objMgr.AllocateNetwork(networkViewName, cidr, uint(prefixUint), networkName)
	if err != nil {
		return fmt.Errorf("Creation of network block failed in network view (%s) : %s", networkViewName, err)
	}

	// Check whether gateway or ip address already allocated
	gatewayIP, err := objMgr.GetFixedAddress(networkViewName, nwname.Cidr, gateway, "")
	if err == nil && gatewayIP != nil {
		fmt.Printf("Gateway already created")
	} else if gatewayIP == nil {
		gatewayIP, err = objMgr.AllocateIP(networkViewName, nwname.Cidr, gateway, ZeroMacAddr, "", ea)
		if err != nil {
			return fmt.Errorf("Gateway Creation failed in network block(%s) error: %s", nwname.Cidr, err)
		}
	}

	for i := 1; i <= reserveIP; i++ {
		_, err = objMgr.AllocateIP(networkViewName, nwname.Cidr, gateway, ZeroMacAddr, "", ea)
		if err != nil {
			return fmt.Errorf("Reservation in network block failed in network view(%s):%s", networkViewName, err)
		}
	}

	d.Set("gateway", gatewayIP.IPAddress)
	d.SetId(nwname.Ref)

	log.Printf("[DEBUG] %s: Creation on network block complete", resourceNetworkAllocationIDString(d))
	return resourceNetworkRead(d, m)
}
func resourceNetworkAllocationRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Reading the required network block", resourceNetworkAllocationIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetNetworkwithref(d.Id())
	if err != nil {
		return fmt.Errorf("Getting Network block from network view (%s) failed : %s", networkViewName, err)
	}
	d.SetId(obj.Ref)
	log.Printf("[DEBUG] %s: Completed reading network block", resourceNetworkAllocationIDString(d))
	return nil
}
func resourceNetworkAllocationUpdate(d *schema.ResourceData, m interface{}) error {

	return fmt.Errorf("network updation is not supported")
}

func resourceNetworkAllocationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Deletion of network block", resourceNetworkAllocationIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	_, err := objMgr.DeleteNetwork(d.Id(), d.Get("network_view_name").(string))
	if err != nil {
		return fmt.Errorf("Deletion of Network block failed from network view(%s): %s", networkViewName, err)
	}
	d.SetId("")

	log.Printf("[DEBUG] %s: Deletion of network block complete", resourceNetworkAllocationIDString(d))
	return nil
}

type resourceNetworkAllocationIDStringInterface interface {
	Id() string
}

func resourceNetworkAllocationIDString(d resourceNetworkAllocationIDStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("infoblox_ip_allocation (ID = %s)", id)
}
