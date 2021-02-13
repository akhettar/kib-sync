package es

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	host       string
	user       string
	password   string
}

func NewClient(url, user, password string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Timeout: time.Second * 10, Transport: tr,
	}
	return &Client{c, url, user, password}
}

func (c *Client) Do(path, method string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.host, path)
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	req.SetBasicAuth(c.user, c.password)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	return response, nil
}
