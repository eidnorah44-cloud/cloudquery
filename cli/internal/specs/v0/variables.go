package specs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/thoas/go-funk"
)

type PluginVariables struct {
	Connection string `json:"connection"`
}

type Variables struct {
	Plugins map[string]PluginVariables `json:"plugins"`
}

var reVariables = regexp.MustCompile(`@@(plugins\.[a-zA-Z0-9_\.-]+)`)

// ReplaceVariables replaces variables starting with @@ in the src string
// with the values from the values from variables by dot notation.
// Example: @@plugins.aws.connection will be replaced with the value of variables.Plugins["aws"].Connection
func ReplaceVariables(src string, variables Variables, shouldReplaceLocalhost bool) (string, error) {
	// Performance optimization: fast-path check to avoid regex & lookup overhead when no variables exist in src
	if !strings.Contains(src, "@@") {
		return src, nil
	}

	var lastErr error
	// Direct struct lookup instead of converting struct to JSON and using reflection map lookups (go-funk)
	result := reVariables.ReplaceAllStringFunc(src, func(s string) string {
		variablePath := s[2:] // strip @@ prefix

		var resString string
		var found bool

		// Handle plugins.<plugin_name>.connection lookups
		if strings.HasPrefix(variablePath, "plugins.") && strings.HasSuffix(variablePath, ".connection") && len(variablePath) >= len("plugins.")+len(".connection") {
			pluginName := variablePath[len("plugins.") : len(variablePath)-len(".connection")]
			if plug, ok := variables.Plugins[pluginName]; ok {
				resString = plug.Connection
				found = true
			}
		}

		if !found {
			// Fallback using json serialization + funk.Get in case new fields/structures are added to Variables in the future
			bytes, err := json.Marshal(variables)
			if err != nil {
				lastErr = err
				return s
			}
			variablesMap := make(map[string]any)
			if err := json.Unmarshal(bytes, &variablesMap); err != nil {
				lastErr = err
				return s
			}
			res := funk.Get(variablesMap, variablePath)
			if res == nil {
				lastErr = fmt.Errorf("variable %s not found", variablePath)
				return s
			}
			var ok bool
			resString, ok = res.(string)
			if !ok {
				lastErr = fmt.Errorf("variable %s is not a string", variablePath)
				return s
			}
		}

		// Edge case: if the plugin whose spec's variables are being replaced is a docker plugin,
		// it won't be able to connect to localhost, so we replace localhost with host.docker.internal
		if strings.HasPrefix(variablePath, "plugins.") && strings.HasSuffix(variablePath, ".connection") && shouldReplaceLocalhost {
			for _, needle := range []string{"localhost", "0.0.0.0", "127.0.0.1"} {
				resString = strings.ReplaceAll(resString, needle, "host.docker.internal")
			}
		}

		// make safe for replacement into JSON string
		v, err := json.Marshal(resString)
		if err != nil {
			lastErr = err
			return s
		}
		resString = string(v[1 : len(v)-1])
		return resString
	})
	return result, lastErr
}
