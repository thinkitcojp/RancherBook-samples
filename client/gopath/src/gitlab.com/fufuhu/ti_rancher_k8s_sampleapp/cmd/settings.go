package cmd

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// ClientSetting はToDoクライアントが
// ToDoサーバにアクセスする際の情報を格納します。
type ClientSetting struct {
	// Protocol クライアントがサーバにアクセスする際のプロトコル(http/https)
	Protocol func() (string, error)
	// クライアントがサーバにアクセスする際のホスト名(IPアドレス/ホスト名)
	Host func() (string, error)
	// クライアントがサーバにアクセスする際のポート番号(0より大きい整数の値)
	Port func() (int, error)
	// クライアントがサーバにアクセスする際のユーザ名、loginコマンド時に利用する
	Username func() (string, error)
	// クライアントがサーバにアクセスする際のパスワード、loginコマンド時に利用する
	Password func() (string, error)
	// クライアントがサーバにアクセスする際の認証トークン(JWT)
	Token func() (string, error)
}

// SettingErrorMessageUsernameNotFound はユーザ名がusernameオプションで定義
// されていない場合に発生するエラーに含まれるエラーメッセージです。
const SettingErrorMessageUsernameNotFound = "ユーザ名情報が設定されていません。"

// SettingErrorMessagePasswordNotFound はパスワードがpassswordオプションで定義
// されていない場合に発生するエラーに含まれるエラーメッセージです。
const SettingErrorMessagePasswordNotFound = "パスワード情報が指定されていません。"

// SettingErrorMessageTokenNotFound はJWTによる認証トークンが設定ファイルから
// 取得できない場合に発生するエラーに含まれるエラーメッセージです。
const SettingErrorMessageTokenNotFound = "トークン情報が見つかりません。"

var clientSetting ClientSetting

func init() {

	// clientSetting.Protocol 設定ファイルおよびコマンドラインオプション(--protocol)
	// からToDoサーバにアクセスする際のプロトコル(http or https)でアクセスします。
	// 指定が何も無い場合はhttpでアクセスします。
	clientSetting.Protocol = func() (string, error) {
		var protocol string

		// 設定ファイルからの読み込み
		if protocolFromConfig := viper.GetString("protocol"); protocolFromConfig != "" {
			protocol = protocolFromConfig
		}

		// コマンドオプションからの読み込み
		protocolFromOption, err := rootCmd.PersistentFlags().GetString("protocol")
		if err != nil {
			log.Println(err)
		}

		if protocolFromOption != "" {
			protocol = protocolFromOption
		}

		// いずれの場合も値が得られなければデフォルトの値(http)を設定する。
		if protocol == "" {
			protocol = "http"
		}
		return protocol, err
	}

	// clientSetting.Host 設定ファイルおよびコマンドラインオプション(--host)から
	// ToDoサーバのFQDNを取得します。指定が何も無い場合は"127.0.0.1"を利用します。
	clientSetting.Host = func() (string, error) {

		var host string

		// 設定ファイルからの読み込み
		if hostFromConfig := viper.GetString("host"); hostFromConfig != "" {
			host = hostFromConfig
		}

		// コマンドオプションからの読み込み
		hostFromOption, err := rootCmd.PersistentFlags().GetString("host")
		if hostFromOption != "" {
			host = hostFromOption
		}

		if err != nil {
			log.Println(err)
		}

		// いずれの場合も値が得られなければデフォルトの値(127.0.0.1)を設定する。
		if host == "" {
			host = "127.0.0.1"
		}
		return host, nil
	}

	// clientSetting.Port 設定ファイルおよびコマンドラインオプション(--port)から
	// ToDoサーバにアクセスする際の宛先TCPポート番号を指定します。
	// いずれも値が得られない場合は80番ポートを利用します。
	clientSetting.Port = func() (int, error) {
		port := 0

		// 設定ファイルからの読み込み
		if portFromConfig := viper.GetInt("port"); portFromConfig != 0 {
			port = portFromConfig
		}

		// コマンドオプションからの読み込み
		portFromConfig, err := rootCmd.PersistentFlags().GetInt("port")

		if err != nil {
			log.Println(err)
		}

		if portFromConfig != 0 {
			port = portFromConfig
		}

		// いずれも値が得られなければデフォルトの値(80)を設定する。
		if port == 0 {
			port = 80
		}
		return port, err
	}

	// clientSetting.Username コマンドラインオプション(--username)からユーザ名情報を読み込む
	clientSetting.Username = func() (string, error) {
		var username string
		//コマンドオプションからの読み込み
		usernameFromConfig, err := loginCmd.Flags().GetString("username")
		if err != nil {
			log.Fatal(err)
		}

		if usernameFromConfig != "" {
			username = usernameFromConfig
		} else {
			err = errors.New(SettingErrorMessageUsernameNotFound)
		}

		return username, err
	}

	// clientSetting.Password コマンドラインオプション(--password)からパスワード情報を読み込む
	clientSetting.Password = func() (string, error) {
		var password string

		//コマンドオプションからの読み込み
		passwordFromConfig, err := loginCmd.Flags().GetString("password")
		if err != nil {
			log.Fatal(err)
		}

		if passwordFromConfig != "" {
			password = passwordFromConfig
		} else {
			err = errors.New(SettingErrorMessagePasswordNotFound)
		}

		return password, err
	}

	// Token JWTの認証トークン情報を設定ファイルから取得する
	clientSetting.Token = func() (string, error) {

		var token string
		var err error
		if tokenFromConfig := viper.GetString("token"); tokenFromConfig != "" {
			token = tokenFromConfig
		} else {
			err = errors.New(SettingErrorMessageTokenNotFound)
		}
		return token, err
	}
}
