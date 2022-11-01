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
	"log"
	"net/http"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// existApplication check if the application exists
func (c *Client) existApplication(app *appv1.Application) bool {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/applications/%s", c.url, app.Name),
		nil)
	if err != nil {
		log.Printf("create request failed, err: %v", err)
		return false
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	response, err := c.client.Do(request)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return false
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(app)
	if err != nil {
		log.Printf("decode response failed, err: %v", err)
		return false
	}
	return true
}

func (c *Client) CreateApplication(app *appv1.Application) error {
	log.Printf("create application %s start \n", app.Name)
	if c.existApplication(app) {
		return nil
	}
	app.Status = appv1.ApplicationStatus{}
	app.ObjectMeta = metav1.ObjectMeta{
		Name:      app.Name,
		Namespace: app.Namespace,
	}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(app)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/applications", c.url),
		buf)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	response, err := c.client.Do(request)
	if err != nil {
		log.Printf("request failed, err: %v", err)
		return err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(app)
	if err != nil {
		log.Printf("decode response failed, err: %v", err)
		return err
	}
	return nil
}
