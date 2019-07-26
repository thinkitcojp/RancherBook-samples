// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"log"

	"gitlab.com/fufuhu/ti_rancher_k8s_sampleapp/service"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping-poing APIを使ってToDoサーバの起動を確認します",
	Long: `ToDoサーバに対してping-pong APIを使ったアクセス確認を行うことで
ToDoサーバへの通信の疎通を確認します。`,
	Run: ping,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("ping called")
	// },
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ping(cmd *cobra.Command, args []string) {

	protocol, err := clientSetting.Protocol()
	if err != nil {
		log.Fatal(err)
	}
	host, err := clientSetting.Host()
	if err != nil {
		log.Fatal(err)
	}
	port, err := clientSetting.Port()
	if err != nil {
		log.Fatal(err)
	}

	pongMessage := service.RequestPing(protocol, host, port)
	fmt.Println(pongMessage.Message)
}
