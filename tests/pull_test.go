// This test was moved to standalone tests/ dir (which is against common practice), but
// due to cyclomatic dependency being unavoidable

package pull

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/utils/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	pull.Client = &mocks.MockClient{}
}

func testGetFederations(t *testing.T) []pkg.Federation {
	jsonString := `[
		{
			"id": 1,
			"auth_type": "bearer_token",
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

func TestFederationPull(t *testing.T) {
	federations := testGetFederations(t)

	assert.EqualValues(t, 1, federations[0].ID)
	assert.EqualValues(t, 18, federations[0].Team[0].ID)
}

func TestGenerateHeaders(t *testing.T) {
	fed := testGetFederations(t)

	p := pull.NewPull(
		fmt.Sprintf("%s%s", fed[0].EndpointBaseURL, fed[0].EndpointDatasets),
		fmt.Sprintf("%s%s", fed[0].EndpointBaseURL, fed[0].EndpointDataset),
		"",
		"",
		"TEST-BEARER-TOKEN",
		fed[0].AuthType,
		false,
	)

	assert.EqualValues(t, "bearer_token", p.Method)
	assert.EqualValues(t, "TEST-BEARER-TOKEN", p.AccessToken)
}
