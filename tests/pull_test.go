// This test was moved to standalone tests/ dir (which is against common practice), but
// due to cyclomatic dependency being unavoidable

package pull

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/validator"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type PullTestSuite struct {
	suite.Suite
}

var jsonStringList = `{
	"items":[
		{
			"name":"HDRUK Simulation Test Server 1",
			"@schema":"https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
			"description":"Sample description here",
			"type":"dataset",
			"persistentId":"e96e36ba-30ca-4c25-bc55-fab02d72a51c",
			"self":"http://example-url.com/api/datasets/e96e36ba-30ca-4c25-bc55-fab02d72a51c",
			"version":"0.0.0",
			"issued":"2023-02-13T14:10:31.640Z",
			"modified":"2023-09-04T10:53:01.239Z",
			"source":"HDRUK"
		}
	],
	"query":{
		"q":"",
		"total":1,
		"limit":0,
		"offset":0
	}
}`

var jsonStringDataset = `{
	"@schema":"https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
	"identifier":"e96e36ba-30ca-4c25-bc55-fab02d72a51c",
	"version":"1.0.0",
	"issued":"2021-08-30T21:00:00+00:00",
	"modified":"2021-08-30T21:00:00+00:00",
	"revisions":[],
	"summary":{
		"title":"Bones Dataset",
		"abstract":"Test description",
		"publisher":{
			"identifier":"http://bones.com",
			"name":"Bones",
			"logo":"http://example.com",
			"description":"A publisher",
			"contactPoint":[],
			"memberOf":"ALLIANCE",
			"accessRights":[],
			"accessService":"Many many options",
			"accessRequestCost":"Thousands of pounds"
		},
		"contactPoint":"test@bones.com",
		"keywords":["bones","Blood"],
		"doiName":"10.1093/ajae/aaq063"
	},
	"coverage":{
		"spatial":[],
		"typicalAgeRange":"25-80"
	},
	"provenance":{
		"origin":{
			"purpose":"COVID-19",
			"source":"Routine Surveilance"
		},
		"temporal":{
			"accrualPeriodicity":"IRREGULAR",
			"distributionReleaseDate":"2021-08-30T21:00:00+00:00",
			"startDate":"2021-08-10T21:00:00+00:00",
			"endDate":"2021-08-30T21:00:00+00:00",
			"timeLag":"1-2 WEEKS"
		}
	},
	"accessibility":{
		"usage":{
			"dataUseRequirements":"none",
			"resourceCreator":"Someone",
			"investigations":[],
			"isReferencedBy":[]
		},
		"access":{
			"accessRights":["http://test.bones.com"],
			"accessService":"Many many options",
			"accessRequestCost":"Thousands",
			"jurisdiction":["GB-NIR"],
			"dataProcessor":"Some guy somewhere",
			"dataController":"This is the data controller"
		},
		"formatAndStandards":{
			"vocabularyEncodingScheme":["OPCS4"],
			"conformsTo":["HL7 CDA"],
			"language":["en"],
			"format":["audio"]
		}
	},
	"enrichmentAndLinkage":{
		"qualifiedRelation":[],
		"derivation":[],
		"tools":[]
	},
	"observations":[
		{
			"observedNode":"PERSONS",
			"measuredValue":3,
			"observationDate":"2021-08-10T21:00:00+00:00",
			"measuredProperty":"Count"
		}
	]
}`

func (t *PullTestSuite) SetUpTest() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
}

func (t *PullTestSuite) testGetFederations() []pkg.Federation {
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

	_ = json.Unmarshal(([]byte(jsonString)), &feds)
	return feds
}

func testGetList(t *PullTestSuite) pkg.FederationResponse {
	var resp pkg.FederationResponse

	err := json.Unmarshal(([]byte(jsonStringList)), &resp)
	t.Equal(nil, err)

	return resp
}

func (t *PullTestSuite) TestGetGatewayFederations() {
	federations := t.testGetFederations()

	t.Equal(1, federations[0].ID)
	t.Equal(18, federations[0].Team[0].ID)
}

func (t *PullTestSuite) TestGenerateHeaders() {
	fed := t.testGetFederations()

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

	t.Equal("bearer", p.Method)
	t.Equal("TEST-BEARER-TOKEN", p.AccessToken)
}

func (t *PullTestSuite) TestCallForList() {
	list := testGetList(t)

	t.Equal("HDRUK Simulation Test Server 1", list.Items[0].Name)
	t.Equal("e96e36ba-30ca-4c25-bc55-fab02d72a51c", list.Items[0].PersistentID)
	t.Equal("HDRUK", list.Items[0].Source)
	t.Equal("dataset", list.Items[0].Type)
}

func (t *PullTestSuite) TestItCanValidateAgainstOurSchema() {
	verdict, err := validator.ValidateSchema(jsonStringList)
	t.Nil(err)

	t.Equal(true, verdict)
	t.Equal(nil, err)
}

func TestPullTestSuite(t *testing.T) {
	suite.Run(t, new(PullTestSuite))
}

// REMOVED FOR NOW - ASK LOKI WHY...
// func TestItCanRetrieveAPIKeySecrets(t *testing.T) {
// 	p := pull.NewPull(
// 		1,
// 		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets"),
// 		fmt.Sprintf("%s%s", "https://fma-custodian-test-server-pljgro4dzq-nw.a.run.app", "/api/v1/datasets/{id}"),
// 		"",
// 		"",
// 		"FMA_UAT_fma_test_team_new",
// 		"api_key",
// 		true,
// 	)

// 	sec := secrets.NewSecrets(p.AccessToken, "")
// 	ret, err := sec.GetSecret(p.Method)

// 	assert.Equal(t, ret.(secrets.APIKeyResponse).ClientID, "ce26859054ec0c9")
// 	assert.EqualValues(t, nil, err)
// }

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
