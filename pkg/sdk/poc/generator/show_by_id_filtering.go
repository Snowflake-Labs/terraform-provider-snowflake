package generator

import (
	"fmt"
	"log"
)

type ShowByIDFilteringKind uint

const (
	// Enables filtering with: Like
	ShowByIDLikeFiltering ShowByIDFilteringKind = iota
	// Enables filtering with: In
	// Based on the identifier Kind
	ShowByIDInFiltering
	// Enables filtering with: In
	// Based on the identifier Kind
	ShowByIDExtendedInFiltering
	// Enables filtering with: Limit
	ShowByIDLimitFiltering
)

type ShowByIDFiltering interface {
	WithFiltering() string
}

type showByIDFilter struct {
	Name string
	Kind string
	Args string
}

func (s *showByIDFilter) WithFiltering() string {
	return fmt.Sprintf("With%s(%s{%s})", s.Name, s.Kind, s.Args)
}

var filteringMapping = map[ShowByIDFilteringKind]func(string) ShowByIDFiltering{
	ShowByIDLikeFiltering:       newShowByIDLikeFiltering,
	ShowByIDInFiltering:         newShowByIDInFiltering,
	ShowByIDExtendedInFiltering: newShowByIDExtendedInFiltering,
	ShowByIDLimitFiltering:      newShowByIDLimitFiltering,
}

func newShowByIDFiltering(name, kind, args string, identifierKind *string) ShowByIDFiltering {
	filter := &showByIDFilter{
		Name: name,
		Kind: kind,
		Args: args,
	}
	if identifierKind != nil {
		filter.Args = fmt.Sprintf(args, *identifierKind)
	}
	return filter
}

func newShowByIDLikeFiltering(string) ShowByIDFiltering {
	return newShowByIDFiltering("Like", "Like", "Pattern: String(id.Name())", nil)
}

func newShowByIDInFiltering(identifierKind string) ShowByIDFiltering {
	return newShowByIDFiltering("In", "In", "%[1]v: id.%[1]vId()", &identifierKind)
}

func newShowByIDExtendedInFiltering(identifierKind string) ShowByIDFiltering {
	return newShowByIDFiltering("In", "In", "In: In{%[1]v: id.%[1]vId()}", &identifierKind)
}

func newShowByIDLimitFiltering(string) ShowByIDFiltering {
	return newShowByIDFiltering("Limit", "LimitFrom", "Rows: Int(1)", nil)
}

func (s *Operation) withFiltering(filtering ...ShowByIDFilteringKind) *Operation {
	for _, filteringKind := range filtering {
		if filter, ok := filteringMapping[filteringKind]; ok {
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, filter(s.ObjectInterface.ObjectIdentifierKind()))
		} else {
			log.Println("No showByID filtering found for kind:", filteringKind)
		}
	}
	return s
}
