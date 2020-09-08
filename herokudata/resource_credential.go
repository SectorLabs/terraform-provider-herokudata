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
		Update: resourceCredentialUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

	credential, err := api.CreateCredential(addonID, name, permission)
	if err != nil {
		return err
	}

	d.SetId(credential.ID)
	return resourceCredentialRead(d, m)
}

func resourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	log.Printf("[INFO] Fetching credential resource: %s", d.Id())

	result, err := api.FetchCredential(d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("name", result.Name)
	d.Set("addon_id", result.AddonID)
	d.Set("permission", result.Permission)
	return nil
}

func resourceCredentialDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	log.Printf("[INFO] Deleting credential resource: %s", d.Id())

	err := api.DestroyCredential(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceCredentialUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*Config).API
	log.Printf("[INFO] Updating credential resource: %s", d.Id())

	if d.HasChange("permission") {
		credential, err := api.FetchCredential(d.Id())

		if err != nil {
			return err
		}

		permission := d.Get("permission").(string)
		err = api.setPermission(credential, permission)

		if err != nil {
			return err
		}
	}

	return resourceCredentialRead(d, m)
}
