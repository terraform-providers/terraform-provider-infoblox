package infoblox

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func resourceZoneAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneAuthCreate,
		Read:   resourceZoneAuthGet,
		Update: resourceZoneAuthUpdate,
		Delete: resourceZoneAuthDelete,

		Schema: map[string]*schema.Schema{

			"fqdn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fqdn of the auth zone to create.",
			},

			"dns_view": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Dns View under which the zone has been created.",
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},
		},
	}
}

func resourceZoneAuthCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning to create auth zone from  required network block", resourceZoneAuthIDString(d))

	fqdn := d.Get("fqdn").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	ea := make(ibclient.EA)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	ZoneAuth, err := objMgr.CreateZoneAuth(fqdn, ea)

	if err != nil {
		return fmt.Errorf("Error creating auth zone (%s): %s", fqdn, err)
	}

	d.SetId(ZoneAuth.Ref)

	log.Printf("[DEBUG] %s: Creation of auth zone complete", resourceZoneAuthIDString(d))

	return nil
	return resourceZoneAuthGet(d, m)
}

func resourceZoneAuthGet(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning to Get auth zone", resourceZoneAuthIDString(d))

	fqdn := d.Get("fqdn").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetZoneAuthByRef(d.Id())
	if err != nil {
		return fmt.Errorf("Getting auth zone failed from dns view (%s) : %s", fqdn, err)
	}
	d.SetId(obj.Ref)

	log.Printf("[DEBUG] %s: Completed reading required auth zone ", resourceZoneAuthIDString(d))
	return nil
}

func resourceZoneAuthUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Updating an auth zone is not supported")
}

func resourceZoneAuthDelete(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning Deletion of auth zone", resourceZoneAuthIDString(d))

	fqdn := d.Get("fqdn").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	zaList, err := objMgr.GetZoneAuth()
	if err != nil {
		return fmt.Errorf("Getting a list of all current AuthZones failed")
	}

	if hasSubdomain(ibclient.ZoneAuth{Fqdn: fqdn}, zaList) {
		return fmt.Errorf("Cannot delete an AuthZone that has a sub-domain: %s", fqdn)
	}

	_, err = objMgr.DeleteZoneAuth(d.Id())
	if err != nil {
		return fmt.Errorf("Deletion of auth zone failed from dns view(%s) : %s", fqdn, err)
	}
	d.SetId("")

	log.Printf("[DEBUG] %s: Deletion of auth zone complete", resourceZoneAuthIDString(d))
	return nil
}

type resourceZoneAuthIDStringInterface interface {
	Id() string
}

func resourceZoneAuthIDString(d resourceZoneAuthIDStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("infoblox_auth_zone (ID = %s)", id)
}

func hasSubdomain(target ibclient.ZoneAuth, list []ibclient.ZoneAuth) bool {
	for _, za := range list {
		if za.Fqdn != target.Fqdn && strings.Contains(za.Fqdn, target.Fqdn) {
			return true
		}
	}
	return false
}
