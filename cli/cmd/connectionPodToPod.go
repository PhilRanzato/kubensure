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

var podNsToPod string
var targetNsToPod string
var targetPortToPod int

// connectionPodToPodCmd represents the connectionPodToPod command
var connectionPodToPodCmd = &cobra.Command{
	Use:   "pod-to-pod",
	Short: "Check connection from a pod to another pod.",
	Long: `
Check connection from a pod to another pod.

Usage examples:

  # Ensure pod 'example' of namespace 'test' can connect to pod 'target' in namespace 'pod-test'

  kubensure connection pod-to-pod example -n test target -t pod-test

  # Ensure pod 'example' of namespace 'test' can connect to pod 'target' in namespace 'pod-test' on port 8000

  kubensure connection pod-to-pod example -n test target -t pod-test -p 8000

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cs := backend.GetClientSet()
			pod := backend.GetPodByName(backend.GetPods(cs), args[0], podNsToPod)
			trgt := backend.GetPodByName(backend.GetPods(cs), args[1], targetNsToPod)
			if backend.ConnectionPodToPod(cs, pod, trgt, targetPortToPod) {
				fmt.Printf("Pod %s can connect to %s", args[0], args[1])
			} else {
				fmt.Printf("Pod %s cannot connect to %s", args[0], args[1])
			}
		} else {
			fmt.Printf(`'kubensure connection pod-to-pod' needs at least two arguments: <PodName> and <PodNamespace>.
		See 'kubensure connection pod-to-pod -h' for more information`)
		}
	},
}

func init() {
	connectionCmd.AddCommand(connectionPodToPodCmd)

	connectionPodToPodCmd.Flags().StringVarP(&podNsToPod, "pod-ns", "n", "default", "Pod namespace")
	connectionPodToPodCmd.Flags().StringVarP(&targetNsToPod, "target-ns", "t", "default", "Target Pod namespace")
	connectionPodToPodCmd.Flags().IntVarP(&targetPortToPod, "target-port", "p", 0, "Target Pod port")
	connectionPodToPodCmd.SuggestionsMinimumDistance = 2

}
