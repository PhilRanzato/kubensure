package cmd

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

import (
	"github.com/spf13/cobra"
)

// connectionCmd represents the connection command
var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Check connection from a pod to a service or to another pod or to an external endpoint.",
	Long: `
Check connection from a pod to a service or to another pod or to an external endpoint.

`,
	// Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(connectionCmd)
	rootCmd.SuggestionsMinimumDistance = 2
}
