package codegen

import (
	"fmt"
	"slices"

	"github.com/dave/jennifer/jen"

	"github.com/srdtrk/go-codegen/pkg/schemas"
	"github.com/srdtrk/go-codegen/pkg/types"
)

func GenerateResponses(f *jen.File, responses map[string]*schemas.JSONSchema) {
	for key, schema := range responses {
		RegisterDefinitions(schema.Definitions)
		switch {
		case len(schema.Type) == 1:
			switch schema.Type[0] {
			case schemas.TypeNameObject:
				if schema.Title == "" {
					panic(fmt.Sprintf("response schema for %s must have a title", key))
				}
				duplicate, found := GetDefinition(schema.Title)
				if found {
					if duplicate.Description != schema.Description || !slices.Contains([]string{key, schema.Title}, duplicate.Title) {
						types.DefaultLogger().Warn().Msgf("found duplicate definition `%s` with differing implementations", schema.Title)
						types.DefaultLogger().Warn().Msgf("renaming the duplicate definition to `%s`", schema.Title+"_2")

						schema.Title += "_2"
					}
				}

				RegisterDefinition(schema.Title, schema)
			case schemas.TypeNameArray:
				if len(schema.Definitions) == 0 {
					types.DefaultLogger().Error().Msgf("response schema for %s must have a definition for array, skipping... Please create an issue in https://github.com/srdtrk/go-codegen", key)
				}

				RegisterDefinition(key, schema)
			case schemas.TypeNameString:
				// Do nothing
			case schemas.TypeNameNumber:
				// Do nothing
			case schemas.TypeNameInteger:
				// Do nothing
			case schemas.TypeNameBoolean:
			// Do nothing
			case schemas.TypeNameNull:
				types.DefaultLogger().Warn().Msgf("response schema for %s is of type null", key)
			default:
				types.DefaultLogger().Error().Msgf("response schema for %s is of unknown type %s", key, schema.Type[0])
			}
		case len(schema.Type) == 2 && slices.Contains(schema.Type, schemas.TypeNameNull):
			switch {
			case slices.Contains(schema.Type, schemas.TypeNameString):
				// Do nothing
			case slices.Contains(schema.Type, schemas.TypeNameNumber):
				// Do nothing
			case slices.Contains(schema.Type, schemas.TypeNameInteger):
				// Do nothing
			case slices.Contains(schema.Type, schemas.TypeNameBoolean):
				// Do nothing
			default:
				types.DefaultLogger().Error().Msgf("response schema for %s is not supported, skipping... Please create an issue in https://github.com/srdtrk/go-codegen", key)
			}
		case len(schema.OneOf) != 0:
			RegisterDefinition(key, schema)
		case len(schema.AllOf) == 1:
			RegisterDefinition(key, schema)
		case schema.Ref != nil:
			RegisterDefinition(key, schema)
		case len(schema.AnyOf) != 0:
			RegisterDefinition(key, schema)
			RegisterDefinitions(schema.Definitions)
		default:
			types.DefaultLogger().Error().Msgf("response schema for %s is not supported, skipping... Please create an issue in https://github.com/srdtrk/go-codegen", key)
		}
	}
}
