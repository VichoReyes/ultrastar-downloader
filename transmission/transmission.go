package transmission

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const rpc = "http://localhost:9091/transmission/rpc"
const contentType = "application/x-www-form-urlencoded"
const sessionIDHeader = "X-Transmission-Session-Id"

var client = http.Client{Timeout: 10 * time.Second}
var sessionID string

type message struct {
	Method    string `json:"method"`
	Tag       int    `json:"tag"`
	arguments `json:"arguments"`
	Result    string `json:"result,omitempty"`
}

type arguments struct {
	Filename    string    `json:"filename,omitempty"`
	DownloadDir string    `json:"download-dir,omitempty"`
	Fields      []string  `json:"fields,omitempty"`
	Torrents    []torrent `json:"torrents,omitempty"`
}

type torrent struct {
	Name string `json:"name"`
}

// Add adds a torrent to transmission.
// If target is non-empty, the download will be stored there
// instead of the default download location.
func Add(uri, target string) error {
	mes := message{
		Method: "torrent-add",
		Tag:    8,
	}
	mes.Filename = uri
	mes.DownloadDir = target

	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("adding %v: %v", uri, err)
	}

	response, err := submit(body)
	if err != nil {
		return err
	} else if response.Result != "success" {
		resp, _ := json.MarshalIndent(response, "", "\t")
		return errors.New(string(resp))
	}
	return nil
}

// List returns a set of all torrent names
func List() (map[string]bool, error) {
	req := message{
		Method: "torrent-get",
		Tag:    4,
	}
	req.Fields = []string{"name"}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := submit(body)
	if err != nil {
		return nil, err
	} else if response.Result != "success" {
		resp, _ := json.MarshalIndent(response, "", "\t")
		return nil, errors.New(string(resp))
	}

	ret := make(map[string]bool)
	for _, t := range response.Torrents {
		ret[t.Name] = true
	}
	return ret, nil
}

func submit(body []byte) (*message, error) {
	req, err := http.NewRequest("POST", rpc, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if sessionID == "" {
		err := putSessionID(req)
		if err != nil {
			return nil, err
		}
		return submit(body)
	}
	req.Header.Set(sessionIDHeader, sessionID)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := new(message)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
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
