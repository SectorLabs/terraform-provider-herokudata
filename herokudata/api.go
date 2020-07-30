package herokudata

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HerokuDataAPI struct {
	APIKey string
	URL    string
}

func (api HerokuDataAPI) ReadCredential(addonID, name string) (*ReadCredentialResult, error) {
	resultRef := &ReadCredentialJSONResult{}
	data := getCredentialReadQuery(addonID)
	err := api.Post("graphql", data, resultRef)

	if err != nil {
		return nil, err
	}

	credentials := resultRef.Data.Postgres.CredentialsList
	if credentials == nil {
		return nil, nil
	}

	// check if specified name exists in the list of credentials
	for _, credential := range credentials {
		// TODO: check state
		if name == credential.Name {
			result := ReadCredentialResult{
				ID: name,
				Name: name,
				AddonID: addonID,
			}
			return &result, nil
		}
	}
	return nil, nil
}

func (api HerokuDataAPI) CreateCredential(addonID, name string) (bool, error) {
	resultRef := &CreateCredentialJSONResult{}
	data := getCredentialCreateQuery(addonID, name)
	err := api.Post("graphql", data, resultRef)

	if err != nil {
		return false, err
	}

	// TODO: create permissions

	return resultRef.Data.CreateCredential != "", nil
}

func (api HerokuDataAPI) Post(path string, data JSONDataMap, resultRef interface{}) error {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodPost,
		api.URL+"/"+path,
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Add("Authorization", "Bearer "+api.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// if response is not needed, return
	if resultRef == nil {
		return nil
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, resultRef); err != nil {
		return err
	}

	return nil
}

func getCredentialCreateQuery(addonID, name string) JSONDataMap {
	return JSONDataMap{
		"query": `
			mutation CreateCredential($addonUUID: ID!, $name: String!) {
				createCredential(
					addon_uuid: $addonUUID
					name: $name
				)
			}`,
		"variables": JSONDataMap{
			"addonUUID": addonID,
			"name":      name,
		},
	}
}

func getCredentialReadQuery(addonID string) JSONDataMap {
	return JSONDataMap{
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
		"variables": JSONDataMap{
			"addonUUID": addonID,
		},
	}
}
