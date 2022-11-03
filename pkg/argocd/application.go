/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 06:27:34
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 06:27:34
 */

package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// existApplication check if the application exists
func (c *Client) existApplication(app *appv1.Application) bool {
	name := app.Name
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/applications/%s", c.url, name),
		nil)
	if err != nil {
		log.Printf("create request failed, err: %v", err)
		return false
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	response, err := c.GetClient().Do(request)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		log.Printf("application %s not found", name)
		return false
	}
	newApp := &appv1.Application{}
	err = json.NewDecoder(response.Body).Decode(newApp)
	if err != nil {
		log.Printf("decode response failed, err: %v", err)
		return false
	}
	return name == app.Name
}

func (c *Client) CreateApplication(app *appv1.Application) error {
	log.Printf("create application %s start \n", app.Name)
	if c.existApplication(app) {
		log.Printf("application %s already exists \n", app.Name)
		return nil
	}
	return c.createApplication(app)
}

func (c *Client) createApplication(app *appv1.Application) error {
	app.Status = appv1.ApplicationStatus{}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(app)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v1/applications", c.url)
	method := http.MethodPost
	//log.Printf("method %s, url %s", method, url)
	fmt.Println(buf.String())
	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return err
	}
	defer response.Body.Close()
	log.Printf("create application %s response status code %d", app.Name, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		log.Printf("create application %s failed, status code: %d", app.Name, response.StatusCode)
		rb, _ := io.ReadAll(response.Body)
		log.Println(string(rb))
		return fmt.Errorf("create application %s failed, status code: %d", app.Name, response.StatusCode)
	}
	err = json.NewDecoder(response.Body).Decode(app)
	if err != nil {
		log.Printf("decode response failed, err: %v", err)
		return err
	}
	return nil
}
