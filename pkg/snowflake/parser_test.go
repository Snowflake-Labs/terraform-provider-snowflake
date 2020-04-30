package snowflake

import (
	"fmt"
	"testing"
)

func TestViewSelectStatementExtractor_Extract(t *testing.T) {
	basic := "create view foo as select * from bar;"
	caps := "CREATE VIEW FOO AS SELECT * FROM BAR;"
	parens := "create view foo as (select * from bar);"
	multiline := `
create view foo as
select *
from bar;`

	multilineComment := `
create view foo as
-- comment
select *
from bar;`

	secure := "create secure view foo as select * from bar;"
	replace := "create or replace view foo as select * from bar;"
	recursive := "create recursive view foo as select * from bar;"
	ine := "create view if not exists foo as select * from bar;"

	comment := `create view foo comment='asdf' as select * from bar;`
	commentEscape := `create view foo comment='asdf\'s are fun' as select * from bar;`
	identifier := `create view "foo"."bar"."bam" comment='asdf\'s are fun' as select * from bar;`

	full := `CREATE SECURE VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" COMMENT = 'Terraform test resource' AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`

	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{basic}, "select * from bar;", false},
		{"caps", args{caps}, "SELECT * FROM BAR;", false},
		{"parens", args{parens}, "(select * from bar);", false},
		{"multiline", args{multiline}, "select *\nfrom bar;", false},
		{"multilineComment", args{multilineComment}, "-- comment\nselect *\nfrom bar;", false},
		{"secure", args{secure}, "select * from bar;", false},
		{"replace", args{replace}, "select * from bar;", false},
		{"recursive", args{recursive}, "select * from bar;", false},
		{"ine", args{ine}, "select * from bar;", false},
		{"comment", args{comment}, "select * from bar;", false},
		{"commentEscape", args{commentEscape}, "select * from bar;", false},
		{"identifier", args{identifier}, "select * from bar;", false},
		{"full", args{full}, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewViewSelectStatementExtractor(tt.args.input)
			got, err := e.Extract()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewSelectStatementExtractor.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ViewSelectStatementExtractor.Extract() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeToken(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	type args struct {
		t string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		posAfter int
	}{
		{"basic - found", fields{[]rune("foo"), 0}, args{"foo"}, 3},
		{"basic - not found", fields{[]rune("foo"), 0}, args{"bar"}, 0},
		{"basic - not found", fields{[]rune("fob"), 0}, args{"foo"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeToken(tt.args.t)

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeSpace(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	tests := []struct {
		name     string
		fields   fields
		posAfter int
	}{
		{"simple", fields{[]rune("   foo"), 0}, 3},
		{"empty", fields{[]rune(""), 0}, 0},
		{"middle", fields{[]rune("foo \t\n bar"), 3}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeSpace()

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeComment(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	tests := []struct {
		name     string
		fields   fields
		posAfter int
	}{
		{"basic", fields{[]rune("comment='foo'"), 0}, 13},
		{"escaped", fields{[]rune(`comment='fo\'o'`), 0}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeComment()

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}
