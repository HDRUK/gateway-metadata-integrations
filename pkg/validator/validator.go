package validator

import (
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateSchema Attempts to validate a returned json object against
// our json schema for federation services. Returns true on success,
// false otherwise. Upon error, errors are output to stdout
func ValidateSchema(document string) (bool, error) {
	var schemaUrl = os.Getenv("GMI_DEFAULT_SCHEMA_VALIDATION_URL")
	if schemaUrl == "" {
		schemaUrl = "https://raw.githubusercontent.com/HDRUK/schemata-2/master/hdr_schemata/models/GMI/gmi.schema.json"
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(document)
	fmt.Printf("schemaUrl %s: \n", schemaUrl)
	fmt.Printf("document %s: \n", document)

	fmt.Printf("schemaLoader %s: \n", schemaLoader)
	fmt.Printf("documentLoader %s: \n", documentLoader)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
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
