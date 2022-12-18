/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"examples/scrappy/internal/es"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Elasticsearch cluster configuration
const usernameFlagKey = "es_username"
const passwordFlagKey = "es_password"
const esURLFlagKey = "es_url"
const caCertFlagKey = "es_cacert"
const defaultCaCertPath = "./elasticsearch_ca.crt"

var esFlags = []string{usernameFlagKey, passwordFlagKey, esURLFlagKey, caCertFlagKey}

var (
	ErrMissingConfig = errors.New("missing ES config value")
	ErrMissingCACert = errors.New("failed to load CA certificate")
)

// esCmd represents the es command
var esCmd = &cobra.Command{
	Use:   "es",
	Short: "Elastic Search commands",
}

func init() {
	rootCmd.AddCommand(esCmd)

	flags := []struct {
		key          string
		defaultValue string
		usage        string
	}{
		{key: caCertFlagKey,
			defaultValue: defaultCaCertPath,
			usage:        "Elasticsearch cluster CA certificate"},
		{key: usernameFlagKey, usage: "Elasticsearch username"},
		{key: passwordFlagKey, usage: "Elasticsearch password"},
		{key: esURLFlagKey, usage: "Elasticsearch URL"},
	}

	for _, flag := range flags {
		esCmd.PersistentFlags().String(flag.key, flag.defaultValue, flag.usage)
		viper.BindPFlag(flag.key, esCmd.PersistentFlags().Lookup(flag.key))
	}
}

func esConfig() (*es.Config, error) {
	// Get config values from Viper
	values := map[string]string{}

	for _, flag := range esFlags {
		value := viper.GetString(flag)
		if value == "" {
			return nil, fmt.Errorf("%w: %s", ErrMissingConfig, flag)
		}
		values[flag] = value
	}

	// Load CA certificate from path
	caCert, err := os.ReadFile(values[caCertFlagKey])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMissingCACert, err)
	}

	config := es.Config{
		Username:  values[usernameFlagKey],
		Password:  values[passwordFlagKey],
		CACert:    caCert,
		Addresses: []string{values[esURLFlagKey]},
	}

	return &config, nil
}
