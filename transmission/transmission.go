package transmission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const rpc = "http://localhost:9091/transmission/rpc"
const contentType = "application/x-www-form-urlencoded"
const sessionIDHeader = "X-Transmission-Session-Id"

var client http.Client
var sessionID string

type message struct {
	Method    string `json:"method"`
	Tag       int    `json:"tag"`
	arguments `json:"arguments"`
}

type arguments struct {
	Filename string `json:"filename"`
}

// Add adds a torrent to transmission
func Add(uri string) error {
	mes := message{
		Method: "torrent-add",
		Tag:    8,
	}
	mes.Filename = uri
	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("adding %v: %v", uri, err)
	}

	req, err := http.NewRequest("POST", rpc, bytes.NewReader(body))
	if err != nil {
		return err
	}

	if sessionID == "" {
		putSessionID(req)
	}
	req.Header.Set(sessionIDHeader, sessionID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func putSessionID(req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	sessionID = resp.Header.Get(sessionIDHeader)
	if sessionID == "" {
		return fmt.Errorf("session ID header not returned")
	}
	req.Header.Set(sessionIDHeader, sessionID)
	return nil
}
