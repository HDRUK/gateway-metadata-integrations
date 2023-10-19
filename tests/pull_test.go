// This test was moved to standalone tests/ dir (which is against common practice), but
// due to cyclomatic dependency being unavoidable

package pull

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/utils/mocks"
	"hdruk/federated-metadata/pkg/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonString = `{
	"items": [{
		"name": "HDR syndication test 1 - stand cat. ",
		"@schema": "https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
		"description": "This dataset is being used to test the standard cat for syndicating with HDR UK. The dataset is set to private, with no access request. It includes 2 csv files and associated dictionaries with no lookups. ",
		"type": "dataset",
		"persistentId": "hdr_syndication_test_1___stand_cat__",
		"self": "https://fair.preview.aridhia.io/api/datasets/hdr_syndication_test_1___stand_cat__",
		"version": "1.0.0",
		"issued": "2022-10-10T09:23:48.492Z",
		"modified": "2023-03-14T13:03:17.736Z",
		"source": "FAIR PM "
	}, {
		"name": "HDR syndication test 2 - stand cat. ",
		"@schema": "https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
		"description": "This dataset is being used to test the standard cat for syndicating with HDR UK. The dataset is set to public, with no access request. It includes 1 csv file and associated dictionary with no lookups. ",
		"type": "dataset",
		"persistentId": "hdr_syndication_test_2___stand_cat__",
		"self": "https://fair.preview.aridhia.io/api/datasets/hdr_syndication_test_2___stand_cat__",
		"version": "1.0.0",
		"issued": "2022-10-10T09:45:36.074Z",
		"modified": "2023-03-14T13:03:38.401Z",
		"source": "FAIR "
	}, {
		"name": "hdr syndication test 3 - custom cat",
		"@schema": "https://raw.githubusercontent.com/HDRUK/schemata/master/schema/dataset/2.1.0/dataset.schema.json",
		"description": "Sample description here",
		"type": "dataset",
		"persistentId": "hdr_syndication_test_3___custom_cat",
		"self": "https://fair.preview.aridhia.io/api/datasets/hdr_syndication_test_3___custom_cat",
		"version": "0.0.0",
		"issued": "2023-02-13T14:10:31.640Z",
		"modified": "2023-09-04T10:53:01.239Z",
		"source": "Aridhia DRE"
	}],
	"query": {
		"q": "",
		"total": 3,
		"limit": 0,
		"offset": 0
	}
}`

func init() {
	pull.Client = &mocks.MockClient{}
}

func testGetFederations(t *testing.T) []pkg.Federation {
	jsonString := `[
		{
			"id": 1,
			"auth_type": "bearer",
			"auth_secret_key": "projects/987760029877/secrets/dev-gateway-mfs-aridhia/versions/latest",
			"endpoint_baseurl": "https://fair.preview.aridhia.io/api/syndication/hdruk",
			"endpoint_datasets": "/datasets?assigned=true",
			"endpoint_dataset": "/datasets/{id}?assigned=true",
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

	assert.EqualValues(t, "HDR syndication test 1 - stand cat. ", list.Items[0].Name)
	assert.EqualValues(t, "hdr_syndication_test_1___stand_cat__", list.Items[0].PersistentID)
	assert.EqualValues(t, "FAIR PM ", list.Items[0].Source)
	assert.EqualValues(t, "dataset", list.Items[0].Type)

	assert.EqualValues(t, "HDR syndication test 2 - stand cat. ", list.Items[1].Name)
	assert.EqualValues(t, "hdr_syndication_test_2___stand_cat__", list.Items[1].PersistentID)
	assert.EqualValues(t, "FAIR ", list.Items[1].Source)
	assert.EqualValues(t, "dataset", list.Items[1].Type)

	assert.EqualValues(t, "hdr syndication test 3 - custom cat", list.Items[2].Name)
	assert.EqualValues(t, "hdr_syndication_test_3___custom_cat", list.Items[2].PersistentID)
	assert.EqualValues(t, "Aridhia DRE", list.Items[2].Source)
	assert.EqualValues(t, "dataset", list.Items[2].Type)
}

func TestItCanValidateAgainstOurSchema(t *testing.T) {
	verdict, err := validator.ValidateSchema(jsonString)

	assert.EqualValues(t, true, verdict)
	assert.EqualValues(t, nil, err)
}
