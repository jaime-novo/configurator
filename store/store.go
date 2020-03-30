package store

import "github.com/banknovo/configurator/core"

// Store represents storage of parameters
type Store interface {
	// Fetches all the values which start with given path
	FetchAll(path string) ([]*core.Config, error)
}
