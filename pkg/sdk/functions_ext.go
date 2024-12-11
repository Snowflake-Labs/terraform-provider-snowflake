package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

const DefaultFunctionComment = "user-defined function"

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// FunctionDetails contains aggregated describe results for the given function.
type FunctionDetails struct {
	Signature                  string  // present for all function types
	Returns                    string  // present for all function types
	Language                   string  // present for all function types
	Body                       *string // present for all function types (hidden when SECURE)
	NullHandling               *string // present for all function types but SQL
	Volatility                 *string // present for all function types but SQL
	ExternalAccessIntegrations *string // list present for python, java, and scala
	Secrets                    *string // map present for python, java, and scala
	Imports                    *string // list present for python, java, and scala (hidden when SECURE)
	Handler                    *string // present for python, java, and scala (hidden when SECURE)
	RuntimeVersion             *string // present for python, java, and scala (hidden when SECURE)
	Packages                   *string // list // present for python, java, and scala
	TargetPath                 *string // present for scala and java (hidden when SECURE)
	InstalledPackages          *string // list present for python (hidden when SECURE)
	IsAggregate                *bool   // present for python

	NormalizedImports    []NormalizedPath
	NormalizedTargetPath *NormalizedPath
	ReturnDataType       datatypes.DataType
	NormalizedArguments  []NormalizedArgument
}

type NormalizedPath struct {
	// StageLocation is a normalized (fully-quoted id or `~`) stage location
	StageLocation string
	// PathOnStage is path to the file on stage without opening `/`
	PathOnStage string
}

type NormalizedArgument struct {
	name         string
	dataType     datatypes.DataType
	defaultValue string // TODO [next PR]: handle when adding default values
}

func functionDetailsFromRows(rows []FunctionDetail) (*FunctionDetails, error) {
	v := &FunctionDetails{}
	var errs []error
	for _, row := range rows {
		switch row.Property {
		case "signature":
			errs = append(errs, row.setStringValueOrError("signature", &v.Signature))
		case "returns":
			errs = append(errs, row.setStringValueOrError("returns", &v.Returns))
		case "language":
			errs = append(errs, row.setStringValueOrError("language", &v.Language))
		case "null handling":
			v.NullHandling = row.Value
		case "volatility":
			v.Volatility = row.Value
		case "body":
			v.Body = row.Value
		case "external_access_integrations":
			v.ExternalAccessIntegrations = row.Value
		case "secrets":
			v.Secrets = row.Value
		case "imports":
			v.Imports = row.Value
		case "handler":
			v.Handler = row.Value
		case "runtime_version":
			v.RuntimeVersion = row.Value
		case "packages":
			v.Packages = row.Value
		case "installed_packages":
			v.InstalledPackages = row.Value
		case "is_aggregate":
			errs = append(errs, row.setOptionalBoolValueOrError("is_aggregate", &v.IsAggregate))
		case "target_path":
			v.TargetPath = row.Value
		}
	}
	if e := errors.Join(errs...); e != nil {
		return nil, e
	}

	if functionDetailsImports, err := parseFunctionDetailsImport(*v); err != nil {
		errs = append(errs, err)
	} else {
		v.NormalizedImports = functionDetailsImports
	}

	if v.TargetPath != nil {
		if p, err := parseStageLocationPath(*v.TargetPath); err != nil {
			errs = append(errs, err)
		} else {
			v.NormalizedTargetPath = p
		}
	}

	if dt, returnNotNull, err := parseFunctionAndProcedureReturns(v.Returns); err != nil {
		errs = append(errs, err)
	} else {
		v.ReturnDataType = dt
		_ = returnNotNull // TODO [next PR]: used when adding return nullability to the resource
	}

	if args, err := parseFunctionAndProcedureSignature(v.Signature); err != nil {
		errs = append(errs, err)
	} else {
		v.NormalizedArguments = args
	}

	return v, errors.Join(errs...)
}

func parseFunctionDetailsImport(details FunctionDetails) ([]NormalizedPath, error) {
	functionDetailsImports := make([]NormalizedPath, 0)
	if details.Imports == nil || *details.Imports == "" || *details.Imports == "[]" {
		return functionDetailsImports, nil
	}
	if !strings.HasPrefix(*details.Imports, "[") || !strings.HasSuffix(*details.Imports, "]") {
		return functionDetailsImports, fmt.Errorf("could not parse imports from Snowflake: %s, wrapping brackets not found", *details.Imports)
	}
	raw := (*details.Imports)[1 : len(*details.Imports)-1]
	imports := strings.Split(raw, ",")
	for _, imp := range imports {
		p, err := parseStageLocationPath(imp)
		if err != nil {
			return nil, fmt.Errorf("could not parse imports from Snowflake: %s, err: %w", *details.Imports, err)
		}
		functionDetailsImports = append(functionDetailsImports, *p)
	}
	return functionDetailsImports, nil
}

func parseStageLocationPath(location string) (*NormalizedPath, error) {
	log.Printf("[DEBUG] parsing stage location path part: %s", location)
	idx := strings.Index(location, "/")
	if idx < 0 {
		return nil, fmt.Errorf("part %s cannot be split into stage and path", location)
	}
	stageRaw := strings.TrimPrefix(strings.TrimSpace(location[:idx]), "@")
	if stageRaw != "~" {
		stageId, err := ParseSchemaObjectIdentifier(stageRaw)
		if err != nil {
			return nil, fmt.Errorf("part %s contains incorrect stage location: %w", location, err)
		}
		stageRaw = stageId.FullyQualifiedName()
	}
	pathRaw := strings.TrimPrefix(strings.TrimSpace(location[idx:]), "/")
	if pathRaw == "" {
		return nil, fmt.Errorf("part %s contains empty path", location)
	}
	return &NormalizedPath{stageRaw, pathRaw}, nil
}

func parseFunctionAndProcedureReturns(returns string) (datatypes.DataType, bool, error) {
	var returnNotNull bool
	trimmed := strings.TrimSpace(returns)
	if strings.HasSuffix(trimmed, " NOT NULL") {
		returnNotNull = true
		trimmed = strings.TrimSuffix(trimmed, " NOT NULL")
	}
	dt, err := datatypes.ParseDataType(trimmed)
	return dt, returnNotNull, err
}

// Format in Snowflake DB is: (argName argType [DEFAULT defaultValue], argName argType [DEFAULT defaultValue], ...).
func parseFunctionAndProcedureSignature(signature string) ([]NormalizedArgument, error) {
	normalizedArguments := make([]NormalizedArgument, 0)
	trimmed := strings.TrimSpace(signature)
	if trimmed == "" {
		return normalizedArguments, fmt.Errorf("could not parse signature from Snowflake: %s, can't be empty", signature)
	}
	if trimmed == "()" {
		return normalizedArguments, nil
	}
	if !strings.HasPrefix(trimmed, "(") || !strings.HasSuffix(trimmed, ")") {
		return normalizedArguments, fmt.Errorf("could not parse signature from Snowflake: %s, wrapping parentheses not found", trimmed)
	}
	raw := (trimmed)[1 : len(trimmed)-1]
	args := strings.Split(raw, ",")

	for _, arg := range args {
		a, err := parseFunctionOrProcedureArgument(arg)
		if err != nil {
			return nil, fmt.Errorf("could not parse signature from Snowflake: %s, err: %w", trimmed, err)
		}
		normalizedArguments = append(normalizedArguments, *a)
	}
	return normalizedArguments, nil
}

// TODO [next PR]: adjust after tests for strange arg names and defaults
func parseFunctionOrProcedureArgument(arg string) (*NormalizedArgument, error) {
	log.Printf("[DEBUG] parsing argument: %s", arg)
	trimmed := strings.TrimSpace(arg)
	idx := strings.Index(trimmed, " ")
	if idx < 0 {
		return nil, fmt.Errorf("arg %s cannot be split into arg name, data type, and default", arg)
	}
	argName := trimmed[:idx]
	rest := strings.TrimSpace(trimmed[idx:])
	split := strings.Split(rest, " DEFAULT ")
	var dt datatypes.DataType
	var defaultValue string
	var err error
	switch len(split) {
	case 1:
		dt, err = datatypes.ParseDataType(split[0])
		if err != nil {
			return nil, fmt.Errorf("arg type %s cannot be parsed, err: %w", split[0], err)
		}
	case 2:
		dt, err = datatypes.ParseDataType(split[0])
		if err != nil {
			return nil, fmt.Errorf("arg type %s cannot be parsed, err: %w", split[0], err)
		}
		defaultValue = strings.TrimSpace(split[1])
	default:
		return nil, fmt.Errorf("cannot parse arg %s, part: %s", arg, rest)
	}
	return &NormalizedArgument{argName, dt, defaultValue}, nil
}

func (v *functions) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return functionDetailsFromRows(rows)
}

func (v *functions) ShowParameters(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Function: id,
		},
	})
}

func (d *FunctionDetail) setStringValueOrError(property string, field *string) error {
	if d.Value == nil {
		return fmt.Errorf("value expected for field %s", property)
	} else {
		*field = *d.Value
	}
	return nil
}

func (d *FunctionDetail) setOptionalBoolValueOrError(property string, field **bool) error {
	if d.Value != nil && *d.Value != "" {
		v, err := strconv.ParseBool(*d.Value)
		if err != nil {
			return fmt.Errorf("invalid value for field %s, err: %w", property, err)
		} else {
			*field = Bool(v)
		}
	}
	return nil
}

func (s *CreateForJavaFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForJavaFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func (s *CreateForPythonFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForPythonFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func (s *CreateForScalaFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForScalaFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func NewCreateForSQLFunctionRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	returns FunctionReturnsRequest,
	functionDefinition string,
) *CreateForSQLFunctionRequest {
	s := CreateForSQLFunctionRequest{}
	s.name = name
	s.Returns = returns
	s.FunctionDefinition = fmt.Sprintf(`$$%s$$`, functionDefinition)
	return &s
}

func NewCreateForJavascriptFunctionRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	returns FunctionReturnsRequest,
	functionDefinition string,
) *CreateForJavascriptFunctionRequest {
	s := CreateForJavascriptFunctionRequest{}
	s.name = name
	s.Returns = returns
	s.FunctionDefinition = fmt.Sprintf(`$$%s$$`, functionDefinition)
	return &s
}
