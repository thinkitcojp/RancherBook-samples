package service

import (
	"fmt"
	"os"
	"testing"
)

func TestRequestPing(t *testing.T) {

	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	fmt.Println(testTarget)

	pongMessage := RequestPing("http", testTarget, 8000)

	if pongMessage.Message != "pong" {
		t.Fail()
	}

}
