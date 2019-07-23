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
	"log"

	"gitlab.com/fufuhu/ti_rancher_k8s_sampleapp/service"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Todo Serverにログインし、JWTトークンを取得します。",
	Long:  `Todo Serverにログインし、JWTトークンを取得します。`,
	Run:   login,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("login called")
	//},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().String("username", "", "Todoサーバにログインするためのユーザ名")
	loginCmd.Flags().String("password", "", "Todoサーバにログインするためのパスワード")
}

func login(cmd *cobra.Command, args []string) {

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

	username, err := clientSetting.Username()
	if err != nil {
		log.Fatal(err)
	}

	password, err := clientSetting.Password()
	if err != nil {
		log.Fatal(err)
	}

	//loginMessage := service.RequestPing(protocol, host, port)
	loginConfig, err := service.Login(protocol, host, port, username, password)
	if err != nil {
		log.Fatal(err)
	}

	if loginConfig.Filepath, err = rootCmd.PersistentFlags().GetString("config"); err != nil {
		log.Fatal(err)
	}
	config, err := service.CreateConfigFile(loginConfig)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("認証トークンを取得しました。")
	log.Println(config.Token)

}
