// This test was moved to standalone tests/ dir (which is against common practice), but
// due to cyclomatic dependency being unavoidable

package pull

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/secrets"
	"hdruk/federated-metadata/pkg/validator"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var jsonString = `{
	"items": [{
		"@schema": "https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
		"type": "dataset",
		"persistentId": "edad35a6-3be0-4907-acf1-cc44b82b2342",
		"self": "https://fair.preview.aridhia.io/api/datasets/edad35a6-3be0-4907-acf1-cc44b82b2342",
		"name": "CHRIS - Return wind usually.",
		"description": "Power their however power produce woman. Section drop successful. White within factor bring wear.",
		"version": "11.0.0",
		"issued": "2010-01-18T10:34:17Z",
		"modified": "2015-12-09T05:21:42Z",
		"source": "NHSD"
	}],
	"query": {
		"q": "",
		"total": 1,
		"limit": 0,
		"offset": 0
	}
}`

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("can't read .env file. resorting to os variables\n")
	}
}

func testGetFederations(t *testing.T) []pkg.Federation {
	jsonString := `[
		{
			"id": 1,
			"auth_type": "bearer",
			"auth_secret_key": "projects/987760029877/secrets/FMA_UAT_fma_test_team_new/versions/latest",
			"endpoint_baseurl": "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app",
			"endpoint_datasets": "/api/v1/datasets",
			"endpoint_dataset": "/api/v1/datasets/{id}",
			"run_time_hour": 12,
			"enabled": true,
			"created_at": "2023-09-15T16:13:30.000000Z",
			"updated_at": "2023-09-15T16:13:34.000000Z",
			"deleted_at": null,
			"team": [
				{
					"id": 18,
					"created_at": "2023-09-15T10:43:23.000000Z",
					"updated_at": "2023-09-15T10:43:23.000000Z",
					"deleted_at": null,
					"name": "Tremblay PLC",
					"enabled": false,
					"allows_messaging": true,
					"workflow_enabled": false,
					"access_requests_management": true,
					"uses_5_safes": true,
					"is_admin": true,
					"member_of": 5049,
					"contact_point": "burdette84@schuppe.com",
					"application_form_updated_by": "aut consequatur",
					"application_form_updated_on": "1980-12-26 00:52:35",
					"mdm_folder_id": null,
					"pivot": {
						"federation_id": 1,
						"team_id": 18
					}
				}
			]
		}
	]`

	var feds []pkg.Federation

	err := json.Unmarshal(([]byte(jsonString)), &feds)
	if err != nil {
		t.Fatal(err.Error())
	}

	return feds
}

func testGetList(t *testing.T) pkg.FederationResponse {
	var resp pkg.FederationResponse

	err := json.Unmarshal(([]byte(jsonString)), &resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	return resp
}

func TestGetGatewayFederations(t *testing.T) {
	federations := testGetFederations(t)

	assert.EqualValues(t, 1, federations[0].ID)
	assert.EqualValues(t, 18, federations[0].Team[0].ID)
}

func TestGenerateHeaders(t *testing.T) {
	fed := testGetFederations(t)

	p := pull.NewPull(
		fed[0].ID,
		fmt.Sprintf("%s%s", fed[0].EndpointBaseURL, fed[0].EndpointDatasets),
		fmt.Sprintf("%s%s", fed[0].EndpointBaseURL, fed[0].EndpointDataset),
		"",
		"",
		"TEST-BEARER-TOKEN",
		fed[0].AuthType,
		false,
	)

	assert.EqualValues(t, "bearer", p.Method)
	assert.EqualValues(t, "TEST-BEARER-TOKEN", p.AccessToken)
}

func TestCallForList(t *testing.T) {
	list := testGetList(t)

	assert.EqualValues(t, "CHRIS - Return wind usually.", list.Items[0].Name)
	assert.EqualValues(t, "edad35a6-3be0-4907-acf1-cc44b82b2342", list.Items[0].PersistentID)
	assert.EqualValues(t, "NHSD", list.Items[0].Source)
	assert.EqualValues(t, "dataset", list.Items[0].Type)
}

func TestItCanValidateAgainstOurSchema(t *testing.T) {
	verdict, err := validator.ValidateSchema(jsonString)

	assert.EqualValues(t, true, verdict)
	assert.EqualValues(t, nil, err)
}

func TestItCanRetrieveAPIKeySecrets(t *testing.T) {
	p := pull.NewPull(
		1,
		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets"),
		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets/{id}"),
		"",
		"",
		"FMA_UAT_fma_test_team_new",
		"api_key",
		true,
	)

	sec := secrets.NewSecrets(p.AccessToken, "")
	ret, err := sec.GetSecret(p.Method)

	assert.Equal(t, ret.(secrets.APIKeyResponse).ClientID, "ce26859054ec0c9")
	assert.EqualValues(t, nil, err)
}

// REMOVED FOR NOW - ASK LOKI WHY...
// func TestItReturns200OnValidCredentials(t *testing.T) {
// 	sec := secrets.NewSecrets("FMA_UAT_fma_test_team_new", "")
// 	ret, err := sec.GetSecret("api_key")
// 	assert.EqualValues(t, err, nil)

// 	var accessToken string
// 	if reflect.TypeOf(ret).String() == "secrets.BearerTokenResponse" {
// 		accessToken = ret.(secrets.BearerTokenResponse).BearerToken
// 	} else if reflect.TypeOf(ret).String() == "secrets.APIKeyResponse" {
// 		accessToken = ret.(secrets.APIKeyResponse).APIKey
// 	} else { // NO_AUTH
// 		accessToken = ""
// 	}

// 	p := pull.NewPull(
// 		1,
// 		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets"),
// 		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets/{id}"),
// 		"",
// 		"",
// 		accessToken,
// 		"api_key",
// 		true,
// 	)

// 	feds, err := p.CallForList()
// 	assert.NotEmpty(t, feds)
// 	assert.NotEmpty(t, err)
// }
