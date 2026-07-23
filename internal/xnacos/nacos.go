package xnacos

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
	"github.com/kitex-contrib/registry-nacos/nacos"
	nacosregistry "github.com/kitex-contrib/registry-nacos/registry"
	nacosresolver "github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewRegistry(addr string) (kitexregistry.Registry, error) {
	host, port, ok, err := parseAddr(addr)
	if err != nil || !ok {
		return nil, err
	}

	cli, err := nacos.NewDefaultNacosClient(nacos.Option{F: func(param *vo.NacosClientParam) {
		param.ServerConfigs = []constant.ServerConfig{
			*constant.NewServerConfig(host, uint64(port)),
		}
	}})
	if err != nil {
		return nil, err
	}
	return nacosregistry.NewNacosRegistry(cli), nil
}

func NewResolver(addr string) (discovery.Resolver, error) {
	host, port, ok, err := parseAddr(addr)
	if err != nil || !ok {
		return nil, err
	}

	cli, err := nacos.NewDefaultNacosClient(nacos.Option{F: func(param *vo.NacosClientParam) {
		param.ServerConfigs = []constant.ServerConfig{
			*constant.NewServerConfig(host, uint64(port)),
		}
	}})
	if err != nil {
		return nil, err
	}
	return nacosresolver.NewNacosResolver(cli), nil
}

func parseAddr(addr string) (string, int, bool, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0, false, nil
	}
	host, portText, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, false, fmt.Errorf("parse nacos addr %q: %w", addr, err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port <= 0 {
		return "", 0, false, fmt.Errorf("parse nacos port %q: %w", portText, err)
	}
	return host, port, true, nil
}
