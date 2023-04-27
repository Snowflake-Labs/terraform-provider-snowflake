package sdk

import (
	"errors"
	"strings"
)

var ErrObjectNotExistOrAuthorized = errors.New("object does not exist or not authorized")

func decodeError(err error) error {
	if err == nil {
		return nil
	}
	m := map[string]error{
		"does not exist or not authorized": ErrObjectNotExistOrAuthorized,
	}
	for k, v := range m {
		if strings.Contains(err.Error(), k) {
			return v
		}
	}

	return err
}
