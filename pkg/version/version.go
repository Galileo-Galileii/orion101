package version

import (
	"github.com/orion101-ai/nah/pkg/version"
)

var (
	Tag = "v0.0.0-dev"
)

func Get() version.Version {
	return version.NewVersion(Tag)
}
