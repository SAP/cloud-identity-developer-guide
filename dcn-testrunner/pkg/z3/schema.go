package z3

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type InputType byte

const (
	STRING InputType = iota
	BOOLEAN
	NUMBER
	STRING_ARRAY
	NUMBER_ARRAY
	BOOLEAN_ARRAY
	STRUCTURE
	UNDEFINED
)

type Schema struct {
	inputTypes map[string]InputType
}

func SchemaFromDCN(sc []dcn.Schema) Schema {
	result := Schema{
		inputTypes: map[string]InputType{
			"$dcl":          STRUCTURE,
			"$dcl.action":   STRING,
			"$dcl.resource": STRING,
		},
	}

	for _, s := range sc {
		if s.Definition.Nested != nil {
			result.buildSchemaAttributes(s.Definition, []string{})
		}
	}
	return result
}
func (s *Schema) buildSchemaAttributes(a dcn.SchemaAttribute, path []string) {
	for k, v := range a.Nested {
		newPath := append(path, k) //nolint:gocritic
		if v.Nested != nil {
			s.inputTypes[util.StringifyQualifiedName(newPath)] = STRUCTURE
			s.buildSchemaAttributes(v, newPath)
		} else {
			s.inputTypes[util.StringifyQualifiedName(newPath)] = mapType(v.Type)
		}
	}
}

func mapType(dcnType string) InputType {
	switch dcnType {
	case "String":
		return STRING
	case "Boolean":
		return BOOLEAN
	case "Number":
		return NUMBER
	case "String[]":
		return STRING_ARRAY
	case "Boolean[]":
		return BOOLEAN_ARRAY
	case "Number[]":
		return NUMBER_ARRAY
	case "Structure":
		return STRUCTURE
	}
	return UNDEFINED
}

func (s Schema) GetTypeOfReference(ref string) InputType {
	t, ok := s.inputTypes[ref]
	if !ok {
		return UNDEFINED
	}
	return t
}
