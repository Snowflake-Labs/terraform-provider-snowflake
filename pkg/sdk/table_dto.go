package sdk

type TableCreateDto struct {
	OrReplace             bool
	Scope                 *TableScope
	Kind                  *TableKind
	IfNotExists           bool
	Name                  SchemaObjectIdentifier
	ClusterBy             []string
	EnableSchemaEvolution *bool
	tageFileFormat        []StageFileFormat
	tageCopyOptions       []StageCopyOptions
}
