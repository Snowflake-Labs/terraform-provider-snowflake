package generator

type ShowByIDFiltering struct {
	Kind            string
	Args            string
	IdentifierBased bool
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
				Kind:            "Like",
				Args:            "Pattern: String(id.Name())",
				IdentifierBased: false,
			})
		case ShowByIDInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Kind:            "In",
				Args:            "%[1]v: id.%[1]vId()",
				IdentifierBased: true,
			})
		case ShowByIDExtendedInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Kind:            "ExtendedIn",
				Args:            "In: In{%[1]v: id.%[1]vId()}",
				IdentifierBased: true,
			})
		}
	}
	return s
}
