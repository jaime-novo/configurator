package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/banknovo/configurator/core"
	"github.com/banknovo/configurator/export"
	"github.com/banknovo/configurator/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	exportMode    string
	blueprintFile string
	exportFormat  string
	exportOutput  string

	exportCmd = &cobra.Command{
		Use:   "export [flags]",
		Short: "Export the config",
		RunE:  runExport,
	}
)

func init() {
	exportCmd.Flags().StringVarP(&exportMode, "mode", "m", "",
		`Mode of export required
	flat: Keys are exported as-is in a flat structure
	hierarchical: Keys are exported in a hierarchical structure, keys broken down by separator
	blueprint: A blueprint file is taken which has keys are values, and real values are replaced in it`)

	exportCmd.Flags().StringVarP(&blueprintFile, "blueprint-file", "b", "", "Path to blueprint file (only required if exportMode is blueprint)")

	exportCmd.Flags().StringVarP(&format, "format", "f", "json",
		`Format of the export
	json`)

	exportCmd.Flags().StringVarP(&exportOutput, "output-file", "o", "", "Output file (default is standard output)")

	RootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	var err error

	env, err := getEnvironment()
	if err != nil {
		return err
	}

	c, err := getExporter(exportMode)
	if err != nil {
		return err
	}

	s, err := store.NewAWSPMStore()
	if err != nil {
		return err
	}

	allConfigs := getConfigs(env, s)
	configMap, err := c.Export(allConfigs)
	if err != nil {
		return err
	}

	file := os.Stdout
	if exportOutput != "" {
		if file, err = os.OpenFile(exportOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return errors.Wrap(err, "Failed to open output file for writing")
		}
		defer file.Close()
		defer file.Sync()
	}

	w := bufio.NewWriter(file)
	defer w.Flush()

	switch exportFormat {
	case "json":
	default:
		err = exportAsJSON(configMap, w)
	}

	if err != nil {
		return errors.Wrap(err, "Unable to export configs")
	}

	return nil
}

func exportAsJSON(configMap map[string]interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(configMap)
}

// getConfig fetches the configs from the store
func getConfigs(env string, s store.Store) []*core.Config {
	additionalConfig = append(additionalConfig, app)
	allConfigs := make([]*core.Config, 0)
	for _, config := range additionalConfig {
		key := "/" + env + "/" + config
		fmt.Printf("Fetching config for %s\n", key)
		configs, err := s.FetchAll(key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// cleansing - remove the env prefix from key names
		for _, config := range configs {
			config.Key = config.Key[len(env)+2:]
		}
		allConfigs = append(allConfigs, configs...)
	}
	return allConfigs
}

// getExporter returns the exporter based on mode
func getExporter(mode string) (export.Exporter, error) {
	var e export.Exporter
	switch mode {
	case "flat":
		e = &export.FlatExporter{}
	case "hierarchical":
		e = &export.HierarchicalExporter{
			Separator: "/",
		}
	case "blueprint":
		{
			data, err := getBlueprint()
			if err != nil {
				return nil, err
			}

			var blueprintMap map[string]interface{}
			err = json.Unmarshal(data, &blueprintMap)
			if err != nil {
				return nil, err
			}

			e = &export.BlueprintBasedExporter{
				Blueprint: blueprintMap,
			}
		}
	default:
		return nil, fmt.Errorf("Invalid export mode `%s`", mode)
	}
	return e, nil
}

func getBlueprint() ([]byte, error) {
	if blueprintFile == "" {
		return nil, fmt.Errorf("blueprint-file is required when export mode is blueprint")
	}
	return ioutil.ReadFile(blueprintFile)
}