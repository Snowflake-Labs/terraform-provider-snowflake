package builder

func foo() {
	create := QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf().
		Field().
		End().
}
