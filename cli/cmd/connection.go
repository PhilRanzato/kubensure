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

	"github.com/PhilRanzato/kubensure/backend"
	"github.com/spf13/cobra"
)

var podNs string
var svcNs string
var svcPort int

// connectionCmd represents the connection command
var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Check connection from a pod to a service",
	Long: `
Check connection from a pod to a service.

Usage examples:

  # Ensure pod 'example' of namespace 'test' can connect to service 'svc-example' in namespace 'svc-test'

  kubensure connection example -n test svc-example -s svc-test

`,
	// Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cs := backend.GetClientSet()
			pod := backend.GetPodByName(backend.GetPods(cs), args[0], podNs)
			svc := backend.GetServiceByName(backend.GetServices(cs), args[1], svcNs)
			if backend.TestConnectionPodToService(cs, pod, svc, svcPort) {
				fmt.Printf("Pod %s can connect to %s", args[0], args[1])
			} else {
				fmt.Printf("Pod %s cannot connect to %s", args[0], args[1])
			}
		} else {
			fmt.Printf(`'kubensure connection ok' needs at least two arguments: <PodName> and <ServiceName>.
See 'kubensure connection -h' for more information`)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectionCmd)
	connectionCmd.Flags().StringVarP(&podNs, "pod-ns", "n", "default", "Pod namespace")
	connectionCmd.Flags().StringVarP(&svcNs, "svc-ns", "s", "default", "Service namespace")
	connectionCmd.Flags().IntVarP(&svcPort, "svc-port", "p", 0, "Service port")
}
