package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HEROKUDATA_API_KEY", nil),
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HEROKUDATA_API_URL", DefaultURL),
			},
			"poll_timeout_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "60s",
			},
			"poll_wait_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2s",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"herokudata_credential": resourceCredential(),
		},

		ConfigureFunc: InitializeConfig,
	}
}
