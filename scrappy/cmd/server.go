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
	"fmt"
	"net"
	"strconv"
	"time"

	"examples/scrappy/internal/es"
	"examples/scrappy/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Default port to start listening on.
const defaultPort = 8080
const portFlagKey = "port"

// Default host address to start listening on.
const hostFlagKey = "host"
const defaultHost = "localhost"

// Default request timeout
const timeoutFlagKey = "timeout"
const defaultTimeout = 10000 // milliseconds

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start a server for querying company information",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString(hostFlagKey)
		port := viper.GetInt(portFlagKey)

		timeoutMillis := viper.GetInt(timeoutFlagKey)
		timeout := time.Duration(timeoutMillis) * time.Millisecond

		return serverAction(host, port, timeout)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Flags
	// Port
	portUsage := fmt.Sprintf("Port server will listen on. (default %d)", defaultPort)
	serverCmd.Flags().Int(portFlagKey, defaultPort, portUsage)
	viper.BindPFlag(portFlagKey, serverCmd.Flags().Lookup(portFlagKey))

	// Host
	hostUsage := fmt.Sprintf("Host to start server on (default %q)", defaultHost)
	serverCmd.Flags().String(hostFlagKey, defaultHost, hostUsage)
	viper.BindPFlag(hostFlagKey, serverCmd.Flags().Lookup(hostFlagKey))

	// Timeout
	timeoutUsage := fmt.Sprintf("Request timeout in milliseconds (default %d)", defaultTimeout)
	serverCmd.Flags().Int(timeoutFlagKey, defaultTimeout, timeoutUsage)
	viper.BindPFlag(timeoutFlagKey, serverCmd.Flags().Lookup(timeoutFlagKey))
}

func serverAction(host string, port int, timeout time.Duration) error {
	// Get ElasticSearch config
	config, err := esConfig()
	if err != nil {
		return err
	}

	// Initialize a new ES client
	client, err := es.NewClient(config)
	if err != nil {
		return err
	}

	// Initialize server
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	server := server.NewServer(addr, timeout, client)

	// Listen for HTTP requests
	fmt.Printf("Listening on %s\n", addr)
	return server.ListenAndServe()
}
