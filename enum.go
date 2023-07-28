

//DSL
columnConsraintType := b.QueryEnum[string]("ColumnConstraintType")
		.With("Unique", "UNIQUE")
		.With("PrimaryKey", "PRIMARY KEY")
		.With("ForeignKey", "FOREIGN KEY")

//wynik
type ColumnConstraintType string

const (
	ColumnConstraintTypeUnique     ColumnConstraintType = "UNIQUE"
	ColumnConstraintTypePrimaryKey ColumnConstraintType = "PRIMARY KEY"
	ColumnConstraintTypeForeignKey ColumnConstraintType = "FOREIGN KEY"
)


