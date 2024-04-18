package bifrosts

import (
	"encoding/json"
	"gin-vue-admin/internal/pkg/bifrosts"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_getStreamSrvsProxyBrief(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	client, err := mockBifrostsManager.GetBifrostClient(webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	data, err := client.WebServerConfig().Get(webSrvOpts.ServerName)
	if err != nil {
		t.Fatal(err)
	}
	conf, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	proxyB, err := getStreamSrvsProxyBrief(conf)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func Test_getHttpSrvsProxyBrief(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	client, err := mockBifrostsManager.GetBifrostClient(webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	data, err := client.WebServerConfig().Get(webSrvOpts.ServerName)
	if err != nil {
		t.Fatal(err)
	}
	conf, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	proxyB, err := getHttpSrvsProxyBrief(conf)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}
