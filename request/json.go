package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type HeaderKV struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func POST(url string, data interface{}, headers ...HeaderKV) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 5}
	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	for _, h := range headers {
		req.Header.Set(h.Key, h.Val)
	}
	return client.Do(req)
}

func GET(url string) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []byte(resp.Status), errors.New(resp.Status)
	}
	return io.ReadAll(resp.Body)
}
