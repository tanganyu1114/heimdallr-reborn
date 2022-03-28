package bifrosts

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"testing"
)

func Test_getStreamSrvsProxyBrief(t *testing.T) {
	//conf, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\heimdallr-reborn\\server\\test\\config_test\\nginx.conf")
	conf, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatalf(err.Error())
	}
	proxyB, err := getStreamSrvsProxyBrief(conf)
	if err != nil {
		t.Fatalf(err.Error())
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(jsonData))
}

func Test_getHttpSrvsProxyBrief(t *testing.T) {
	//conf, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\heimdallr-reborn\\server\\test\\config_test\\nginx.conf")
	conf, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatalf(err.Error())
	}
	proxyB, err := getHttpSrvsProxyBrief(conf)
	if err != nil {
		t.Fatalf(err.Error())
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(jsonData))
}
