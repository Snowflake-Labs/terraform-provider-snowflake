package generator

type ShowByIDFiltering struct {
	Kind string
	Args string
}

type ShowByIDFilteringKind uint

const (
	// Enables filtering with: Like
	ShowByIDLikeFiltering ShowByIDFilteringKind = iota
	// Enables filtering with: In
	ShowByIDInFiltering
	// Enables filtering with: ExtendedIn
	ShowByIDExtendedInFiltering
)

func (s *Operation) withFiltering(filtering ...ShowByIDFilteringKind) *Operation {
	for _, f := range filtering {
		switch f {
		case ShowByIDLikeFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Kind: "Like",
				Args: "Pattern: String(id.Name())",
			})
		case ShowByIDInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Kind: "In",
				Args: "%[1]v: id.%[1]vId()",
			})
		case ShowByIDExtendedInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Kind: "ExtendedIn",
				Args: "In: In{%[1]v: id.%[1]vId()}",
			})
		}
	}
	return s
}
