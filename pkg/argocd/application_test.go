/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 06:52:38
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 06:52:38
 */

package argocd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"testing"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func TestClient_CreateApplication(t *testing.T) {
	type fields struct {
		url    string
		token  string
		client *http.Client
	}
	type args struct {
		app *appv1.Application
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test create application",
			fields: fields{
				url:    os.Getenv("ARGOCD_URL"),
				token:  os.Getenv("ARGOCD_TOKEN"),
				client: http.DefaultClient,
			},
			args: args{
				app: &appv1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mss-boot-admin-beta",
						Namespace: "argocd",
						Labels: map[string]string{
							"Project": "devops",
						},
					},
					Spec: appv1.ApplicationSpec{
						Source: &appv1.ApplicationSource{
							RepoURL:        "https://github.com/mss-boot-io/mss-boot-gitops",
							Path:           "beta/mss-boot/admin",
							TargetRevision: "main",
						},
						Destination: appv1.ApplicationDestination{
							Name:      os.Getenv("ARGOCD_TEST_CLUSTER"),
							Namespace: os.Getenv("ARGOCD_TEST_NAMESPACE"),
						},
						Project: "default",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				url:    tt.fields.url,
				token:  tt.fields.token,
				client: tt.fields.client,
			}
			if err := c.CreateApplication(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("CreateApplication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
