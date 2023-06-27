package convert

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/banknovo/configurator/config"
)

// Converter converts the config values into desired format
type Converter interface {
	// Convert converts the array of configs into a map
	Convert(configs []*config.Config) (map[string]interface{}, error)
}

func parseJson(jsonValue string) (interface{}, error) {

	if !json.Valid([]byte(jsonValue)) {
		return nil, fmt.Errorf("The given string is not a valid JSON")
	}

	var parsedJson interface{}

	if err := json.Unmarshal([]byte(jsonValue), &parsedJson); err != nil {
		return nil, err
	}

	return parsedJson, nil

}

// getTypedValue attempts to convert rawValue to its data type
func getTypedValue(rawValue string) interface{} {
	var parsedValue interface{}
	var err error

	// Initially, validate if the JSON prefix is set
	jsonValue, isJson := strings.CutPrefix(rawValue, "json::")

	if isJson == true {
		parsedValue, err := parseJson(jsonValue)
		if err == nil {
			return parsedValue
		} else {
			fmt.Println(err)
		}
	}

	// check if it is a bool
	if rawValue == "true" {
		return true
	}
	if rawValue == "false" {
		return false
	}

	// check if it is an int, if so convert to int64
	matched, err := regexp.Match("^(0|[1-9][0-9]*)$", []byte(rawValue))
	if err == nil && matched {
		parsedValue, err = strconv.ParseInt(rawValue, 10, 64)
		if err == nil {
			return parsedValue
		}
	}

	// for all others return rawValue as string
	return rawValue
}
