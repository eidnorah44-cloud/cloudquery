package specs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	// Optimization: Skip regex and allocations if no variable placeholders are present
	if !strings.Contains(src, "@@") {
		return src, nil
	}

	var lastErr error
	result := reVariables.ReplaceAllStringFunc(src, func(s string) string {
		variablePath := s[2:]
		resString, err := getPluginVariable(variables, variablePath)
		if err != nil {
			lastErr = err
			return s
		}

		// Edge case: if the plugin whose spec's variables are being replaced is a docker plugin,
		// it won't be able to connect to localhost, so we replace localhost with host.docker.internal
		if strings.HasPrefix(variablePath, "plugins.") && strings.HasSuffix(variablePath, ".connection") && shouldReplaceLocalhost {
			for _, needle := range []string{"localhost", "0.0.0.0", "127.0.0.1"} {
				resString = strings.ReplaceAll(resString, needle, "host.docker.internal")
			}
		}

		// Optimization: make safe for replacement into JSON string using strconv.Quote instead of json.Marshal
		v := strconv.Quote(resString)
		return v[1 : len(v)-1]
	})
	return result, lastErr
}

// Optimization: Direct lookup avoids JSON marshal/unmarshal and reflection overhead from go-funk
func getPluginVariable(variables Variables, variablePath string) (string, error) {
	parts := strings.Split(variablePath, ".")
	if len(parts) != 3 || parts[0] != "plugins" {
		if len(parts) == 2 && parts[0] == "plugins" {
			if _, ok := variables.Plugins[parts[1]]; ok {
				return "", fmt.Errorf("variable %s is not a string", variablePath)
			}
		}
		return "", fmt.Errorf("variable %s not found", variablePath)
	}

	pluginName := parts[1]
	prop := parts[2]

	pluginVars, ok := variables.Plugins[pluginName]
	if !ok {
		return "", fmt.Errorf("variable %s not found", variablePath)
	}

	switch prop {
	case "connection":
		return pluginVars.Connection, nil
	default:
		return "", fmt.Errorf("variable %s not found", variablePath)
	}
}
