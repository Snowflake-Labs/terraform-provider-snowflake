// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"fmt"
	"regexp"
	"strings"
)

// Generic Errors.
var (
	ErrNoRowInRS = "sql: no rows in result set"
)

func IsResourceNotExistOrNotAuthorized(errorString string, resourceType string) bool {
	regexStr := fmt.Sprintf("SQL compilation error:%s '.*' does not exist or not authorized", resourceType)
	userNotExistOrNotAuthorizedRegEx, _ := regexp.Compile(regexStr)
	return userNotExistOrNotAuthorizedRegEx.MatchString(strings.ReplaceAll(errorString, "\n", ""))
}
