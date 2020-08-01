package herokudata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HerokuDataAPI struct {
	APIKey              string
	URL                 string
	PollWaitDuration    time.Duration
	PollTimeoutDuration time.Duration
}

const (
	PollMaxStep           = 3
	CredentialActiveState = "active"
)

const (
	PermissionReadonly  = "readonly"
	PermissionReadWrite = "readwrite"
)

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
			if name == credential.Name && credential.State == CredentialActiveState {
				result := Credential{
					ID:       name,
					Name:     name,
					AddonID:  addonID,
					Database: credential.Database,
				}
				return &result, nil
			}
		}
	}

	// haven't found credential, return error
	return nil, fmt.Errorf("HerokuDataAPI: Credential not found")
}

func (api HerokuDataAPI) CreateCredential(addonID, name, permission string) error {
	resultRef := &apiCredentialCreateResult{}
	data := getAPICredentialCreateQuery(addonID, name)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return err
	}

	if len(resultRef.Errors) > 0 {
		return fmt.Errorf(
			"HerokuDataAPI: Create credential failed: %s",
			resultRef.Errors[0].Message,
		)
	}

	// we need to wait until the credential is created, in order to set permissions
	credential, err := api.poll(
		func() (interface{}, error) {
			return api.FetchCredential(addonID, name)
		},
		func() error {
			return fmt.Errorf("HerokuDataAPI: Create credential has timed out")
		},
	)
	if err != nil {
		return err
	}

	err = api.setPermission(credential.(*Credential), permission)
	if err != nil {
		return err
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

func (api HerokuDataAPI) setPermission(credential *Credential, permission string) error {
	if permission == "" {
		return nil
	}

	resultRef := &apiCredentialSetPermissionsResult{}
	data := getAPICredentialPermissionQuery(
		credential.AddonID, credential.Name, credential.Database, permission,
	)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return err
	}

	if len(resultRef.Errors) > 0 {
		return fmt.Errorf(
			"HerokuDataAPI: Set credential permission failed: %s",
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

func (api HerokuDataAPI) poll(
	callback func() (interface{}, error), timeoutError func() error,
) (interface{}, error) {
	start := time.Now()
	elapsed := time.Now().Sub(start)
	step := 0

	for ; elapsed < api.PollTimeoutDuration; elapsed = time.Now().Sub(start) {
		waitModifier := 1.0
		if step < PollMaxStep {
			waitModifier = float64(step) / PollMaxStep
		}
		time.Sleep(time.Duration(waitModifier) * api.PollWaitDuration)
		step += 1

		if result, err := callback(); err == nil {
			return result, nil
		}
	}

	return nil, timeoutError()
}
