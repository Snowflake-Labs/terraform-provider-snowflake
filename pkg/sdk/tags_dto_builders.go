package sdk

func NewCreateTagRequest(name SchemaObjectIdentifier) *CreateTagRequest {
	s := CreateTagRequest{}
	s.name = name
	return &s
}

func (s *CreateTagRequest) WithOrReplace(orReplace bool) *CreateTagRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTagRequest) WithIfExists(ifExists bool) *CreateTagRequest {
	s.ifNotExists = ifExists
	return s
}

func (s *CreateTagRequest) WithComment(comment *string) *CreateTagRequest {
	s.comment = comment
	return s
}

func (s *CreateTagRequest) WithAllowedValues(values []string) *CreateTagRequest {
	if len(values) > 0 {
		s.allowedValues = createAllowedValues(values)
	}
	return s
}

func createAllowedValues(values []string) *AllowedValues {
	items := make([]AllowedValue, 0, len(values))
	for _, value := range values {
		items = append(items, AllowedValue{
			Value: value,
		})
	}
	return &AllowedValues{
		Values: items,
	}
}

func NewAlterTagRequest(name SchemaObjectIdentifier) *AlterTagRequest {
	s := AlterTagRequest{}
	s.name = name
	return &s
}

func (s *AlterTagRequest) WithAdd(values []string) *AlterTagRequest {
	if len(values) > 0 {
		s.add = &TagAdd{createAllowedValues(values)}
	}
	return s
}

func (s *AlterTagRequest) WithDrop(values []string) *AlterTagRequest {
	if len(values) > 0 {
		s.drop = &TagDrop{createAllowedValues(values)}
	}
	return s
}

func NewTagSetRequest() *TagSetRequest {
	return &TagSetRequest{}
}

func (s *TagSetRequest) WithMaskingPolicies(maskingPolicies []string) *TagSetRequest {
	s.maskingPolicies = maskingPolicies
	return s
}

func (s *TagSetRequest) WithForce(force bool) *TagSetRequest {
	s.force = Bool(force)
	return s
}

func (s *TagSetRequest) WithComment(comment string) *TagSetRequest {
	s.comment = String(comment)
	return s
}

func createTagMaskingPolicies(policies []string) []TagMaskingPolicy {
	items := make([]TagMaskingPolicy, 0, len(policies))
	for _, policy := range policies {
		items = append(items, TagMaskingPolicy{
			Name: policy,
		})
	}
	return items
}

func (s *AlterTagRequest) WithSet(request *TagSetRequest) *AlterTagRequest {
	set := &TagSet{
		Comment: request.comment,
	}
	if len(request.maskingPolicies) > 0 {
		set.MaskingPolicies = &TagSetMaskingPolicies{
			MaskingPolicies: createTagMaskingPolicies(request.maskingPolicies),
			Force:           request.force,
		}
	}
	s.set = set
	return s
}

func NewTagUnsetRequest() *TagUnsetRequest {
	return &TagUnsetRequest{}
}

func (s *TagUnsetRequest) WithMaskingPolicies(maskingPolicies []string) *TagUnsetRequest {
	s.maskingPolicies = maskingPolicies
	return s
}

func (s *TagUnsetRequest) WithAllowedValues(allowedValues bool) *TagUnsetRequest {
	s.allowedValues = Bool(allowedValues)
	return s
}

func (s *TagUnsetRequest) WithComment(comment bool) *TagUnsetRequest {
	s.comment = Bool(comment)
	return s
}

func (s *AlterTagRequest) WithUnset(request *TagUnsetRequest) *AlterTagRequest {
	unset := &TagUnset{
		AllowedValues: request.allowedValues,
		Comment:       request.comment,
	}
	if len(request.maskingPolicies) > 0 {
		unset.MaskingPolicies = &TagUnsetMaskingPolicies{
			MaskingPolicies: createTagMaskingPolicies(request.maskingPolicies),
		}
	}
	s.unset = unset
	return s
}

func (s *AlterTagRequest) WithRename(name SchemaObjectIdentifier) *AlterTagRequest {
	s.rename = &TagRename{
		Name: name,
	}
	return s
}

func NewShowTagRequest() *ShowTagRequest {
	return &ShowTagRequest{}
}

func (s *ShowTagRequest) WithLike(pattern string) *ShowTagRequest {
	s.like = &Like{
		Pattern: String(pattern),
	}
	return s
}

func (s *ShowTagRequest) WithIn(in *In) *ShowTagRequest {
	s.in = in
	return s
}

func NewDropTagRequest(name SchemaObjectIdentifier) *DropTagRequest {
	s := DropTagRequest{}
	s.name = name
	return &s
}

func (s *DropTagRequest) WithIfNotExists(ifNotExists bool) *DropTagRequest {
	s.ifExists = ifNotExists
	return s
}

func NewUndropTagRequest(name SchemaObjectIdentifier) *UndropTagRequest {
	s := UndropTagRequest{}
	s.name = name
	return &s
}
