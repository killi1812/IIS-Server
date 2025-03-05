package xmlvalidator

import (
	"errors"
	"fmt"

	"github.com/killi1812/libxml2"
	"github.com/killi1812/libxml2/relaxng"
	"github.com/killi1812/libxml2/types"
	"github.com/killi1812/libxml2/xsd"
)

type ValidationMethod = string

const (
	XSD ValidationMethod = "xsd"
	RNG ValidationMethod = "rng"
)

func Validate(data []byte, method ValidationMethod) error {
	var schema types.Schema
	switch method {
	case XSD:
		s, err := validateXsd()
		if err != nil {
			return err
		}
		schema = s

	case RNG:
		s, err := validateRng()
		if err != nil {
			return err
		}
		schema = s
	default:
		return errors.New("validation method not supported")
	}

	defer schema.Free()

	doc, err := libxml2.Parse(data)
	if err != nil {
		fmt.Printf("error parsing XML: %v\nType:%T\n", err, err)
		return &ErrInvalidXML{Reason: err.Error()}
	}
	defer doc.Free() // Free memory

	if err := schema.Validate(doc); err != nil {
		// TODO: return ErrInvalidXML istead of  types.SchemaValidationError
		return err
	}

	fmt.Println("XML is valid!")
	return nil
}

const (
	_XSD_FILE_NAME = "schemas/schema.xsd"
	_RNG_FILE_NAME = "schemas/schema.rng"
)

type ErrInvalidXML struct {
	Reason string
}

func (e *ErrInvalidXML) Error() string {
	return e.Reason
}

func validateXsd() (types.Schema, error) {
	// NOTE: writes a nice output to console but doesn't have great error returns
	schema, err := xsd.ParseFromFile(_XSD_FILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("error parsing XSD: %v", err)
	}
	return schema, nil
}

func validateRng() (types.Schema, error) {
	// NOTE: Not really the best error output
	schema, err := relaxng.ParseFromFile(_RNG_FILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("error parsing RNG: %v", err)
	}
	return schema, nil
}
