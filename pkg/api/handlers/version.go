package handlers

import (
	"os"

	"github.com/orion101-ai/orion101/pkg/api"
	"github.com/orion101-ai/orion101/pkg/version"
	"sigs.k8s.io/yaml"
)

func GetVersion(req api.Context) error {
	return req.Write(getVersionResponse())
}

func getVersionResponse() map[string]string {
	values := make(map[string]string)
	versions := os.Getenv("ORION101_SERVER_VERSIONS")
	if versions != "" {
		if err := yaml.Unmarshal([]byte(versions), &values); err != nil {
			values["error"] = err.Error()
		}
	}
	values["orion101"] = version.Get().String()
	return values
}
