package pull

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/secrets"
	"hdruk/federated-metadata/pkg/utils"
	"hdruk/federated-metadata/pkg/validator"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// HTTPClient Defines an HTTPClient object
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

// Pull Defines a Pull object
type Pull struct {
	ID          int
	DatasetsUri string
	DatasetUri  string
	Username    string
	Password    string
	AccessToken string
	Method      string
	Verbose     bool
	Dataset     string
}

// NewPull Creates a new instance of Pull
func NewPull(id int, datasetsUri, datasetUri, username, password, accessToken, method string, verbose bool) *Pull {
	pull := Pull{
		ID:          id,
		DatasetsUri: datasetsUri,
		DatasetUri:  datasetUri,
		Username:    username,
		Password:    password,
		AccessToken: accessToken,
		Method:      method,
		Verbose:     verbose,
	}

	return &pull
}

func init() {
	_ = godotenv.Load()

	timeoutSeconds, err := strconv.Atoi(os.Getenv("FMA_DEFAULT_TIMEOUT_SECONDS"))
	if err != nil {
		utils.WriteGatewayAudit(fmt.Sprintf("unabled to determine default timeout value %v", err.Error()), "CONFIG")
		timeoutSeconds = 10
	}

	Client = &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
}

// GetFederations Retrieves a list of active federations from the gateway-api
// to run against during this pull cycle
func GetGatewayFederations() ([]pkg.Federation, error) {
	var customMsg string
	customAction := "GetGatewayFederations"

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", os.Getenv("GATEWAY_API_URL"), "federations"), nil)
	if err != nil {
		customMsg = "unable to create new request for gateway api pull"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
		return []pkg.Federation{}, fmt.Errorf("%s: %v", customMsg, err)
	}

	res, err := Client.Do(req)
	if os.IsTimeout(err) {
		customMsg = "http call timed out"
		utils.WriteGatewayAudit(fmt.Sprintf("%s %v", customMsg, err.Error()), customAction)
		return []pkg.Federation{}, fmt.Errorf("%s %v", customMsg, err)
	}
	if err != nil {
		customMsg = "unable to pull active federations from gateway api"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
		return []pkg.Federation{}, fmt.Errorf("%s: %v", customMsg, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		customMsg = "unable to read body of response from gateway api"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
		return []pkg.Federation{}, fmt.Errorf("%s: %v", customMsg, err)
	}

	var feds []pkg.Federation
	json.Unmarshal(body, &feds)

	return feds, nil
}

// InvalidateFederationDueToFailure Attempts to invalidate the federation object
// held within gateway api, due to a failure in processing. Sets enabled, tested
// to false - so that it's updated in gateway frontend and the user can determine
// the cause of the issue before testing again
func InvalidateFederationDueToFailure(fed int) bool {
	var customMsg string
	customAction := "InvalidateFederationDueToFailure"

	enabled, err := strconv.Atoi(os.Getenv("MARK_DISABLED_ON_ERROR"))
	if err != nil {
		enabled = 0
	}

	if enabled != 1 {
		return true
	}

	body := []byte(`{
		"enabled": 0,
		"tested": 0
	}`)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s/%d", os.Getenv("GATEWAY_API_URL"), "federations", fed),
		bytes.NewBuffer(body))
	if err != nil {
		customMsg = "unable to create new request for gateway api update %v"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := Client.Do(req)
	if err != nil {
		customMsg = "unable to update federation via gateway api %v"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		customMsg = "unable to read body of response from gateway api %v"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)
	}

	return true
}

// GenerateHeaders Returns headers primed on the Request pointer ready
// for authentication
func (p *Pull) GenerateHeaders(req *http.Request) {
	var customMsg string
	customAction := "GenerateHeaders"

	switch strings.ToUpper(p.Method) {
	case "BEARER":
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken))
	case "API_KEY":
		req.Header.Add("apikey", p.AccessToken)
	default:
		customMsg = fmt.Sprintf("unknown auth method %s. aborting", p.Method)
		utils.WriteGatewayAudit(customMsg, customAction)

		if p.Verbose {
			fmt.Printf("%s", customMsg)
		}
	}
}

// TestCredentials Tests that we can access an external site given
// the provided details. Returns true if the returned status code
// is 200. False otherwise.
func (p *Pull) TestCredentials() gin.H {
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		return utils.FormResponse(pkg.ERROR_INVALID_HTTP_REQUEST, false, "Credentials Test", err.Error())
	}

	p.GenerateHeaders(req)

	if p.Verbose {
		fmt.Printf("%v", req)
	}

	result, err := Client.Do(req)
	if err != nil {
		return utils.FormResponse(http.StatusBadRequest, false, "Credentials Test", err.Error())
	}
	defer result.Body.Close()

	return checkStatus(result.StatusCode)
}

// TestDatasetsEndpoint Tests that we can access an external
// datasets collection given the provided details. Returns true
// if the returned status code is 200. False otherwise.
func (p *Pull) TestDatasetsEndpoint() gin.H {
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		return utils.FormResponse(pkg.ERROR_INVALID_HTTP_REQUEST, false, "Endpoints Test", err.Error())
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		return utils.FormResponse(http.StatusBadRequest, false, "Endpoints Tests", err.Error())
	}
	defer result.Body.Close()

	list, err := p.CallForList()
	if err != nil {
		return returnFailedValidation()
	}

	if len(list.Items) >= 0 {
		return checkStatus(result.StatusCode)
	}

	return checkStatus(result.StatusCode)
}

// CallForList Attempts to authenticate against an external source and call
// recorded endpoints for data
func (p *Pull) CallForList() (pkg.FederationResponse, error) {
	var customMsg string
	customAction := "CallForList"

	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		customMsg = "unable to form new request: %v"
		utils.WriteGatewayAudit(fmt.Sprintf(customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Println(fmt.Sprintf(customMsg, err.Error()))
		}
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if os.IsTimeout(err) {
		customMsg = "http call timed out"
		utils.WriteGatewayAudit(fmt.Sprintf("%s %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("http call timedout %v", err.Error())
		}

		return pkg.FederationResponse{}, err
	}
	if err != nil {
		customMsg = "auth call failed with error: "
		utils.WriteGatewayAudit(fmt.Sprintf("%s %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("auth call failed with error %s\n", err)
		}
	}
	defer result.Body.Close()

	if !utils.IsSuccessfulStatusCode(result.StatusCode) {
		InvalidateFederationDueToFailure(p.ID)

		customMsg = "non 200 status returned %d - flagging federation as invalid. error: %v"
		utils.WriteGatewayAudit(fmt.Sprintf(customMsg, result.StatusCode, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("non 200 status returned %d flagging federation as invalid. Error: %v", result.StatusCode, err)
		}

		return pkg.FederationResponse{}, err
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetsUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		customMsg = "unable to read body of call"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("unable to read body of call %v\n", err)
		}
	}

	fmt.Println(string(body))

	// Ensure the returned payload from http call can be
	// validated against our schema
	res, err := validator.ValidateSchema(string(body))
	if !res {
		if p.Verbose {
			fmt.Printf("unable to validate incoming data against our schema %v\n", err)
		}
		if err != nil {
			return pkg.FederationResponse{}, err
		}
	}

	var fedList pkg.FederationResponse
	err = json.Unmarshal(body, &fedList)
	if err != nil {
		customMsg = "unable to unmarshal body response of call"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("unable to unmarshal body response of call %v\n", err)
		}
	}

	return fedList, nil
}

// CallForDataset Is a subsequent step in the data pulling process. Issues
// an HTTP request against an individual endpoint to probe for data
func (p *Pull) CallForDataset(id string) (pkg.FederationDataset, error) {
	var customMsg string
	customAction := "CallForDataset"

	datasetUriWithId := strings.ReplaceAll(p.DatasetUri, "{id}", id)

	req, err := http.NewRequest("GET", datasetUriWithId, nil)
	if err != nil {
		customMsg = "unable to form new request with following error"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		return pkg.FederationDataset{}, fmt.Errorf("%s %v", customMsg, err)
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if os.IsTimeout(err) {
		customMsg = "http call timed out"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("http call timedout %v", err.Error())
		}

		return pkg.FederationDataset{}, err
	}

	if err != nil {
		customMsg = "auth call failed with error"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		return pkg.FederationDataset{}, fmt.Errorf("%s %v", customMsg, err)
	}
	defer result.Body.Close()

	if !utils.IsSuccessfulStatusCode(result.StatusCode) {
		InvalidateFederationDueToFailure(p.ID)

		customMsg = "non 200 status returned %s flagging federation as invalid. Error: %v"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		if p.Verbose {
			fmt.Printf("non 200 status returned %d flagging federation as invalid. Error: %v", result.StatusCode, err)
		}

		return pkg.FederationDataset{}, err
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		customMsg = "unable to read body of call"
		utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

		return pkg.FederationDataset{}, fmt.Errorf("%s: %v", customMsg, err)
	}

	fmt.Println(string(body))

	var dataset pkg.FederationDataset
	json.Unmarshal(body, &dataset)

	return dataset, nil
}

// Run Runs the functionality of this process
func Run() {
	var customMsg string
	customAction := "Run"

	// Firstly grab a list of all active federations in the api
	feds, err := GetGatewayFederations()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	for _, fed := range feds {
		// Determine if it is time to run this federation
		// if isTimeToRun(&fed) {
		// Next gather the gcloud secrets for this federation
		sec := secrets.NewSecrets(fed.AuthSecretKey, "")
		ret, err := sec.GetSecret(fed.AuthType)
		if err != nil {
			customMsg = "unable to retrieve secrets from gcloud"
			utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

			continue
		}

		var accessToken string
		if reflect.TypeOf(ret).String() == "secrets.BearerTokenResponse" {
			accessToken = ret.(secrets.BearerTokenResponse).BearerToken
		} else if reflect.TypeOf(ret).String() == "secrets.APIKeyResponse" {
			accessToken = ret.(secrets.APIKeyResponse).APIKey
		} else { // NO_AUTH
			accessToken = ""
		}

		// Create a new Pull object to action the request
		p := NewPull(
			fed.ID,
			fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDatasets),
			fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDataset),
			"",
			"",
			accessToken,
			fed.AuthType,
			true,
		)

		list, err := p.CallForList()
		if err != nil {
			// Invalidate this federation as it has received an error
			InvalidateFederationDueToFailure(fed.ID)

			customMsg = "unable to validate provided payload against our schema"
			utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

			if p.Verbose {
				fmt.Printf("%v\n", fmt.Errorf("unable to validate provided payload against our schema: %v", err))
			}
		}

		for _, item := range list.Items {
			dataset, err := p.CallForDataset(item.PersistentID)
			if err != nil {
				InvalidateFederationDueToFailure(fed.ID)

				customMsg = "unable to pull invidual dataset"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("%v\n", fmt.Errorf("unable to pull individual dataset: %v", err))
				}
			}

			jsonString, err := json.Marshal(dataset)
			if err != nil {
				InvalidateFederationDueToFailure(fed.ID)

				customMsg = "unable to marshal dataset response to json"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("%v\n", fmt.Errorf("unable to marshal dataset to json: %v", err))
				}
			}

			//////////////////////////////////////////////////////////////////////////////////////////////////
			// TODO: This entire part will be broken for federation until federated results are in HDR 2.1.2
			// format or GWDM. Comparitively what we're currently seeing matches nothing we're working
			// towards
			//////////////////////////////////////////////////////////////////////////////////////////////////

			// Determine current status of this federated metadata.
			// i. 	Exists in DB, but not in federated payload - DELETE
			// ii. 	Exists in federated payload, but not in DB - CREATE
			// iii. Exists in both payload and DB, but different version - UPDATE
			// iv.	Exists in both payload and DB, but same version - IGNORE

			// Send the dataset to Gateway API for processing
			body := map[string]string{
				"team_id":       strconv.Itoa(fed.Team[0].ID),
				"user_id":       os.Getenv("GATEWAY_API_USER_ID"),
				"metadata":      string(jsonString),
				"create_origin": "FMA",
				"status":        "ACTIVE",
			}

			jsonPayload, _ := json.Marshal(body)

			fmt.Println(string(jsonPayload))

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("GATEWAY_API_URL"), "federations"),
				bytes.NewBuffer(jsonPayload),
			)

			if os.IsTimeout(err) {
				customMsg = "http call timed out"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("http call timedout %v", err.Error())
				}
			}
			req.Header.Add("Content-Type", "application/json")

			if err != nil {
				customMsg = "unable to prepare gateway api call with processed dataset"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("%v\n", fmt.Errorf("unable to prepare gateway api call with processed dataset: %v", err))
				}
			}
			result, err := Client.Do(req)
			if result.StatusCode > 400 {
				fmt.Printf("%v\n", fmt.Errorf("unable to call gateway api with processed dataset: %v", result))
			}
			if err != nil {
				customMsg = "unable to call gateway api with processed dataset"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("%v\n", fmt.Errorf("unable to call gateway api with processed dataset: %v", err))
				}
			}
			defer result.Body.Close()

			bodyResponse, err := io.ReadAll(result.Body)
			if err != nil {
				customMsg = "unable to parse body of gateway api dataset store call"
				utils.WriteGatewayAudit(fmt.Sprintf("%s: %v", customMsg, err.Error()), customAction)

				if p.Verbose {
					fmt.Printf("%v\n", fmt.Errorf("unable to parse body of gateway api dataset store call: %v", err))
				}
			}

			if p.Verbose {
				fmt.Println(string(bodyResponse))
			}
		}
		// continue
	}
}

func determineOperationRequired(fed *pkg.Federation, dataset *pkg.FederationDataset) {

}

func (p *Pull) writeDataset(fed *pkg.Federation, dataset *pkg.FederationDataset, jsonString *[]byte) {

}

// isTimeToRun Helper function to determine if this federation can
// run based on current time (hour) vs that of configuration in
// federated object
func isTimeToRun(fed *pkg.Federation) bool {
	var customMsg string
	customAction := "isTimeToRun"

	loc, _ := time.LoadLocation("UTC")

	dt := time.Now().In(loc)

	if dt.Hour() != fed.RunTimeHour {
		customMsg = "current federation (%d) is not ready to run (current hour: %d) vs (configured hour: %d)"
		utils.WriteGatewayAudit(fmt.Sprintf(customMsg, fed.ID, fed.RunTimeHour, dt.Hour()), customAction)
		return false
	}

	return true
}

func returnFailedValidation() gin.H {
	return utils.FormResponse(http.StatusOK, false, "Schema Validation Failed",
		fmt.Errorf("%s", "test request failed to validate response against schema definition").Error())
}

// checkStatus Returns based upon the received HTTP status code
// from external server request
func checkStatus(statusCode int) gin.H {
	switch statusCode {
	case 200:
		return utils.FormResponse(statusCode, true, "Test Successful", "nil")
	case 400:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request received HTTP 400 (Bad Request)").Error())
	case 401:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request received HTTP 401 (Unauthorized)").Error())
	case 403:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request received HTTP 403 (Forbidden)").Error())
	case 404:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request received HTTP 404 (Not Found)").Error())
	case 500:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request received HTTP 500 (Internal Server Error)").Error())
	case 501:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request receveid HTTP 501 (Not Implemented)").Error())
	case 503:
		return utils.FormResponse(statusCode, false, "Test Unsuccessful", fmt.Errorf("%s", "request receveid HTTP 503 (Gateway Timeout)").Error())
	}

	return utils.FormResponse(pkg.ERROR_UNKNOWN, false, "Test Unsuccessful", fmt.Errorf("%s", "unknown error received").Error())
}
