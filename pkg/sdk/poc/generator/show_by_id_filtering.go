package generator

import (
	"fmt"
	"log"
)

type ShowByIDFiltering interface {
	WithFiltering() string
}

type ShowByIDFilter struct {
	Name string
	Kind string
	Args string
}

func (s *ShowByIDFilter) WithFiltering() string {
	return fmt.Sprintf("With%s(%s{%s})", s.Name, s.Kind, s.Args)
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

func NewShowByIDFiltering(name, kind, args string, identifierKind *string) ShowByIDFiltering {
	filter := &ShowByIDFilter{
		Name: name,
		Kind: kind,
		Args: args,
	}
	if identifierKind != nil {
		filter.Args = fmt.Sprintf(args, *identifierKind)
	}
	return filter
}

func NewShowByIDLikeFiltering(string) ShowByIDFiltering {
	return NewShowByIDFiltering("Like", "Like", "Pattern: String(id.Name())", nil)
}

func NewShowByIDInFiltering(identifierKind string) ShowByIDFiltering {
	return NewShowByIDFiltering("In", "In", "%[1]v: id.%[1]vId()", &identifierKind)
}

func NewShowByIDExtendedInFiltering(identifierKind string) ShowByIDFiltering {
	return NewShowByIDFiltering("In", "ExtendedIn", "In: In{%[1]v: id.%[1]vId()}", &identifierKind)
}

var filteringMap = map[ShowByIDFilteringKind]func(string) ShowByIDFiltering{
	ShowByIDLikeFiltering:       NewShowByIDLikeFiltering,
	ShowByIDInFiltering:         NewShowByIDInFiltering,
	ShowByIDExtendedInFiltering: NewShowByIDExtendedInFiltering,
}

func (s *Operation) withFiltering(filtering ...ShowByIDFilteringKind) *Operation {
	for _, filteringKind := range filtering {
		if filter, ok := filteringMap[filteringKind]; ok {
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, filter(s.ObjectInterface.ObjectIdentifierKind()))
		} else {
			log.Println("No showByID filtering found for kind:", filteringKind)
		}
	}
	return s
}
