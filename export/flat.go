package export

import (
	"github.com/banknovo/configurator/core"
)

// ensure FlatExporter confirms to Exporter interface
var _ Exporter = &FlatExporter{}

// FlatExporter converts the configs into flat key-value pairs
type FlatExporter struct{}

// Export converts the configs into a flat map with key value pairs
func (e *FlatExporter) Export(configs []*core.Config) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for _, config := range configs {
		configMap[config.Key] = getTypedValue(config.Value)
	}
	return configMap, nil
}
