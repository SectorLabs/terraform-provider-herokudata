package herokudata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

const DefaultURL = "https://data-api.heroku.com"

type Config struct {
	API    *HerokuDataAPI
	APIKey string
	URL    string
}

func InitializeConfig(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing HerokuData configuration")

	config := &Config{}

	if apiKey, ok := d.GetOk("api_key"); ok {
		config.APIKey = apiKey.(string)
	}

	if url, ok := d.GetOk("url"); ok {
		config.URL = url.(string)
	}

	pollTimeoutDuration, err := time.ParseDuration(d.Get("poll_timeout_duration").(string))
	if err != nil {
		return nil, err
	}

	pollWaitDuration, err := time.ParseDuration(d.Get("poll_wait_duration").(string))
	if err != nil {
		return nil, err
	}

	config.API = &HerokuDataAPI{
		APIKey:              config.APIKey,
		URL:                 config.URL,
		PollTimeoutDuration: pollTimeoutDuration,
		PollWaitDuration:    pollWaitDuration,
	}

	return config, nil
}
