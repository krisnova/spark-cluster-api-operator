// Copyright © 2019 Kris Nova <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/kris-nova/logger"
	"github.com/kris-nova/spark-cluster-api-operator/operator"
	"os"

	"github.com/spf13/cobra"
)

var sc = &operator.ServerConfiguration{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spark-cluster-api-operator",
	Short: "Autoscale your spark cluster with Kubernetes cluster API",
	Long: `Autoscale your spark cluster with Kubernetes cluster API`,
	Run: func(cmd *cobra.Command, args []string) {

		//
		err := operator.ListenAndWait(sc)
		if err != nil {
			logger.Critical("Error: %s", err)
			os.Exit(1)
		}
		logger.Always("Yay!")
		os.Exit(0)

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().IntVarP(&sc.Port, "port", "p", 80, "The port number for the HTTP REST server to listen on")
	rootCmd.Flags().StringVarP(&sc.BindAddress, "bind", "b", "", "The default bind address to listen on")
}
