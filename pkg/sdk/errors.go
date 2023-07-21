package sdk

import (
	"errors"
	"log"
	"strings"
)

var (
	// go-snowflake errors.
	ErrObjectNotExistOrAuthorized = errors.New("object does not exist or not authorized")
	ErrAccountIsEmpty             = errors.New("account is empty")

	// snowflake-sdk errors.
	ErrInvalidObjectIdentifier = errors.New("invalid object identifier")
	ErrNilOptions              = errors.New("options cannot be nil")
)

func decodeDriverError(err error) error {
	if err == nil {
		return nil
	}
	log.Printf("[DEBUG] err: %v\n", err)
	m := map[string]error{
		"does not exist or not authorized": ErrObjectNotExistOrAuthorized,
		"account is empty":                 ErrAccountIsEmpty,
	}
	for k, v := range m {
		if strings.Contains(err.Error(), k) {
			return v
		}
	}

	return err
}
