package snowflake

import (
	"testing"
)

func TestViewSelectStatementExtractor_Extract(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// {"basic", args{"create view foo as select * from bar;"}, "select * from bar;", false},
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
				t.Errorf("ViewSelectStatementExtractor.Extract() = '%v', want %v", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_remainingString(t *testing.T) {
	type fields struct {
		input string
		pos   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"basic1", fields{"foo", 0}, "foo"},
		{"basic2", fields{"foo", 1}, "oo"},
		{"empty", fields{"", 0}, ""},
		{"overflow", fields{"foo", 3}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			if got := e.remainingString(); got != tt.want {
				t.Errorf("ViewSelectStatementExtractor.remainingString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_nextN(t *testing.T) {
	type fields struct {
		input string
		pos   int
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"empty", fields{"", 0}, args{7}, ""},
		{"basic 1", fields{"foo", 0}, args{1}, "f"},
		{"basic 2", fields{"foo", 1}, args{1}, "o"},
		{"basic 3", fields{"foo", 0}, args{2}, "fo"},
		{"basic 4", fields{"foo", 1}, args{2}, "oo"},
		{"basic 5", fields{"foo", 2}, args{1}, "o"},
		{"overflow 1", fields{"foo", 3}, args{3}, ""},
		{"overflow 2", fields{"foo", 1}, args{4}, "oo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			if got := e.nextN(tt.args.n); got != tt.want {
				t.Errorf("ViewSelectStatementExtractor.nextN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeToken(t *testing.T) {
	type fields struct {
		input string
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
		{"basic - found", fields{"foo", 0}, args{"foo"}, 3},
		{"basic - not found", fields{"foo", 0}, args{"bar"}, 0},
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
