package sdk

func NewCreateTagRequest(name AccountObjectIdentifier) *CreateTagRequest {
	s := CreateTagRequest{}
	s.name = name
	return &s
}

func (s *CreateTagRequest) WithOrReplace(orReplace bool) *CreateTagRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTagRequest) WithIfNotExists(ifNotExists bool) *CreateTagRequest {
	s.ifNotExists = ifNotExists
	return s
}

func (s *CreateTagRequest) WithComment(comment *string) *CreateTagRequest {
	s.comment = comment
	return s
}

func (s *CreateTagRequest) WithAllowedValues(values []string) *CreateTagRequest {
	if len(values) > 0 {
		s.allowedValues = &AllowedValues{
			Values: make([]AllowedValue, 0, len(values)),
		}
		for _, value := range values {
			s.allowedValues.Values = append(s.allowedValues.Values, AllowedValue{
				Value: value,
			})
		}
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

func NewDropTagRequest(name AccountObjectIdentifier) *DropTagRequest {
	s := DropTagRequest{}
	s.name = name
	return &s
}

func (s *DropTagRequest) WithIfNotExists(ifNotExists bool) *DropTagRequest {
	s.ifNotExists = ifNotExists
	return s
}

func NewUndropTagRequest(name AccountObjectIdentifier) *UndropTagRequest {
	s := UndropTagRequest{}
	s.name = name
	return &s
}
