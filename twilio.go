package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	HostName   = "https://api.twilio.com"
	ApiVersion = "/2010-04-01"
)

type Client struct {
	AccountSid, AuthToken string
}

func NewClient(accountSid, authToken string) *Client {
	return &Client{AccountSid: accountSid, AuthToken: authToken}
}

func (c *Client) CreateCall(attrs url.Values) ([]byte, error) {
	return c.Create("Calls", attrs)
}

func (c *Client) Create(resourceName string, attrs url.Values) ([]byte, error) {
	rawurl := c.resourcePath(resourceName) + ".json"
	println("Raw url", rawurl)
	u, err := url.Parse(rawurl)
	if err != nil {
		return []byte{}, err
	}
	u.User = url.UserPassword(c.AccountSid, c.AuthToken)
	resp, err := http.PostForm(u.String(), attrs)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) resourcePath(resourceName string) string {
	return HostName + ApiVersion + "/Accounts/" + c.AccountSid + "/" + resourceName
}
