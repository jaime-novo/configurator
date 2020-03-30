package export

import (
	"strings"

	"github.com/banknovo/configurator/core"
)

// ensure HierarchicalExporter confirms to Exporter interface
var _ Exporter = &HierarchicalExporter{}

// HierarchicalExporter converts the configs into a hierarchical structure
type HierarchicalExporter struct {
	Separator string
}

// Export converts the configs into a nested map of key value pairs
func (e *HierarchicalExporter) Export(configs []*core.Config) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for _, config := range configs {
		keys := strings.Split(config.Key, e.Separator)
		value := getTypedValue(config.Value)
		addToMap(configMap, keys, value)
	}
	return configMap, nil
}

// addToMap is a recursive function that creates the hierarchical map of keys and values
func addToMap(configMap map[string]interface{}, keys []string, value interface{}) map[string]interface{} {
	if len(keys) == 1 {
		configMap[keys[0]] = value
		return configMap
	}
	currKey := keys[0]
	var newMap map[string]interface{}
	if configMap[currKey] != nil {
		newMap = configMap[currKey].(map[string]interface{})
	} else {
		newMap = make(map[string]interface{})
	}
	configMap[currKey] = addToMap(newMap, keys[1:], value)
	return configMap
}
