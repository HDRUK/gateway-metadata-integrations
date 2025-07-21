package validator

import (
	"fmt"
	"hdruk/federated-metadata/pkg/utils"
	"log/slog"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateSchema Attempts to validate a returned json object against
// our json schema for federation services. Returns true on success,
// false otherwise. Upon error, errors are output to stdout
func ValidateSchema(document string, logging string) (bool, error) {
	method_name := utils.MethodName(0)
	slog.Debug(
		"ValidateSchema", 
		"x-request-session-id", logging,
		"method_name", method_name,
	)

	var schemaUrl = os.Getenv("GMI_DEFAULT_SCHEMA_VALIDATION_URL")
	if schemaUrl == "" {
		schemaUrl = "https://raw.githubusercontent.com/HDRUK/schemata-2/master/hdr_schemata/models/GMI/gmi.schema.json"
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("Error validating schema: %v", err.Error()), 
			"x-request-session-id", logging,
			"method_name", method_name,
		)
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	for _, desc := range result.Errors() {
		fmt.Printf("- %s\n", desc)
	}

	return false, err
}
