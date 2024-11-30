package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

func output(format string, obj any) (bool, error) {
	if format == "" || format == "table" {
		return false, nil
	}
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return true, enc.Encode(obj)
	case "yaml":
		data, err := yaml.Marshal(obj)
		if err != nil {
			return false, err
		}
		_, err = os.Stdout.Write(data)
		return true, err
	default:
		return false, fmt.Errorf("unsupported output format: %s", format)
	}
}
