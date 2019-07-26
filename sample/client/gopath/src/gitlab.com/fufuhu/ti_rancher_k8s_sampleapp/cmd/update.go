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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "ユーザに紐づくTODOタスクを更新します。",
	Long: `ユーザに紐づくTODOタスクを更新します。
		--id指定なしの場合は、エラーを返します。`,
	Run: update,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().Int("id", 0, "更新したいタスクのID")
	updateCmd.Flags().String("title", "", "更新後のタスクの名前")
	updateCmd.Flags().String("description", "", "更新後のタスクの説明")
	updateCmd.Flags().String("status", "", "更新後のタスクのステータス")
}

func update(cmd *cobra.Command, args []string) {
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
	title, _ := taskRequestSetting.Title()
	description, _ := taskRequestSetting.Description()
	status, _ := taskRequestSetting.Status()

	if id != 0 {
		task, err := service.UpdateTask(protocol, host, port, token, id, title, description, status)

		if err != nil {
			log.Println(err)
		}

		log.Printf("Task(ID=%d) is updated.\n", task.ID)
		fmt.Printf("ID\tTitle\tStatus\tDescription\n")
		fmt.Printf("%d\t%s\t%s\t%s\n", task.ID, task.Title, task.Status, task.Description)
	} else {
		log.Fatal("更新対象のタスクのIDが正しく指定されていません(--id)")
	}
}
