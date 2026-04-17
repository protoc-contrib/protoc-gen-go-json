package generator

import (
	"fmt"
	"strconv"
)

// Options controls how MarshalJSON and UnmarshalJSON methods are emitted.
//
// All options are exposed as plugin parameters; they can be set via
// `--go-json_opt=<name>=<value>` on the protoc command line, or the equivalent
// `opt:` stanza in buf.gen.yaml.
type Options struct {
	// EnumsAsInts renders enum values as integers instead of their string names.
	// Maps to protojson.MarshalOptions.UseEnumNumbers.
	EnumsAsInts bool

	// EmitDefaults emits fields whose value is the zero value for their type.
	// Maps to protojson.MarshalOptions.EmitUnpopulated.
	EmitDefaults bool

	// OrigName uses the original .proto field names instead of the lowerCamelCase
	// JSON names. Maps to protojson.MarshalOptions.UseProtoNames.
	OrigName bool

	// AllowUnknownFields discards unknown fields during unmarshaling instead of
	// erroring. Maps to protojson.UnmarshalOptions.DiscardUnknown.
	AllowUnknownFields bool
}

// Set applies a single `name=value` plugin parameter to the options. The signature
// matches what protogen.Options.ParamFunc expects.
func (o *Options) Set(name, value string) error {
	switch name {
	case "enums_as_ints":
		return parseBool(value, &o.EnumsAsInts)
	case "emit_defaults":
		return parseBool(value, &o.EmitDefaults)
	case "orig_name":
		return parseBool(value, &o.OrigName)
	case "allow_unknown":
		return parseBool(value, &o.AllowUnknownFields)
	default:
		return fmt.Errorf("unknown plugin option %q", name)
	}
}

func parseBool(value string, dst *bool) error {
	if value == "" {
		*dst = true
		return nil
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid boolean value %q: %w", value, err)
	}
	*dst = b
	return nil
}
