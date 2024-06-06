package sdk

func NewCreateCortexSearchServiceRequest(
	name SchemaObjectIdentifier,
	on string,
	warehouse AccountObjectIdentifier,
	targetLag string,
	query string,
) *CreateCortexSearchServiceRequest {
	s := CreateCortexSearchServiceRequest{}
	s.name = name
	s.on = on
	s.warehouse = warehouse
	s.targetLag = targetLag
	s.query = query
	return &s
}

func (s *CreateCortexSearchServiceRequest) WithOrReplace(orReplace bool) *CreateCortexSearchServiceRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateCortexSearchServiceRequest) WithIfNotExists(ifNotExists bool) *CreateCortexSearchServiceRequest {
	s.ifNotExists = ifNotExists
	return s
}

func (s *CreateCortexSearchServiceRequest) WithAttributes(attributes []string) *CreateCortexSearchServiceRequest {
	s.attributes = attributes
	return s
}

func (s *CreateCortexSearchServiceRequest) WithComment(comment *string) *CreateCortexSearchServiceRequest {
	s.comment = comment
	return s
}

func NewAlterCortexSearchServiceRequest(
	name SchemaObjectIdentifier,
) *AlterCortexSearchServiceRequest {
	s := AlterCortexSearchServiceRequest{}
	s.name = name
	return &s
}

func (s *AlterCortexSearchServiceRequest) WithSet(set *CortexSearchServiceSetRequest) *AlterCortexSearchServiceRequest {
	s.set = set
	return s
}

func NewCortexSearchServiceSetRequest() *CortexSearchServiceSetRequest {
	return &CortexSearchServiceSetRequest{}
}

func (s *CortexSearchServiceSetRequest) WithTargetLag(targetLag string) *CortexSearchServiceSetRequest {
	s.targetLag = &targetLag
	return s
}

func (s *CortexSearchServiceSetRequest) WithWarehouse(warehouse AccountObjectIdentifier) *CortexSearchServiceSetRequest {
	s.warehouse = &warehouse
	return s
}

func NewDropCortexSearchServiceRequest(
	name SchemaObjectIdentifier,
) *DropCortexSearchServiceRequest {
	s := DropCortexSearchServiceRequest{}
	s.name = name
	return &s
}

func (s *DropCortexSearchServiceRequest) WithIfExists(ifExists bool) *DropCortexSearchServiceRequest {
	s.IfExists = &ifExists
	return s
}

func NewDescribeCortexSearchServiceRequest(
	name SchemaObjectIdentifier,
) *DescribeCortexSearchServiceRequest {
	s := DescribeCortexSearchServiceRequest{}
	s.name = name
	return &s
}

func NewShowCortexSearchServiceRequest() *ShowCortexSearchServiceRequest {
	return &ShowCortexSearchServiceRequest{}
}

func (s *ShowCortexSearchServiceRequest) WithLike(like *Like) *ShowCortexSearchServiceRequest {
	s.like = like
	return s
}

func (s *ShowCortexSearchServiceRequest) WithIn(in *In) *ShowCortexSearchServiceRequest {
	s.in = in
	return s
}

func (s *ShowCortexSearchServiceRequest) WithStartsWith(startsWith *string) *ShowCortexSearchServiceRequest {
	s.startsWith = startsWith
	return s
}

func (s *ShowCortexSearchServiceRequest) WithLimit(limit *LimitFrom) *ShowCortexSearchServiceRequest {
	s.limit = limit
	return s
}
