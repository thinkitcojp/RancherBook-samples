package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

const DefaultTestFilename = "test_config"

// TestProtocolWithDefaultValue は特に何も指定しなかった場合に、
// デフォルトの値(http)がProtocolから取得できることを確認する。
func TestProtocolWithDefaultValue(t *testing.T) {
	protocol, err := clientSetting.Protocol()

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if protocol != "http" {
		t.Fail()
	}
}

// TestHostWithDefaultValue は特に何も指定しなかった場合に、
// デフォルトの値(127.0.0.1)がHostから取得できることを確認する。
func TestHostWithDefaultValue(t *testing.T) {
	host, err := clientSetting.Host()

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if host != "127.0.0.1" {
		t.Fail()
	}
}

// TestPortWithDefaultValue は特に何も指定しなかった場合に、
// デフォルトの値(127.0.0.1)がHostから取得できることを確認する。
func TestPortWithDefaultValue(t *testing.T) {
	port, err := clientSetting.Port()

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if port != 80 {
		t.Fail()
	}
}

// TestProtocolWithOptionOverride は--protocolオプションで
// Protocolの返り値が上書きされることを確認する
func TestProtocolWithOptionOverride(t *testing.T) {
	flags := rootCmd.PersistentFlags()
	err := flags.Set("protocol", "https")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	protocol, err := clientSetting.Protocol()
	if err != nil {
		t.Fail()
	}
	if protocol != "https" {
		t.Fail()
	}
}

// TestHostWithOptionOverride は--hostオプションで
// Hostの返り値が上書きされることを確認する
func TestHostWithOptionOveride(t *testing.T) {
	flags := rootCmd.PersistentFlags()
	err := flags.Set("host", "127.0.0.1")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	host, err := clientSetting.Host()
	if err != nil {
		t.Fail()
	}
	if host != "127.0.0.1" {
		t.Fail()
	}
}

// TestPortWithOptionOverride は--portオプションで
// Portの返り値が上書きされることを確認する
func TestPortWithOptionOverride(t *testing.T) {
	flags := rootCmd.PersistentFlags()
	err := flags.Set("port", "8000")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	port, err := clientSetting.Port()
	if err != nil {
		t.Fail()
	}
	if port != 8000 {
		t.Fail()
	}
}

// Load設定ファイルをオーバーライドするための関数です。
func loadConfigForConfigFileOveride(filename string) {
	currentDirectory, _ := os.Getwd()
	viper.AddConfigPath(currentDirectory)
	// viper.SetConfigName("test_config")
	viper.SetConfigName(filename)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// TestProtocolWithConfigFileOveride は設定ファイルから
// Protocolの返り値を設定できることを確認する
func TestProtocolWithConfigFileOveride(t *testing.T) {
	//protocolを空文字列にしてフラグの当該部分を初期化
	flags := rootCmd.PersistentFlags()
	err := flags.Set("protocol", "")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	//テスト用の設定ファイルの読み込み
	loadConfigForConfigFileOveride(DefaultTestFilename)

	protocol, err := clientSetting.Protocol()
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	if protocol != "https" {
		t.Fail()
	}
}

// TestHostWithConfigFileOveride は設定ファイルから
// Hostの返り値を設定できることを確認する
func TestHostWithConfigFileOverride(t *testing.T) {
	//hostを空文字列にしてフラグの当該部分を初期化
	flags := rootCmd.PersistentFlags()
	err := flags.Set("host", "")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	/*テスト用の設定ファイルの読み込み
	本来は不要だが、将来的にテストコードが肥大化した際に
	テスト用の関数間の依存関係を排除するために必要。
	*/
	loadConfigForConfigFileOveride(DefaultTestFilename)

	host, err := clientSetting.Host()
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	if host != "example.com" {
		log.Fatal(host)
		t.Fail()
	}
}

// TestPortWithConfigFileOverride は設定ファイルから
// Portの返り値を設定できることを確認する
func TestPortWithConfigFileOverride(t *testing.T) {
	//TCPポートに0を設定して初期化。
	flags := rootCmd.PersistentFlags()
	err := flags.Set("port", "0")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	/*テスト用の設定ファイルの読み込み
	本来は不要だが、将来的にテストコードが肥大化した際に
	テスト用の関数間の依存関係を排除するために必要。
	*/
	loadConfigForConfigFileOveride(DefaultTestFilename)

	port, err := clientSetting.Port()
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	if port != 1000 {
		t.Fail()
	}

}

// TestHostWithConfigOverrideWithFlag は設定ファイルの内容よりも
// --hostオプションがHostの返り値として優先されることを確認する
func TestHostWithConfigOverrideWithFlag(t *testing.T) {

	//テスト用設定ファイル読み込み
	loadConfigForConfigFileOveride(DefaultTestFilename)

	flags := rootCmd.PersistentFlags()
	err := flags.Set("host", "example.com")

	host, err := clientSetting.Host()
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if host != "example.com" {
		t.Fail()
	}
}

// TestPortWithConfigOverrideWithFlag は設定ファイルの内容よりも
// --portオプションがPortの返り値として優先されることを確認する
func TestPortWithConfigOverrideWithFlag(t *testing.T) {

	//テスト用設定ファイル読み込み
	loadConfigForConfigFileOveride(DefaultTestFilename)

	flags := rootCmd.PersistentFlags()
	err := flags.Set("port", "30000")

	port, err := clientSetting.Port()
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if port != 30000 {
		t.Fail()
	}
}

// TestProtocolWithConfigOverrideWithFlag は設定ファイルの内容よりも
// --protocolオプションがProtocolの返り値として優先されることを確認する
func TestProtocolWithConfigOverrideWithFlag(t *testing.T) {

	//テスト用設定ファイル読み込み
	loadConfigForConfigFileOveride(DefaultTestFilename)

	flags := rootCmd.PersistentFlags()
	err := flags.Set("protocol", "http")

	protocol, err := clientSetting.Protocol()
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if protocol != "http" {
		t.Fail()
	}
}

// TestUsernameWithCommandLineOpstion は正常系のテストです。
// usernameオプションでユーザ名が指定されている場合に、
// usernameオプションで指定されたユーザ名が
// clientSettingから取得できることを確認します。
func TestUsernameWithCommandLineOption(t *testing.T) {

	expect := "testuser001"

	flags := loginCmd.Flags()
	err := flags.Set("username", expect)

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	var username string
	if username, err = clientSetting.Username(); err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if username != expect {
		t.Fail()
	}
}

// TestUsernameWithoutCommandLineOption は異常系のテストです。
// usernameオプションの指定が存在しない場合に、
// clientSetting.Usernameからエラーが返されることを期待します。
func TestUsernameWithoutCommandLineOption(t *testing.T) {
	flags := loginCmd.Flags()
	err := flags.Set("username", "")
	_, err = clientSetting.Username()

	if err == nil {
		t.Fail()
	}

	log.Println(err)
	if err.Error() != SettingErrorMessageUsernameNotFound {
		t.Fail()
	}
}

// TestPasswordWithCommandLineOption は正常系のテストです。
// --passwordオプションでパスワードが指定されている場合に、
// passwordオプションで指定されたパスワードが
// clientSettingから取得可能かを確認します。
func TestPasswordWithCommandLineOption(t *testing.T) {

	expect := "testpassword001"

	flags := loginCmd.Flags()
	err := flags.Set("password", expect)

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	var password string
	if password, err = clientSetting.Password(); err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if password != expect {
		t.Fail()
	}
}

// TestPasswordWithoutCommandLineOption は異常系のテストです。
// passwordオプションの指定が存在しない場合に、
// clientSetting.Passwordからエラーが返されることを期待します。
func TestPasswordWithoutCommandLineOption(t *testing.T) {
	flags := loginCmd.Flags()
	err := flags.Set("password", "")
	_, err = clientSetting.Password()

	if err == nil {
		t.Fail()
	}

	log.Println(err)
	if err.Error() != SettingErrorMessagePasswordNotFound {
		t.Fail()
	}
}

// TestTokenWithConfigFile は正常系のテストです。
// 設定ファイルが存在する場合に、
// clientSettingから認証トークンを取得できることを確認します。
func TestTokenWithConfigFile(t *testing.T) {
	loadConfigForConfigFileOveride(DefaultTestFilename)

	expect := "test_token"

	token, err := clientSetting.Token()
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	if token != expect {
		t.Fail()
	}
}

// TestTokenWithoutConfigFile は異常系のテストです。
// 設定ファイルにトークン情報が存在しない場合に、
// clientSettingから認証トークンを取得できることを確認します。
func TestTokenWithoutConfigFile(t *testing.T) {
	loadConfigForConfigFileOveride("test_config_without_token")

	token, err := clientSetting.Token()

	if token != "" {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}

	log.Println(err)
	if err.Error() != SettingErrorMessageTokenNotFound {
		t.Fail()
	}

}
