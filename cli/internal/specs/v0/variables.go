package specs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
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
// Optimization: Direct reflection traversal (getVariableValue) instead of expensive
// json.Marshal, json.Unmarshal, and third-party reflection library (funk.Get).
func ReplaceVariables(src string, variables Variables, shouldReplaceLocalhost bool) (string, error) {
	var lastErr error
	varsValue := reflect.ValueOf(variables)
	result := reVariables.ReplaceAllStringFunc(src, func(s string) string {
		variablePath := s[2:]
		parts := strings.Split(variablePath, ".")
		res := getVariableValue(varsValue, parts)
		if res == nil {
			lastErr = fmt.Errorf("variable %s not found", variablePath)
			return s
		}
		resString, ok := res.(string)
		if !ok {
			lastErr = fmt.Errorf("variable %s is not a string", variablePath)
			return s
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

func getVariableValue(val reflect.Value, parts []string) any {
	if len(parts) == 0 {
		if !val.IsValid() {
			return nil
		}
		return val.Interface()
	}

	if !val.IsValid() {
		return nil
	}

	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	key := parts[0]
	switch val.Kind() {
	case reflect.Map:
		for _, mapKey := range val.MapKeys() {
			if mapKey.String() == key {
				return getVariableValue(val.MapIndex(mapKey), parts[1:])
			}
		}
		return nil
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			tagName := strings.Split(tag, ",")[0]
			if tagName == key || strings.EqualFold(f.Name, key) {
				return getVariableValue(val.Field(i), parts[1:])
			}
		}
		return nil
	default:
		return nil
	}
}
