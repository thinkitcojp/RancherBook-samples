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

	"github.com/spf13/cobra"
	"gitlab.com/fufuhu/ti_rancher_k8s_sampleapp/service"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "TODOタスクを作成します。",
	Long:  `TODOタスクを作成します。`,
	Run:   create,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("create called")
	// },
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().String("title", "", "タスクのタイトル")
	createCmd.Flags().String("description", "", "タスクの概要")
}

func create(cmd *cobra.Command, args []string) {

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

	token, err := clientSetting.Token()

	title, err := taskRequestSetting.Title()
	if err != nil {
		log.Println("タスクの名前指定(--title)が不正です。")
		log.Fatal(err)
	}

	description, err := taskRequestSetting.Description()
	if err != nil {
		log.Println("タスクの概要指定(--description)が不正です。")
		log.Fatal(err)
	}

	task, err := service.CreateTask(protocol, host, port, token, title, description)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %d\n", task.ID)
	fmt.Println("TITLE: " + task.Title)
	fmt.Println("DESCRIPTION: ")
	fmt.Println(task.Description)

}
