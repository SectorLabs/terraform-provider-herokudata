package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
		},

		ResourcesMap: map[string]*schema.Resource{
			"herokudata_credential": resourceCredential(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing HerokuData provider")
	config := &Config{}

	if apiKey, ok := d.GetOk("api_key"); ok {
		config.APIKey = apiKey.(string)
	}

	if url, ok := d.GetOk("url"); ok {
		config.URL = url.(string)
	}

	config.InitializeAPI()

	return config, nil
}
