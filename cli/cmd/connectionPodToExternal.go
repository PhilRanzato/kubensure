/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

var podNsToExternal string
var extPortToExternal int

// connectionPodToExternalCmd represents the connectionPodToExternal command
var connectionPodToExternalCmd = &cobra.Command{
	Use:   "pod-to-ext",
	Short: "Check connection from a pod to an external endpoint.",
	Long: `
Check connection from a pod to an external endpoint.

Usage examples:

  # Ensure pod 'example' of namespace 'test' can connect to https://kubernetes.io

  kubensure connection pod-to-ext example -n test https://kubernetes.io

  # Ensure pod 'example' of namespace 'test' can connect to http://192.168.100.112:90

  kubensure connection pod-to-ext example -n test http://192.168.100.112 --ext-port 90

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cs := backend.GetClientSet()
			pod := backend.GetPodByName(backend.GetPods(cs), args[0], podNsToExternal)
			if backend.ConnectionPodToExternal(cs, pod, args[1], extPortToExternal) {
				fmt.Printf("Pod %s can connect to %s", args[0], args[1])
			} else {
				fmt.Printf("Pod %s cannot connect to %s", args[0], args[1])
			}
		} else {
			fmt.Printf(`'kubensure connection pod-to-pod' needs at least two arguments: <PodName> and <ServiceName>.
See 'kubensure connection pod-to-pod -h' for more information`)
		}
	},
}

func init() {
	connectionCmd.AddCommand(connectionPodToExternalCmd)
	connectionPodToExternalCmd.Flags().StringVarP(&podNsToExternal, "pod-ns", "n", "default", "Pod namespace")
	connectionPodToExternalCmd.Flags().IntVarP(&extPortToExternal, "ext-port", "p", 443, "External endpoint port")
	connectionPodToExternalCmd.SuggestionsMinimumDistance = 2

}
