package provider

import (
	"context"
	"regexp"
)

type emptyDescriptions struct {
}

func (emptyDescriptions) Description(ctx context.Context) string {
	return ""
}

func (emptyDescriptions) MarkdownDescription(ctx context.Context) string {
	return ""
}

var errorNotFoundRegexp = regexp.MustCompile("NOT_FOUND|does not exist")

func isNotFoundError(err error) bool {
	return errorNotFoundRegexp.MatchString(err.Error())
}
