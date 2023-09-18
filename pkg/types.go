package pkg

// Federation Defines the shape of a Federation object being returned
// from Gateway API
type Federation struct {
	ID               int    `json:"id"`
	AuthType         string `json:"auth_type"`
	AuthSecretKey    string `json:"auth_secret_key"`
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
