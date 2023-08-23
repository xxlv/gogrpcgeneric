package gogrpcgeneric

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// All service
type serviceMap map[string][]string

var Defaultconfig *NacosRegistryConfig

// TODO this client need closed when system shutdown
var _nacosNameClientMap map[string]*NacosNameClientReadOnly

func GetClientGlobal(config NacosRegistryConfig) (*NacosNameClientReadOnly, error) {
	key := asKey(&config)
	var err error
	var c *NacosNameClientReadOnly
	if client, ok := _nacosNameClientMap[key]; ok {
		c = client
	} else {
		c, err = NewReadOnlyNamingClient(config)
		if err != nil {
			//  TODO  handle error
		}
	}
	return c, err
}

func LoadServiceDefault(service, group string) (*model.Instance, error) {
	if Defaultconfig == nil {
		return nil, fmt.Errorf("can not find default config of nacos,extra info is: [%s/%s]", service, group)
	}
	return LoadServiceLocation(*Defaultconfig, service, group)
}

// LoadServiceLocation query service location from nacos registry
func LoadServiceLocation(config NacosRegistryConfig, service, group string) (*model.Instance, error) {
	if service == "" {
		return nil, fmt.Errorf("service is empty")
	}
	if group == "" {
		group = "DEFAULT_GROUP"
	}
	c, _ := GetClientGlobal(config)
	// look for healthy instance
	inf, err := c.selectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: service,
		GroupName:   group,
	})
	if inf == nil {
		return nil, err
	}
	return inf, nil
}

func asKey(config *NacosRegistryConfig) string {
	if config == nil {
		return ""
	}
	return fmt.Sprintf("%s%s%s", config.IpAddr, fmt.Sprint(config.Port), config.NamespaceId)
}
