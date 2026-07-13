package xnacos

import "errors"

type Config struct {
	Address   string
	Namespace string
	Group     string
}

type Instance struct {
	ServiceName string
	Host        string
	Port        uint64
}

type Client struct {
	config Config
}

func NewClient(config Config) (*Client, error) {
	if config.Address == "" {
		return nil, errors.New("nacos address is required")
	}
	if config.Group == "" {
		config.Group = "DEFAULT_GROUP"
	}
	return &Client{config: config}, nil
}

func (c *Client) Register(instance Instance) error {
	if instance.ServiceName == "" {
		return errors.New("service name is required")
	}
	if instance.Host == "" || instance.Port == 0 {
		return errors.New("service endpoint is required")
	}
	return nil
}

func (c *Client) Config() Config {
	return c.config
}
