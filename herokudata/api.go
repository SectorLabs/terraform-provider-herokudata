package herokudata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HerokuDataAPI struct {
	APIKey              string
	URL                 string
	PollWaitDuration    time.Duration
	PollTimeoutDuration time.Duration
}

const (
	PollMaxStep            = 3
	CredentialActiveState  = "active"
	CredentialNullState    = "null"
	CredentialRotatedState = "rotation_completed"
	IDSeparator            = "/"
	PrivilegeSelect        = "SELECT"
	PermissionObjectType   = "relation"
)

const (
	PermissionReadonly  = "readonly"
	PermissionReadWrite = "readwrite"
)

var CredentialExistsStates = map[string]bool{
	CredentialActiveState:  true,
	CredentialRotatedState: true,
	CredentialNullState:    true,
}

func (api HerokuDataAPI) FetchCredential(id string) (*Credential, error) {
	resultRef := &apiCredentialFetchResult{}

	addonID, name, err := parseID(id)
	if err != nil {
		return nil, err
	}

	data := getAPICredentialFetchQuery(addonID)
	err = api.post("graphql", data, resultRef)
	if err != nil {
		return nil, err
	}

	credentials := resultRef.Data.Postgres.CredentialsList
	if credentials != nil {
		// check if specified name exists in the list of credentials
		for _, credential := range credentials {
			if name != credential.Name {
				continue
			}

			if _, ok := CredentialExistsStates[credential.State]; !ok {
				return nil, fmt.Errorf(
					"HerokuDataAPI: Credential found, but has inactive state: %s",
					credential.State,
				)
			}

			permission, err := api.getPermission(addonID, name)
			if err != nil {
				return nil, err
			}

			result := Credential{
				ID:         makeID(addonID, name),
				Name:       name,
				AddonID:    addonID,
				Database:   credential.Database,
				Permission: permission,
			}
			return &result, nil
		}
	}

	// haven't found credential, return error
	return nil, fmt.Errorf("HerokuDataAPI: Credential not found")
}

func (api HerokuDataAPI) CreateCredential(addonID, name, permission string) (*Credential, error) {
	resultRef := &apiCredentialCreateResult{}
	data := getAPICredentialCreateQuery(addonID, name)
	err := api.post("graphql", data, resultRef)

	if err != nil {
		return nil, err
	}

	if len(resultRef.Errors) > 0 {
		return nil, fmt.Errorf(
			"HerokuDataAPI: Create credential failed: %s",
			resultRef.Errors[0].Message,
		)
	}

	id := makeID(addonID, name)
	// we need to wait until the credential is created, in order to set permissions
	credential, err := api.poll(
		func() (interface{}, error) {
			return api.FetchCredential(id)
		},
		func() error {
			return fmt.Errorf("HerokuDataAPI: Create credential has timed out")
		},
	)
	if err != nil {
		return nil, err
	}

	err = api.setPermission(credential.(*Credential), permission)
	if err != nil {
		return nil, err
	}

	return credential.(*Credential), nil
}

func (api HerokuDataAPI) DestroyCredential(id string) error {
	resultRef := &apiCredentialDestroyResult{}

	addonID, name, err := parseID(id)
	if err != nil {
		return err
	}

	data := getAPICredentialDestroyQuery(addonID, name)
	err = api.post("graphql", data, resultRef)
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

func parseID(id string) (string, string, error) {
	idParts := strings.Split(id, IDSeparator)
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("HerokuDataAPI: Given ID was not in a correct format: addon_id%sname", IDSeparator)
	}
	return idParts[0], idParts[1], nil
}

func makeID(addonID, name string) string {
	return addonID + IDSeparator + name
}

func (api HerokuDataAPI) setPermission(credential *Credential, permission string) error {
	if permission == "" {
		return nil
	}

	resultRef := &apiCredentialPermissionsSetResult{}
	data := getAPICredentialPermissionSetQuery(
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

func (api HerokuDataAPI) getPermission(addonID, name string) (string, error) {
	resultRef := &apiCredentialPermissionsGetResult{}

	data := getAPICredentialPermissionGetQuery(addonID)
	err := api.post("graphql", data, resultRef)
	if err != nil {
		return "", err
	}

	acls := resultRef.Data.PostgresSchema.DefaultAcls
	if acls != nil {
		for _, acl := range acls {
			if acl.Role != name || acl.ObjectType != PermissionObjectType {
				continue
			}
			// TODO: make this check better
			if len(acl.Privileges) == 1 && acl.Privileges[0] == PrivilegeSelect {
				return PermissionReadonly, nil
			}
			return PermissionReadWrite, nil
		}
	}

	return "", nil
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
