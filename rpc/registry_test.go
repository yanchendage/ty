package rpc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestNewRegistryServer(t *testing.T) {

	NewLocalRegistryServer("127.0.0.1:8888")

}

func TestRegistryServerPushApi(t *testing.T) {

	resp, err := http.PostForm("http://127.0.0.1:8888/push",
		url.Values{ "addr": {"127.0.0.1:7729"}})

	if err != nil {
		// handle error
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

	NewLocalRegistryServer("127.0.0.1:8888")

}
