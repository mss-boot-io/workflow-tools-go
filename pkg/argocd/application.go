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
		fmt.Printf("create request failed, err: %s\n", err.Error())
		return false
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	response, err := c.GetClient().Do(request)
	if err != nil {
		fmt.Printf("request failed, err: %s\n", err.Error())
		return false
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusForbidden {
		fmt.Printf("application %s not found\n", name)
		return false
	}
	newApp := &appv1.Application{}
	err = json.NewDecoder(response.Body).Decode(newApp)
	if err != nil {
		fmt.Printf("decode response failed, err: %s\n", err)
		return false
	}
	return name == app.Name
}

func (c *Client) CreateApplication(app *appv1.Application) error {
	fmt.Printf("create application %s start \n", app.Name)
	//if c.existApplication(app) {
	//	fmt.Printf("application %s already exists \n", app.Name)
	//	return nil
	//}
	return c.controlApplication(app, c.existApplication(app))
}

func (c *Client) controlApplication(app *appv1.Application, exist bool) error {
	app.Status = appv1.ApplicationStatus{}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(app)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v1/applications", c.url)
	method := http.MethodPost
	if exist {
		method = http.MethodPut
		url = fmt.Sprintf("%s/api/v1/applications/%s", c.url, app.Name)
	}
	fmt.Println(buf.String())
	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		fmt.Printf("request failed, err: %s\n", err.Error())
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("request failed, err: %s\n", err.Error())
		return err
	}
	defer response.Body.Close()
	fmt.Printf("create application %s response status code %d\n", app.Name, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("create application %s failed, status code: %d\n", app.Name, response.StatusCode)
		rb, _ := io.ReadAll(response.Body)
		log.Println(string(rb))
		return fmt.Errorf("create application %s failed, status code: %d", app.Name, response.StatusCode)
	}
	err = json.NewDecoder(response.Body).Decode(app)
	if err != nil {
		fmt.Printf("decode response failed, err: %s\n", err.Error())
		return err
	}
	return nil
}
