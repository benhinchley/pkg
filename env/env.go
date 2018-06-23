// Package env provides methods for getting environment variables from
// appengine configuration files
package env

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// AppFilePath is the path to the app.yaml file
var AppFilePath = "app.yaml"

// Get returns the value for the passed environment variable key
func Get(k string) string {
	env, err := parseEnvFromYAML(AppFilePath)
	if err != nil {
		return os.Getenv(k)
	}
	return env[k]
}

// GetWithFallback returns the value for the passed environment variable key
// returning the fallback value if nothing is found
func GetWithFallback(k, fallback string) string {
	v := Get(k)
	if v == "" {
		return fallback
	}
	return v
}

type envFile struct {
	Includes     []string          `yaml:"include"`
	EnvVariables map[string]string `yaml:"env_variables"`
}

func parseEnvFromYAML(filename string) (map[string]string, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ef envFile
	if err := yaml.NewDecoder(f).Decode(&ef); err != nil {
		return nil, err
	}

	if len(ef.Includes) > 0 {
		for _, inc := range ef.Includes {
			m, err := parseEnvFromYAML(filepath.Join(filepath.Dir(filename), inc))
			if err != nil {
				continue
			}
			for k, v := range m {
				ef.EnvVariables[k] = v
			}
		}
	}

	return ef.EnvVariables, nil
}
