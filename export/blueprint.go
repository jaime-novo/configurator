package export

import (
	"fmt"

	"github.com/banknovo/configurator/core"
)

// ensure BlueprintBasedExporter confirms to Exporter interface
var _ Exporter = &BlueprintBasedExporter{}

// BlueprintBasedExporter takes a blueprint of the required structure and
// aims to map the configs to it based on the lookup of keys as values
type BlueprintBasedExporter struct {
	Blueprint map[string]interface{}
}

// Export replaces the values in blueprint with their actual values from config based on key matching
func (e *BlueprintBasedExporter) Export(configs []*core.Config) (map[string]interface{}, error) {
	// use FlatExporter to get map of configs
	f := &FlatExporter{}
	configMap, err := f.Export(configs)
	if err != nil {
		return nil, err
	}

	// create a copy of blueprint and create the final version using configMap
	err = createMap(configMap, e.Blueprint)
	if err != nil {
		return nil, err
	}

	return e.Blueprint, nil
}

// createMap is a recursive function that iterates through the blueprint and
// replaces its values, which are key names, from the configMap
func createMap(configMap map[string]interface{}, blueprint map[string]interface{}) error {
	for key, val := range blueprint {
		v, ok := val.(map[string]interface{})
		if ok {
			createMap(configMap, v)
		} else {
			v, ok := val.(string)
			if !ok {
				return fmt.Errorf("Expected %s to have string value", val)
			}
			value := configMap[v]
			if value == nil {
				return fmt.Errorf("Config value %s missing", v)
			}
			blueprint[key] = value
		}
	}
	return nil
}
