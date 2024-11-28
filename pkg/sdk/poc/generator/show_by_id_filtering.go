package generator

import (
	"fmt"
)

type ShowByIDFiltering struct {
	Name string
	Kind string
	Args string
}

func (s *ShowByIDFiltering) String() string {
	return fmt.Sprintf("With%s(%s{%s})", s.Name, s.Kind, s.Args)
}

type ShowByIDFilter interface {
	String() string
}

func ShowByIDLikeFilteringFunc() ShowByIDFilter {
	return &ShowByIDFiltering{
		Name: "Like",
		Kind: "Like",
		Args: "Pattern: String(id.Name())",
	}
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
			})
		case ShowByIDInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Name: "In",
				Kind: "In",
				Args: fmt.Sprintf("%[1]v: id.%[1]vId()", s.ObjectInterface.ObjectIdentifierKind()),
			})
		case ShowByIDExtendedInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, ShowByIDFiltering{
				Name: "In",
				Kind: "ExtendedIn",
				Args: fmt.Sprintf("In: In{%[1]v: id.%[1]vId()}", s.ObjectInterface.ObjectIdentifierKind()),
			})
		}
	}
	return s
}
