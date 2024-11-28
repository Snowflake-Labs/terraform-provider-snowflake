package generator

type ShowByIDFilteringType string

const (
	SimpleFiltering          ShowByIDFilteringType = "Simple"
	IdentifierBasedFiltering ShowByIDFilteringType = "IdentifierBased"
)

type ShowByIDFiltering struct {
	Name string
	Kind string
	Args string
	Type ShowByIDFilteringType
}

type ShowByIDFilteringKind uint

const (
	// Enables filtering with: Like
	ShowByIDLikeFiltering ShowByIDFilteringKind = iota
	// Enables filtering with: In
	// Based on the identifier Kind
	ShowByIDInFiltering
	// Enables filtering with: ExtendedIn
	// Based on the identifier Kind
	ShowByIDExtendedInFiltering
)

func (s *Operation) withFiltering(filtering ...ShowByIDFilteringKind) *Operation {
	for _, f := range filtering {
		switch f {
		case ShowByIDLikeFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Name: "Like",
				Kind: "Like",
				Args: "Pattern: String(id.Name())",
				Type: SimpleFiltering,
			})
		case ShowByIDInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Name: "In",
				Kind: "In",
				Args: "%[1]v: id.%[1]vId()",
				Type: IdentifierBasedFiltering,
			})
		case ShowByIDExtendedInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Name: "In",
				Kind: "ExtendedIn",
				Args: "In: In{%[1]v: id.%[1]vId()}",
				Type: IdentifierBasedFiltering,
			})
		}
	}
	return s
}
