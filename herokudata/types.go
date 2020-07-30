package herokudata

type JSONDataMap map[string]interface{}

type ReadCredentialResult struct {
	ID      string
	Name    string
	AddonID string
}

type ReadCredentialJSONResult struct {
	Data struct {
		Postgres struct {
			CredentialsList []struct {
				Name  string `json:"name"`
				UUID  string `json:"uuid"`
				State string `json:"state"`
			} `json:"credentials_list"`
		} `json:"postgres"`
	} `json:"data"`
}

type CreateCredentialJSONResult struct {
	Data struct {
		CreateCredential string `json:"createCredential"`
	} `json:"data"`
}
