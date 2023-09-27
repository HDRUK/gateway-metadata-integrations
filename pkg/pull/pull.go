package pull

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/secrets"
	"hdruk/federated-metadata/pkg/validator"
	"io"
	"net/http"
	"os"
	"strings"
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
	Client = &http.Client{}
}

// GetFederations Retrieves a list of active federations from the gateway-api
// to run against during this pull cycle
func GetGatewayFederations() ([]pkg.Federation, error) {
	req, err := http.NewRequest("GET", os.Getenv("GATEWAY_API_FEDERATIONS_URL"), nil)
	if err != nil {
		return []pkg.Federation{}, fmt.Errorf("unable to create new request for gateway api pull %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return []pkg.Federation{}, fmt.Errorf("unable to pull active federations from gateway api %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []pkg.Federation{}, fmt.Errorf("unable to read body of response from gateway api %v", err)
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

	body := []byte(`{
		"enabled": 0,
		"tested": 0
	}`)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%d", os.Getenv("GATEWAY_API_FEDERATIONS_URL"), fed),
		bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("unable to create new request for gateway api update %v\n", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := Client.Do(req)
	if err != nil {
		fmt.Printf("unable to update federation via gateway api %v\n", err)
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to read body of response from gateway api %v\n", err)
	}

	return true
}

// GenerateHeaders Returns headers primed on the Request pointer ready
// for authentication
func (p *Pull) GenerateHeaders(req *http.Request) {
	switch strings.ToUpper(p.Method) {
	case "BEARER_TOKEN":
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken))
	case "API_KEY":
		req.Header.Add("apikey", p.AccessToken)
	default:
		fmt.Printf("unknown auth method. aborting\n")
	}
}

// TestCredentials Tests that we can access an external site given
// the provided details. Returns true if the returned status code
// is 200. False otherwise.
func (p *Pull) TestCredentials() (int, bool, error) {
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		return pkg.ERROR_INVALID_HTTP_REQUEST, false, err
	}

	p.GenerateHeaders(req)

	if p.Verbose {
		fmt.Printf("%v", req)
	}

	result, err := Client.Do(req)
	if err != nil {
		return result.StatusCode, false, err
	}
	defer result.Body.Close()

	return checkStatus(result.StatusCode)
}

// TestDatasetsEndpoint Tests that we can access an external
// datasets collection given the provided details. Returns true
// if the returned status code is 200. False otherwise.
func (p *Pull) TestDatasetsEndpoint() (int, bool, error) {
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		return pkg.ERROR_INVALID_HTTP_REQUEST, false, err
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		return result.StatusCode, false, err
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
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		fmt.Printf("unable to form new request with following error %s\n", err)
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		fmt.Printf("auth call failed with error %s\n", err)
	}
	defer result.Body.Close()

	status, successful, err := checkStatus(result.StatusCode)
	if !successful {
		InvalidateFederationDueToFailure(p.ID)
		fmt.Printf("non 200 status returned %d flagging federation as invalid. Error: %v", status, err)
		return pkg.FederationResponse{}, err
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetsUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("unable to read body of call %v\n", err)
	}

	// Ensure the returned payload from http call can be
	// validated against our schema
	_, err = validator.ValidateSchema(string(body))
	if err != nil {
		return pkg.FederationResponse{}, err
	}

	var fedList pkg.FederationResponse
	err = json.Unmarshal(body, &fedList)
	if err != nil {
		fmt.Printf("unable to unmarshal body response of call %v\n", err)
	}

	return fedList, nil
}

// CallForDataset Is a subsequent step in the data pulling process. Issues
// an HTTP request against an individual endpoint to probe for data
func (p *Pull) CallForDataset(id string) (string, error) {
	datasetUriWithId := strings.ReplaceAll(p.DatasetUri, "{id}", id)

	req, err := http.NewRequest("GET", datasetUriWithId, nil)
	if err != nil {
		return "", fmt.Errorf("unable to form new request with following error %v", err)
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth call failed with error %v", err)
	}
	defer result.Body.Close()

	status, successful, err := checkStatus(result.StatusCode)
	if !successful {
		InvalidateFederationDueToFailure(p.ID)
		fmt.Printf("non 200 status returned %d flagging federation as invalid. Error: %v", status, err)
		return "", err
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body of call %v", err)
	}

	if p.Verbose {
		p.Dataset = string(body)
	}

	return p.Dataset, nil
}

// Run Runs the functionality of this process
func Run() {
	// Removed for now, as federations are dummy and will fail
	// when calling via HTTP
	//
	// Firstly grab a list of all active federations in the api
	feds, err := GetGatewayFederations()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	for _, fed := range feds {
		// Next gather the gcloud secrets for this federation
		sec := secrets.NewSecrets("", fed.AuthSecretKey)
		ret, err := sec.GetSecret()
		if err != nil {
			fmt.Printf("unable to retrieve secret from gcloud: %v\n", err)
		}

		fmt.Printf("%v", ret)

		// Create a new Pull object to action the request
		p := NewPull(
			fed.ID,
			fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDatasets),
			fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDataset),
			"",
			"",
			ret.BearerToken,
			fed.AuthType,
			true,
		)

		list, err := p.CallForList()
		if err != nil {
			// Invalidate this federation as it has received an error
			InvalidateFederationDueToFailure(fed.ID)

			fmt.Printf("%v\n", fmt.Errorf("unable to validate provided payload against our schema: %v", err))
		}

		for _, item := range list.Items {
			datasetBody, err := p.CallForDataset(item.PersistentID)
			if err != nil {
				InvalidateFederationDueToFailure(fed.ID)
				fmt.Printf("%v\n", fmt.Errorf("unable to pull individual dataset: %v", err))
			}

			if p.Verbose {
				fmt.Printf("%s\n", datasetBody)
			}
		}
	}
}

func returnFailedValidation() (int, bool, error) {
	return 200, false, fmt.Errorf("test request failed to validate response against schema definition")
}

// checkStatus Returns based upon the received HTTP status code
// from external server request
func checkStatus(statusCode int) (int, bool, error) {
	switch statusCode {
	case 200:
		return statusCode, true, nil
	case 400:
		return statusCode, false, fmt.Errorf("request received HTTP 400 (Bad Request)")
	case 401:
		return statusCode, false, fmt.Errorf("request received HTTP 401 (Unauthorized)")
	case 403:
		return statusCode, false, fmt.Errorf("request received HTTP 403 (Forbidden)")
	case 404:
		return statusCode, false, fmt.Errorf("request received HTTP 404 (Not Found)")
	case 500:
		return statusCode, false, fmt.Errorf("request received HTTP 500 (Internal Server Error)")
	case 501:
		return statusCode, false, fmt.Errorf("request receveid HTTP 501 (Not Implemented)")
	case 503:
		return statusCode, false, fmt.Errorf("request receveid HTTP 503 (Gateway Timeout)")
	}

	return pkg.ERROR_UNKNOWN, false, nil
}
