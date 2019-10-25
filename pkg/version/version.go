package version

import (
	"github.com/chanzuckerberg/go-misc/ver"
)

var (
	Version = "undefined"
	GitSha  = "undefined"
	Release = "false"
	Dirty   = "true"
)

func VersionString() (string, error) {
	return ver.VersionString(Version, GitSha, Release, Dirty)
}
