/*
Copyright © 2021 Phil Ranzato philranzato@gmail.com

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

	backend "github.com/PhilRanzato/kubensure/backend"
	"github.com/spf13/cobra"
)

var podNs string
var svcNs string
var svcPort int

// connCheckCmd represents the connCheck command
var connCheckCmd = &cobra.Command{
	Use:   "conn-test",
	Short: "Check connection from a pod to a service",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cs := backend.GetClientSet()
			pod := backend.GetPodByName(backend.GetPods(cs), args[0], podNs)
			svc := backend.GetServiceByName(backend.GetServices(cs), args[1], svcNs)
			if backend.TestConnectionPodToService(cs, pod, svc, svcPort) {
				fmt.Printf("Connection succeded")
			} else {
				fmt.Printf("Cannot Connect to %s", args[1])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(connCheckCmd)

	connCheckCmd.Flags().StringVarP(&podNs, "pod-ns", "n", "default", "Pod namespace")
	connCheckCmd.Flags().StringVarP(&svcNs, "svc-ns", "s", "default", "Service namespace")
	connCheckCmd.Flags().IntVarP(&svcPort, "svc-port", "p", 0, "Service port")
}
