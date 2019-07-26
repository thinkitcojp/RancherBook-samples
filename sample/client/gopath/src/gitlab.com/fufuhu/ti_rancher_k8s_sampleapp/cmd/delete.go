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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "ユーザに紐づくTODOタスクを削除します。",
	Long: `ユーザに紐づくTODOタスクを削除します。
		--id指定なしの場合は、エラーを返します。`,
	Run: delete,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("delete called")
	// },
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().Int("id", 0, "削除したいタスクのID")
}

func delete(cmd *cobra.Command, args []string) {
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
	if err != nil {
		log.Println("JWTトークンの設定が異常です。loginサブコマンドで再取得してください。")
		log.Fatal(err)
	}

	id, _ := taskRequestSetting.ID()

	if id != 0 {
		task, err := service.DeleteTask(protocol, host, port, token, id)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Task(ID=%d) is deleted.\n", task.ID)
		fmt.Printf("ID\tTitle\tDescription\n")
		fmt.Printf("%d\t%s\t%s\n", task.ID, task.Title, task.Description)
	} else {
		log.Fatal("削除対象のタスクのIDが正しく指定されていません(--id)")
	}
}
