/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 06:31:01
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 06:31:01
 */

package argocd

import "net/http"

type Client struct {
	url    string
	token  string
	client *http.Client
}

func (c *Client) GetClient() *http.Client {
	if c.client == nil {
		return http.DefaultClient
	}
	return &(*c.client)
}

// New create a new argocd client
func New(url, token string, client *http.Client) *Client {
	return &Client{
		url:    url,
		token:  token,
		client: client,
	}
}
