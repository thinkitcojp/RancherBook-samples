package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type PongMessage struct {
	Message string `json:message`
}

/*
	RequestPingはToDoサーバに対してのPingRequestを行い、PongMessageを返します。
	protocol: プロトコル(http/https)を指定
	host:
*/
func RequestPing(protocol string, host string, port int) PongMessage {

	var pong PongMessage

	path := "/api/ping"
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(res)
	}

	if err := json.Unmarshal(body, &pong); err != nil {
		log.Fatal(err)
	}

	return pong
}
