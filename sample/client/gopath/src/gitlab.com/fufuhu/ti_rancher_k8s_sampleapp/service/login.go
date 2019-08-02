package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// JWTAuthMessage は認証の結果として返ってくるJSONメッセージを
// 受け取るための構造体です。
type JWTAuthMessage struct {
	Token string `json:"token"` // 受け取ったトークン本体
}

// LoginConfig はToDoクライアントがToDoサーバに対して認証処理および、
// その結果を受け取る際に必要となる情報をまとめた構造体です
type LoginConfig struct {
	Filepath string // 設定ファイルの保管先パス
	Protocol string // ToDoサーバにアクセスする際のプロトコル
	Host     string // ToDoサーバのFQDN
	Port     int    // ToDoサーバにアクセスする際の宛先TCPポート番号
	Username string // 認証に利用するユーザ名
	Password string // 認証に利用するパスワード
	Token    string // ToDoサーバから取得したトークン
}

// Login はToDoクライアントからToDoサーバへの認証処理を行い、
// 認証トークンを含んだLoginConfigを受け取ります。
func Login(protocol string, host string, port int, username string, password string) (LoginConfig, error) {

	// 前準備
	path := "/api/auth"
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	authInfo := fmt.Sprintf(`{ "username": "%s","password": "%s"}`, username, password)

	req, err := http.NewRequest("POST", url, strings.NewReader(authInfo))
	req.Header.Set("Content-Type", "application/json") //ボディに含まれるコンテンツがJSONであることを明示する

	if err != nil {
		log.Fatal(err)
	}

	// ログインリクエストを発行します。
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// 200 OK以外のステータスが返ってきた場合は異常です。
		log.Fatal(res)
	}

	// レスポンスメッセージのボディの中からJWTを取得します。
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var authMessage JWTAuthMessage
	if err := json.Unmarshal(body, &authMessage); err != nil {
		log.Fatal(err)
	}

	// LoginConfigに詰めなおして返却します。
	loginConfig := LoginConfig{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Token:    authMessage.Token,
	}

	return loginConfig, err
}

// Config はyamlファイルとして保管されているToDoクライアントの設定ファイルを表します。
type Config struct {
	Protocol string `yaml:"protocol"` // ToDoサーバにアクセスする際のプロトコル
	Host     string `yaml:"host"`     // ToDoサーバのFQDN
	Port     int    `yaml:"port"`     // ToDoサーバにアクセスする際の宛先TCPポート番号
	Token    string `yaml:"token"`    // ToDoサーバから取得したトークン
}

// CreateConfigFile はLoginConfigの情報を受け取ってYAML形式でToDoクライアントの設定ファイルを作成します。
func CreateConfigFile(loginConfig LoginConfig) (Config, error) {

	//LoginConfigのままだとYAMLファイルにユーザ名とパスワードがセットで出力されるため、
	//セキュリティ的にあまりよろしくないのでConfigに積み替えて出力します。
	config := Config{
		Protocol: loginConfig.Protocol,
		Host:     loginConfig.Host,
		Port:     loginConfig.Port,
		Token:    loginConfig.Token,
	}

	// 構造体からYAML形式への変換処理を行います。
	out, err := yaml.Marshal(config)

	if err != nil {
		log.Fatal(err)
	}

	// ファイルへの出力
	err = ioutil.WriteFile(loginConfig.Filepath, out, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	return config, err
}
