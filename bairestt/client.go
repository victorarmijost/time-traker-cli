package bairestt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type client struct {
	sync.RWMutex
	http.Client
	token  string
	domain string
}

//Creates a new client
func newClient(domain string) *client {
	return &client{
		Client: http.Client{},
		domain: domain,
	}
}

//Adds authorization token
func (c *client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.GetToken()))

	return c.Client.Do(req)
}

//Sets a new auth token
func (c *client) SetToken(token string) {
	c.Lock()
	defer c.Unlock()

	c.token = token
}

func (c *client) GetToken() string {
	c.RLock()
	defer c.RUnlock()

	return c.token
}

func (c *client) IsToken() bool {
	c.RLock()
	defer c.RUnlock()

	return c.token != ""
}

//Generic request
func (c *client) request(verb, path string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(verb, fmt.Sprintf("%s%s", c.domain, path), bytes.NewReader(bodyBytes))

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Post(path string, body interface{}) (*http.Response, error) {
	return c.request("POST", path, body)
}

func (c *client) Put(path string, body interface{}) (*http.Response, error) {
	return c.request("PUT", path, body)
}
