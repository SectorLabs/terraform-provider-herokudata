package herokudata

const DefaultURL = "https://data-api.heroku.com"

type Config struct {
	API    *HerokuDataAPI
	APIKey string
	URL    string
}

func (c *Config) InitializeAPI() {
	c.API = &HerokuDataAPI{
		APIKey: c.APIKey,
		URL:    c.URL,
	}
}
