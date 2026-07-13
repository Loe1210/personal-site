package xnacos

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

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

type Discovery interface {
	RegisterService(ctx context.Context, serviceName string, host string, port uint64) error
	ResolveService(ctx context.Context, serviceName string) (string, error)
}

type Client struct {
	config Config
	memory *MemoryClient
}

func NewClient(config Config) (*Client, error) {
	if config.Address == "" {
		return nil, errors.New("nacos address is required")
	}
	if config.Group == "" {
		config.Group = "DEFAULT_GROUP"
	}
	return &Client{config: config, memory: NewMemoryClient()}, nil
}

func (c *Client) Register(instance Instance) error {
	return c.RegisterService(context.Background(), instance.ServiceName, instance.Host, instance.Port)
}

func (c *Client) RegisterService(ctx context.Context, serviceName string, host string, port uint64) error {
	return c.memory.RegisterService(ctx, serviceName, host, port)
}

func (c *Client) ResolveService(ctx context.Context, serviceName string) (string, error) {
	return c.memory.ResolveService(ctx, serviceName)
}

func (c *Client) Config() Config {
	return c.config
}

type MemoryClient struct {
	mu        sync.RWMutex
	instances map[string]string
}

func NewMemoryClient() *MemoryClient {
	return &MemoryClient{instances: make(map[string]string)}
}

func (c *MemoryClient) RegisterService(ctx context.Context, serviceName string, host string, port uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if serviceName == "" {
		return errors.New("service name is required")
	}
	if host == "" || port == 0 {
		return errors.New("service endpoint is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[serviceName] = fmt.Sprintf("http://%s:%d", host, port)
	return nil
}

func (c *MemoryClient) ResolveService(ctx context.Context, serviceName string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if serviceName == "" {
		return "", errors.New("service name is required")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	endpoint, ok := c.instances[serviceName]
	if !ok {
		return "", errors.New("service not found")
	}
	return endpoint, nil
}
