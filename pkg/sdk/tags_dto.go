package sdk

var (
	_ optionsProvider[createTagOptions] = new(CreateTagRequest)
	_ optionsProvider[showTagOptions]   = new(ShowTagRequest)
	_ optionsProvider[dropTagOptions]   = new(DropTagRequest)
	_ optionsProvider[undropTagOptions] = new(UndropTagRequest)
)

type CreateTagRequest struct {
	orReplace   bool
	ifNotExists bool

	name SchemaObjectIdentifier // required

	// One of
	comment       *string
	allowedValues *AllowedValues
}

type AlterTagRequest struct {
	name SchemaObjectIdentifier // required

	// One of
	add    *TagAdd
	drop   *TagDrop
	set    *TagSet
	unset  *TagUnset
	rename *TagRename
}

type TagSetRequest struct {
	maskingPolicies []string
	force           *bool
	comment         *string
}

type TagUnsetRequest struct {
	maskingPolicies []string
	allowedValues   *bool
	comment         *bool
}

type ShowTagRequest struct {
	like *Like
	in   *In
}

type DropTagRequest struct {
	ifExists bool

	name SchemaObjectIdentifier // required
}

type UndropTagRequest struct {
	name SchemaObjectIdentifier // required
}
