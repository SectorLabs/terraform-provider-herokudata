package herokudata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HerokuDataAPI struct {
	APIKey string
	URL    string
}

func (api HerokuDataAPI) FetchCredential(addonID, name string) (*Credential, error) {
	resultRef := &apiCredentialFetchResult{}
	data := getAPICredentialFetchQuery(addonID)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return nil, err
	}

	credentials := resultRef.Data.Postgres.CredentialsList
	if credentials != nil {
		// check if specified name exists in the list of credentials
		for _, credential := range credentials {
			// TODO: check state
			if name == credential.Name {
				result := Credential{
					ID:      name,
					Name:    name,
					AddonID: addonID,
				}
				return &result, nil
			}
		}
	}

	// haven't found credential, return error
	return nil, fmt.Errorf("HerokuDataAPI: Credential not found.")
}

func (api HerokuDataAPI) CreateCredential(addonID, name string) error {
	resultRef := &apiCredentialCreateResult{}
	data := getAPICredentialCreateQuery(addonID, name)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return err
	}

	// TODO: create permissions

	if len(resultRef.Errors) > 0 {
		return fmt.Errorf(
			"HerokuDataAPI: Create credential failed: %s",
			resultRef.Errors[0].Message,
		)
	}
	return nil
}

func (api HerokuDataAPI) DestroyCredential(addonID, name string) error {
	resultRef := &apiCredentialDestroyResult{}
	data := getAPICredentialDestroyQuery(addonID, name)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return err
	}

	if len(resultRef.Errors) > 0 {
		return fmt.Errorf(
			"HerokuDataAPI: Destroy credential failed: %s",
			resultRef.Errors[0].Message,
		)
	}
	return nil
}

func (api HerokuDataAPI) post(path string, data apiDataMap, resultRef interface{}) error {
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
