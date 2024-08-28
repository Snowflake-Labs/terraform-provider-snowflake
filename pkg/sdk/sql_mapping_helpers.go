package sdk

import (
	"database/sql"
	"log"
	"strconv"
)

func handleNullableBoolString(nullableBoolString sql.NullString, field *bool) {
	if nullableBoolString.Valid && nullableBoolString.String != "" && nullableBoolString.String != "null" {
		parsed, err := strconv.ParseBool(nullableBoolString.String)
		if err != nil {
			// TODO [SNOW-1561641]: address during handling the issue
			log.Printf("[DEBUG] Could not parse text boolean value %v", nullableBoolString.String)
		} else {
			*field = parsed
		}
	}
}
