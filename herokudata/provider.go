package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("HEROKUDATA_API_KEY", nil),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("HEROKUDATA_API_URL", DefaultURL),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"poll_timeout_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "60s",
			},
			"poll_wait_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "4s",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"herokudata_credential": resourceCredential(),
		},

		ConfigureFunc: InitializeConfig,
	}
}
