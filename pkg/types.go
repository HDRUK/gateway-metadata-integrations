package pkg

// Federation Defines the shape of a Federation object being returned
// from Gateway API
type Federation struct {
	ID               int    `json:"id"`
	PID              string `json:"pid"`
	AuthType         string `json:"auth_type"`
	EndpointBaseURL  string `json:"endpoint_baseurl"`
	EndpointDatasets string `json:"endpoint_datasets"`
	EndpointDataset  string `json:"endpoint_dataset"`
	RunTimeHour      int    `json:"run_time_hour"`
	Enabled          bool   `json:"enabled"`
	Team             []Team `json:"team"`
}

// Team Defines the shape of a Team object being returned from Gateway
// API
type Team struct {
	ID                       int    `json:"id"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
	DeletedAt                string `json:"deleted_at"`
	Name                     string `json:"name"`
	Enabled                  bool   `json:"enabled"`
	AllowsMessaging          bool   `json:"allows_messaging"`
	WorkflowEnabled          bool   `json:"workflow_enabled"`
	AccessRequestsManagement bool   `json:"access_requests_management"`
	Uses5Safes               bool   `json:"uses_5_safes"`
	IsAdmin                  bool   `json:"is_admin"`
	MemberOf                 int    `json:"member_of"`
	ContactPoint             string `json:"contact_point"`
	ApplicationFormUpdatedBy string `json:"application_form_updated_by"`
	ApplicationFormUpdatedOn string `json:"application_form_updated_on"`
	MDMFolderID              string `json:"mdm_folder_id"`
}

// FederationResponse Defines the shape of a response coming from external
// Dataset requests
type FederationResponse struct {
	Items []FederationItem `json:"items"`
}

// FederationItem Defines the shape of a federation item object, we use
// to further probe external services for Data
type FederationItem struct {
	Schema       string `json:"@schema"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	PersistentID string `json:"persistentId"`
	Self         string `json:"self"`
	Version      string `json:"version"`
	Issued       string `json:"issued"`
	Modified     string `json:"modified"`
	Source       string `json:"source"`
}

type CreateSecretRequest struct {
	Path     string `json:"path"`
	SecretID string `json:"secret_id"`
	Payload  string `json:"payload"`
}

type DeleteSecretRequest struct {
	SecretID string `json:"secret_id"`
}

type Revisions struct {
	Version        *string       `json:"version"`
	Url            *string       `json:"url"`
}

type FederationDataset struct {
	Identifier        string               `json:"identifier"`
	Version            string               `json:"version"`
	Issued             string               `json:"issued"`
	Modified           string               `json:"modified"`
	Revisions          []Revisions             `json:"revisions"`
	Summary            Summary              `json:"summary"`
	Documentation      Documentation        `json:"documentation"`
	Coverage           Coverage             `json:"coverage"`
	Provenance         Provenance           `json:"provenance"`
	Accessibility      Accessibility        `json:"accessibility"`
	Observations       []Observations       `json:"observations"`
	StructuralMetadata []StructuralMetadata `json:"structuralMetadata"`
}

type Dataset struct {
	Pid              string  						`json:"pid"`
	Version          string 					    `json:"version"`
	Metadata    	 map[string]interface{} 		`json:"metadata"`
}

type DatasetVersions struct {
	Versions []string `json:"versions"`
}

type DatasetsVersions map[string]DatasetVersions


type Summary struct {
	Title        string    `json:"title"`
	Abstract     string    `json:"abstract"`
	ContactPoint string    `json:"contactPoint"`
	Keywords     []string  `json:"keywords"`
	Publisher    Publisher `json:"publisher"`
}

type Publisher struct {
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	Description  string `json:"description"`
	ContactPoint string `json:"contactPoint"`
	MemberOf     string `json:"memberOf"`
}

type Documentation struct {
	Description string `json:"description"`
}

type Coverage struct {
	Spatial string `json:"spatial"`
}

type Provenance struct {
	Temporal Temporal `json:"temporal"`
}

type Temporal struct {
	AccrualPeriodicity string  `json:"accrualPeriodicity"`
	StartDate          string  `json:"startDate"`
	TimeLag            *string `json:"timeLag"`
}

type Accessibility struct {
	Access             Access             `json:"access"`
	Usage              Usage              `json:"usage"`
	FormatAndStandards FormatAndStandards `json:"formatAndStandards"`
}

type Access struct {
	AccessRights   string `json:"accessRights"`
	Jurisdiction   string `json:"jurisdiction"`
	DataController string `json:"dataController"`
}

type Usage struct {
	DataUseLimitation string `json:"dataUseLimitation"`
}

type FormatAndStandards struct {
	VocabularyEncodingScheme string `json:"vocabularyEncodingScheme"`
	ConformsTo               string `json:"conformsTo"`
	Language                 string `json:"language"`
	Format                   string `json:"format"`
}

type Observations struct {
	ObservedNode     string `json:"observedNode"`
	MeasuredValue    int    `json:"measuredValue"`
	DisambiguatingDescription string `json:"disambiguatingDescription"`
	ObservationDate  string `json:"observationDate"`
	MeasuredProperty string `json:"measuredProperty"`
}

type StructuralMetadata struct {
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Elements    []DataElement `json:"elements"`
}

type DataElement struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	DataType    string  `json:"dataType"`
	Sensitive   bool    `json:"sensitive"`
}
