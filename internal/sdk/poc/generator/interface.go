// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package generator

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// NameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	NameSingular string
	// Operations contains all operations for given interface
	Operations []*Operation
	// IdentifierKind keeps identifier of the underlying object (e.g. DatabaseObjectIdentifier)
	IdentifierKind string
}

func NewInterface(name string, nameSingular string, identifierKind string, operations ...*Operation) *Interface {
	return &Interface{
		Name:           name,
		NameSingular:   nameSingular,
		IdentifierKind: identifierKind,
		Operations:     operations,
	}
}

// NameLowerCased returns interface name starting with a lower case letter
func (i *Interface) NameLowerCased() string {
	return startingWithLowerCase(i.Name)
}
