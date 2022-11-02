// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: LicenseRef-Intel

package main

import (
    "github.com/onosproject/onos-lib-go/pkg/logging"
	manager "github.com/tmp/out/roc-api-server-sca-0.1.x/pkg/manager"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var log = logging.GetLogger()

func main() {
	if err := getRootCommand().Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

func getRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run:   runRootCommand,
	}

	cmd.Flags().String("logLevel", "INFO", "Set the log level (DEBUG, INFO)")
	cmd.Flags().String("caPath", "", "path to CA certificate")
	cmd.Flags().String("keyPath", "", "path to client private key")
	cmd.Flags().String("certPath", "", "path to client certificate")
	cmd.Flags().String("gnmiEndpoint", "onos-config:5150", "address of onos-config")
	cmd.Flags().String("topoEndpoint", "onos-topo:5150", "address of onos-topo")
	cmd.Flags().Duration("gnmiTimeout", 10*time.Second, "timeout for the gnmi requests")
	cmd.Flags().Uint("port", 8181, "http port")
	cmd.Flags().Bool("validateResp", true, "Validate response are compliant with OpenAPI3 schema")
	cmd.Flags().StringSlice("allowCorsOrigin", []string{}, "URLs of CORS origins")
	return cmd
}

func runRootCommand(cmd *cobra.Command, args []string) {
	caPath, _ := cmd.Flags().GetString("caPath")
	keyPath, _ := cmd.Flags().GetString("keyPath")
	certPath, _ := cmd.Flags().GetString("certPath")
	gnmiEndpoint, _ := cmd.Flags().GetString("gnmiEndpoint")
	topoEndpoint, _ := cmd.Flags().GetString("topoEndpoint")
	gnmiTimeout, _ := cmd.Flags().GetDuration("gnmiTimeout")
	allowCorsOrigin, _ := cmd.Flags().GetStringSlice("allowCorsOrigin")
	port, _ := cmd.Flags().GetUint("port")
	validateResp, _ := cmd.Flags().GetBool("validateResp")
	logLevel, _ := cmd.Flags().GetString("logLevel")

	cfg := manager.Config{
		CaPath:            caPath,
		KeyPath:           keyPath,
		CertPath:          certPath,
		Port:              port,
		GnmiEndpoint:      gnmiEndpoint,
		GnmiTimeout:       gnmiTimeout,
		ValidateResponses: validateResp,
		AllowCorsOrigins:  allowCorsOrigin,
		TopoEndpoint:      topoEndpoint,
		LogLevel:          logLevel,
	}

	mgr := manager.NewManager(cfg)
	mgr.Run()
}