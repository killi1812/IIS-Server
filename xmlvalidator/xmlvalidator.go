package xmlvalidator

import (
	"fmt"
	"os"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

type ValidationMethod = string

const (
	XSD ValidationMethod = "xsd"
	RNG ValidationMethod = "rng"
)

func Validate(data []byte, method ValidationMethod) error {
	switch method {
	case XSD:
		return validateXsd(data)

	case RNG:
		return validateRng(data)
	}
	return nil
}

const (
	_XSD_FILE_NAME = "schemas/schema.xsd"
	_RNG_FILE_NAME = "schema.xsd"
)

type ErrInvalidXML struct {
	Reason string
}

func (e *ErrInvalidXML) Error() string {
	return e.Reason
}

func validateXsd(data []byte) error {
	schemaData, err := os.ReadFile(_XSD_FILE_NAME)
	if err != nil {
		return fmt.Errorf("error reading XSD file: %v", err)
	}

	schema, err := xsd.Parse(schemaData)
	if err != nil {
		return fmt.Errorf("error parsing XSD: %v", err)
	}
	defer schema.Free() // Free memory after validation

	// Parse XML document
	doc, err := libxml2.Parse(data)
	if err != nil {
		fmt.Printf("error parsing XML: %v", err)
		return &ErrInvalidXML{Reason: err.Error()}
	}
	defer doc.Free() // Free memory

	// Validate XML against XSD
	if err := schema.Validate(doc); err != nil {
		// TODO: check how to unwrap errors

		fmt.Printf("%T\n", err)
		return err
	}

	fmt.Println("XML is valid!")
	return nil
}

// TODO: Change to struct of xml
func validateRng(data any) error {
	return nil
}
