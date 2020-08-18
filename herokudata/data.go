package herokudata

type Credential struct {
	ID         string
	Name       string
	AddonID    string
	Database   string
	Permission string
}

type apiDataMap map[string]interface{}

type apiDataList []interface{}

type apiCredentialFetchResult struct {
	Data struct {
		Postgres struct {
			CredentialsList []struct {
				Name     string `json:"name"`
				UUID     string `json:"uuid"`
				Database string `json:"database"`
				State    string `json:"state"`
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

type apiCredentialPermissionsSetResult struct {
	Data struct {
		SetPostgresPermissions bool `json:"setPostgresPermissions"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type apiCredentialPermissionsGetResult struct {
	Data struct {
		PostgresSchema struct {
			DefaultAcls []struct {
				Role       string   `json:"role"`
				ObjectType string   `json:"object_type"`
				Privileges []string `json:"privileges"`
			} `json:"default_acls"`
		} `json:"postgresSchema"`
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
						database
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

func getAPICredentialPermissionGetQuery(addonID string) apiDataMap {
	return apiDataMap{
		"query": `
			query FetchSchema($addonUUID: ID!) {
				postgresSchema(addon_uuid: $addonUUID) {
					default_acls {...DefaultACLFragment}
				}
			}
			fragment DefaultACLFragment on PostgresSchemaDefaultACL {
				object_type
				role
				privileges
			}`,
		"variables": apiDataMap{
			"addonUUID": addonID,
		},
	}
}

func getAPICredentialPermissionSetQuery(addonID, name, database, permission string) apiDataMap {
	var databasePrivileges, tablePrivileges, sequencePrivileges apiDataList

	if permission == PermissionReadonly {
		databasePrivileges = apiDataList{"CONNECT"}
		tablePrivileges = apiDataList{"SELECT"}
		sequencePrivileges = apiDataList{"SELECT"}
	}
	if permission == PermissionReadWrite {
		databasePrivileges = apiDataList{"CONNECT", "TEMPORARY"}
		tablePrivileges = apiDataList{"SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE"}
		sequencePrivileges = apiDataList{"SELECT", "USAGE"}
	}

	return apiDataMap{
		"query": `
			mutation SetPostgresPermissions($addonUUID: ID!, $role: String!, $acls: [PostgresACLInput]!) {
				setPostgresPermissions(
					addon_uuid: $addonUUID
					role: $role
					acls: $acls
				)
			}`,
		"variables": apiDataMap{
			"addonUUID": addonID,
			"role":      name,
			"acls": apiDataList{
				apiDataMap{
					"kind":       "database",
					"name":       database,
					"privileges": databasePrivileges,
				},
				apiDataMap{
					"kind":       "table",
					"default":    true,
					"privileges": tablePrivileges,
				},
				apiDataMap{
					"kind":       "sequence",
					"default":    true,
					"privileges": sequencePrivileges,
				},
				apiDataMap{
					"kind":       "schema",
					"name":       "public",
					"privileges": apiDataList{"USAGE"},
				},
				apiDataMap{
					"kind":       "table",
					"schema":     "public",
					"privileges": tablePrivileges,
				},
				apiDataMap{
					"kind":       "sequence",
					"schema":     "public",
					"privileges": sequencePrivileges,
				},
			},
		},
	}
}
