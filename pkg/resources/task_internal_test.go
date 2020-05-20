package resources

import "testing"

func Test_splitQualifiedName(t *testing.T) {
	type args struct {
		qualifiedName string
	}
	tests := []struct {
		name         string
		args         args
		wantName     string
		wantSchema   string
		wantDatabase string
	}{
		{
			name: "all-quoted",
			args: args{
				qualifiedName: "\"other\".\"schema\".\"task\"",
			},
			wantName:     "task",
			wantSchema:   "schema",
			wantDatabase: "other",
		},
		{
			name: "schema-unquoted",
			args: args{
				qualifiedName: "\"other\".SCHEMA.\"task\"",
			},
			wantName:     "task",
			wantSchema:   "SCHEMA",
			wantDatabase: "other",
		},
		{
			name: "name-unquoted",
			args: args{
				qualifiedName: "\"other\".\"schema\".TASK",
			},
			wantName:     "TASK",
			wantSchema:   "schema",
			wantDatabase: "other",
		},
		{
			name: "name-only-unquoted",
			args: args{
				qualifiedName: "TASK",
			},
			wantName:     "TASK",
			wantSchema:   "",
			wantDatabase: "",
		},
		{
			name: "name-only-quoted",
			args: args{
				qualifiedName: "\"task\"",
			},
			wantName:     "task",
			wantSchema:   "",
			wantDatabase: "",
		},
		{
			name: "empty",
			args: args{
				qualifiedName: "",
			},
			wantName:     "",
			wantSchema:   "",
			wantDatabase: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotSchema, gotDatabase := splitQualifiedName(tt.args.qualifiedName)
			if gotName != tt.wantName {
				t.Errorf("splitQualifiedName() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotSchema != tt.wantSchema {
				t.Errorf("splitQualifiedName() gotSchema = %v, want %v", gotSchema, tt.wantSchema)
			}
			if gotDatabase != tt.wantDatabase {
				t.Errorf("splitQualifiedName() gotDatabase = %v, want %v", gotDatabase, tt.wantDatabase)
			}
		})
	}
}
