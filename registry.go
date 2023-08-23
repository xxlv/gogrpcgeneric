package gogrpcgeneric

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosNameClientReadOnly struct {
	client naming_client.INamingClient
}

// NacosRegistryConfig
type NacosRegistryConfig struct {
	IpAddr      string
	Port        uint64
	NamespaceId string
}

func NewReadOnlyNamingClient(config NacosRegistryConfig) (*NacosNameClientReadOnly, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(config.IpAddr, config.Port),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(config.NamespaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
	)
	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	return &NacosNameClientReadOnly{client: client}, err

}

// getService.
func (c *NacosNameClientReadOnly) getService(param vo.GetServiceParam) (*model.Service, error) {
	service, err := c.client.GetService(param)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// selectAllInstances.
func (c *NacosNameClientReadOnly) selectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	instances, err := c.client.SelectAllInstances(param)
	if err != nil {
		return nil, err
	}
	return instances, err
}

// selectInstances.
func (c *NacosNameClientReadOnly) selectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	instances, err := c.client.SelectInstances(param)
	if err != nil {
		return nil, err
	}

	return instances, err
}

// selectOneHealthyInstance.
func (c *NacosNameClientReadOnly) selectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instances *model.Instance, err error) {
	instances, err = c.client.SelectOneHealthyInstance(param)
	return
}

// subscribe.
func (c *NacosNameClientReadOnly) subscribe(param *vo.SubscribeParam) {
	c.client.Subscribe(param)
}

// unSubscribe.
func (c *NacosNameClientReadOnly) unSubscribe(param *vo.SubscribeParam) {
	c.client.Unsubscribe(param)
}

// getAllService.
func (c *NacosNameClientReadOnly) getAllService(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	service, err := c.client.GetAllServicesInfo(param)
	return service, err
}
