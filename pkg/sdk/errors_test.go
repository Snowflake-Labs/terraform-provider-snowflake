package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSDKError(t *testing.T) {
	err := NewError("error one")
	err = WrapError(err,
		NewError("error two"),
		//errors.New("error three"),
		//errors.Join(
		//	errors.New("joined err 1"),
		//	NewError("joined err 2"),
		//),
		WrapError(
			NewError("root"),
			NewError("branch 1"),
			NewError("branch 2"),
		),
		NewError("error four"),
	)
	assert.NoError(t, err)
}
