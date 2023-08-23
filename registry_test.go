package gogrpcgeneric

import (
	"testing"

	_ "embed"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func TestNamingClient(t *testing.T) {

	c, _ := NewReadOnlyNamingClient(NacosRegistryConfig{
		IpAddr:      "$nacos",
		Port:        6801,
		NamespaceId: "local",
	})
	xx, _ := c.selectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "nacos.rpc",
		GroupName:   "DEFAULT_GROUP",
	})
	if xx == nil {
		t.Error("fail")
	} else {
		t.Log("ok")
	}
}
