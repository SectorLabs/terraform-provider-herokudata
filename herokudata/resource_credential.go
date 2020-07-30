package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCreate,
		Read:   resourceCredentialRead,
		Delete: resourceCredentialDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"addon_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCredentialCreate(d *schema.ResourceData, m interface{}) error {
	log.Print("[INFO] Creating credential resource")

	api := m.(*Config).API
	name := d.Get("name").(string)
	addonID := d.Get("addon_id").(string)

	ok, err := api.CreateCredential(addonID, name)

	if err == nil && ok {
		d.SetId(name)
	}

	return resourceCredentialRead(d, m)
}

func resourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	log.Print("[INFO] Checking if credential resource exists")

	api := m.(*Config).API
	name := d.Get("name").(string)
	addonID := d.Get("addon_id").(string)

	result, err := api.ReadCredential(addonID, name)

	if err != nil || result == nil {
		log.Print("[WARN] Resource not found")
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("addon_id", result.AddonID)

	return nil
}

func resourceCredentialDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: implement
	return nil
}
