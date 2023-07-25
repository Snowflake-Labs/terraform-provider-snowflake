package sdk

type AlertCreateOptions1 struct {
	create      bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	alert       bool                   `ddl:"static" sql:"ALERT"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`

	Warehouse string   `ddl:"parameter,double_quotes" sql:"WAREHOUSE"`
	Schedule  Schedule `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
}

type Schedule struct {
	// one of
	minute *Minute
	cron   *Cron
}

type Minute struct {
	N      int  `ddl:"keyword"`
	minute bool `ddl:"static" sql:"MINUTE"`
}

type Cron struct {
	usingCron bool   `ddl:"static" sql:"USING CRON"`
	expr      string `ddl:"keyword, no_quotes"`
	timezone  bool   `ddl:"keyword, no_quotes"`
}
