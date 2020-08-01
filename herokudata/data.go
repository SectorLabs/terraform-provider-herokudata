package herokudata

type Credential struct {
	ID      string
	Name    string
	AddonID string
}

type apiDataMap map[string]interface{}

type apiCredentialFetchResult struct {
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

type apiCredentialCreateResult struct {
	Data struct {
		CreateCredential string `json:"createCredential"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type apiCredentialDestroyResult struct {
	Data struct {
		DestroyCredential string `json:"destroyCredential"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func getAPICredentialFetchQuery(addonID string) apiDataMap {
	return apiDataMap{
		"query": `
			query FetchPostgresDetails($addonUUID: ID!) {
	 			postgres (addon_uuid: $addonUUID) {
	 				credentials_list {
	 					name
	 					uuid
	 					state
	 				}
	 			}
	 		}`,
		"variables": apiDataMap{
			"addonUUID": addonID,
		},
	}
}

func getAPICredentialCreateQuery(addonID, name string) apiDataMap {
	return apiDataMap{
		"query": `
			mutation CreateCredential($addonUUID: ID!, $name: String!) {
				createCredential(
					addon_uuid: $addonUUID
					name: $name
				)
			}`,
		"variables": apiDataMap{
			"addonUUID": addonID,
			"name":      name,
		},
	}
}

func getAPICredentialDestroyQuery(addonID, name string) apiDataMap {
	return apiDataMap{
		"query": `
			mutation DestroyCredential($addonUUID: ID!, $name: String!) {
				destroyCredential(
					addon_uuid: $addonUUID
					name: $name
				)
			}`,
		"variables": apiDataMap{
			"addonUUID": addonID,
			"name":      name,
		},
	}
}
