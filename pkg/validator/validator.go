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
	schemaLoader := gojsonschema.NewReferenceLoader(os.Getenv("FMA_DEFAULT_SCHEMA_VALIDATION_URL"))
	documentLoader := gojsonschema.NewStringLoader(document)

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
