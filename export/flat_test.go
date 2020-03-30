package export

import (
	"testing"

	"github.com/banknovo/configurator/core"
	"github.com/stretchr/testify/require"
)

func TestFlatExport(t *testing.T) {
	configs := make([]*core.Config, 0)

	configs = append(configs, &core.Config{
		Key:   "Key1",
		Value: "value1",
	})
	configs = append(configs, &core.Config{
		Key:   "Key2",
		Value: "value2",
	})
	configs = append(configs, &core.Config{
		Key:   "Key3",
		Value: "1",
	})
	configs = append(configs, &core.Config{
		Key:   "Key4",
		Value: "true",
	})
	configs = append(configs, &core.Config{
		Key:   "Key5",
		Value: "false",
	})

	var e Exporter = &FlatExporter{}
	configMap, err := e.Export(configs)

	require.Empty(t, err, "got error while export")
	require.Equal(t, "value1", configMap["Key1"])
	require.Equal(t, "value2", configMap["Key2"])
	require.Equal(t, int64(1), configMap["Key3"])
	require.Equal(t, true, configMap["Key4"])
	require.Equal(t, false, configMap["Key5"])
}
