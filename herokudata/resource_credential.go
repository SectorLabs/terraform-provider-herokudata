package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"permission": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					PermissionReadonly, PermissionReadWrite,
				}, true),
			},
		},
	}
}

func resourceCredentialCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	name := d.Get("name").(string)
	addonID := d.Get("addon_id").(string)
	permission := d.Get("permission").(string)
	log.Printf("[INFO] Creating credential resource: %s", name)

	err := api.CreateCredential(addonID, name, permission)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceCredentialRead(d, m)
}

func resourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	name := d.Get("name").(string)
	addonID := d.Get("addon_id").(string)
	log.Printf("[INFO] Fetching credential resource: %s", name)

	result, err := api.FetchCredential(addonID, name)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("name", result.Name)
	d.Set("addon_id", result.AddonID)
	return nil
}

func resourceCredentialDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	name := d.Get("name").(string)
	addonID := d.Get("addon_id").(string)
	log.Printf("[INFO] Deleting credential resource: %s", name)

	err := api.DestroyCredential(addonID, name)

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
