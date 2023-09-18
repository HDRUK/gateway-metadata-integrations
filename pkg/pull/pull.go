package pull

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
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
func NewPull(datasetsUri, datasetUri, username, password, accessToken, method string, verbose bool) *Pull {
	pull := Pull{
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
func GetGatewayFederations() []pkg.Federation {
	req, err := http.NewRequest("GET", os.Getenv("GATEWAY_API_FEDERATIONS_URL"), nil)
	if err != nil {
		fmt.Printf("unable to create new request for gateway api pull %s\n", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		fmt.Printf("unable to pull active federations from gateway api %s\n", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to read body of response from gateway api %s\n", err)
	}

	var feds []pkg.Federation
	json.Unmarshal(body, &feds)

	return feds
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

// CallForList Attempts to authenticate against an external source and call
// recorded endpoints for data
func (p *Pull) CallForList() pkg.FederationResponse {
	req, err := http.NewRequest("GET", p.DatasetsUri, nil)
	if err != nil {
		fmt.Printf("unable to form new request with following error %s\n", err)
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		fmt.Printf("auth call failed with error %s\n", err)
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetsUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("unable to read body of call %v\n", err)
	}

	var fedList pkg.FederationResponse
	err = json.Unmarshal(body, &fedList)
	if err != nil {
		fmt.Printf("unable to unmarshal body response of call %v\n", err)
	}

	return fedList
}

// CallForDataset Is a subsequent step in the data pulling process. Issues
// an HTTP request against an individual endpoint to probe for data
func (p *Pull) CallForDataset(id string) string {
	datasetUriWithId := strings.ReplaceAll(p.DatasetUri, "{id}", id)

	req, err := http.NewRequest("GET", datasetUriWithId, nil)
	if err != nil {
		fmt.Printf("unable to form new request with following error %s\n", err)
	}

	p.GenerateHeaders(req)

	result, err := Client.Do(req)
	if err != nil {
		fmt.Printf("auth call failed with error %s\n", err)
	}

	if p.Verbose {
		fmt.Printf("running call against %s\n", p.DatasetUri)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("unable to read body of call %v\n", err)
	}

	p.Dataset = string(body)

	return p.Dataset
}

// Run Runs the functionality of this process
func Run() {
	// Firstly grab a list of all active federations in the api
	// feds := GetGatewayFederations()
	// for _, fed := range feds {
	// 	// Next gather the gcloud secrets for this federation
	// 	sec := secrets.NewSecrets("", fed.AuthSecretKey)
	// 	ret, err := sec.GetSecret()
	// 	if err != nil {
	// 		fmt.Printf("unable to retrieve secret from gcloud: %v\n", err)
	// 	}

	// 	// Create a new Pull object to action the request
	// 	p := NewPull(
	// 		fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDatasets),
	// 		fmt.Sprintf("%s%s", fed.EndpointBaseURL, fed.EndpointDataset),
	// 		"",
	// 		"",
	// 		ret.BearerToken,
	// 		fed.AuthType,
	// 		true,
	// 	)

	// 	list := p.CallForList()

	// 	for _, item := range list.Items {
	// 		_ = p.CallForDataset(item.PersistentID)
	// 	}
	// }
}
