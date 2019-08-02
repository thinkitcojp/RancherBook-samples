package service

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestCreateConfigFile(t *testing.T) {
	loginConfig := LoginConfig{
		Filepath: "test.yaml",
		Protocol: "http",
		Host:     "localhost",
		Port:     80,
		Username: "username",
		Password: "password",
		Token:    "token",
	}

	config, err := CreateConfigFile(loginConfig)

	if err != nil {
		log.Fatal(err)
	}

	// LoginConfigからConfigへの積み替えが正しくできていることを確認する。
	if config.Protocol != loginConfig.Protocol {
		t.Fail()
	}

	if config.Host != loginConfig.Host {
		t.Fail()
	}

	if config.Port != loginConfig.Port {
		t.Fail()
	}

	if config.Token != loginConfig.Token {
		t.Fail()
	}

	//Configの内容がファイルに反映されていることを確認する。
	var fileConfig Config
	fileBuffer, err := ioutil.ReadFile(loginConfig.Filepath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(fileBuffer, &fileConfig)
	if err != nil {
		log.Fatal(err)
	}

	if config.Protocol != fileConfig.Protocol {
		t.Fail()
	}

	if config.Host != fileConfig.Host {
		t.Fail()
	}

	if config.Port != fileConfig.Port {
		t.Fail()
	}

	if config.Token != fileConfig.Token {
		t.Fail()
	}

	// (後片付け)作成されたファイルの削除
	if _, err := os.Stat(loginConfig.Filepath); err != nil {
		// ファイルの存在確認
		log.Fatal(err)
	}
	// ファイルの削除
	err = os.Remove(loginConfig.Filepath)
	if err != nil {
		log.Fatal(err)
	}
}

// TestLogin は正常系のテストです。
// テスト用ユーザ(test_user)を使って
// JWTトークンを取得できるかを確認します。
func TestLogin(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	expect := LoginConfig{
		Protocol: "http",
		Host:     testTarget,
		Port:     8000,
		Username: "test_user",
		Password: "test_password",
	}
	loginConfig, err := Login(expect.Protocol,
		expect.Host,
		expect.Port,
		expect.Username,
		expect.Password)

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if loginConfig.Token == "" {
		t.Fail()
	}
	log.Println(loginConfig.Token)

}
