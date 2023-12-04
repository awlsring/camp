package agent

import (
	"context"
	"sync"
	"time"

	campagent "github.com/awlsring/camp/packages/camp_agent"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Client interface {
	PowerOffMachine(ctx context.Context, id, endpoint, token string) error
	RebootMachine(ctx context.Context, id, endpoint, token string) error
	CheckMachineConnectivity(ctx context.Context, id, endpoint, token string) (bool, error)
}

type AgentTargetConfig struct {
	Identifier string
	Endpoint   string
	ApiKey     string
}

type CachedAgentClientCache struct {
	cache *expirable.LRU[string, *campagent.Client]
	mu    sync.Mutex
}

func NewAgentClientCache() (Client, error) {
	cache := expirable.NewLRU[string, *campagent.Client](20, nil, time.Second*60)
	return &CachedAgentClientCache{
		cache: cache,
	}, nil
}

func (c *CachedAgentClientCache) createClientInCache(ctx context.Context, id, endpoint, token string) (*campagent.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, err := c.makeClient(ctx, endpoint, token)
	if err != nil {
		return nil, err
	}

	c.cache.Add(id, client)
	return client, nil
}

func (c *CachedAgentClientCache) makeClient(ctx context.Context, endpoint, key string) (*campagent.Client, error) {
	client, err := campagent.NewClient(endpoint, NewStaticAuthKeyProvider(key))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *CachedAgentClientCache) GetClient(ctx context.Context, id, endpoint, token string) (*campagent.Client, error) {
	client, ok := c.cache.Get(id)
	if ok {
		return client, nil
	}

	client, err := c.createClientInCache(ctx, id, endpoint, token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *CachedAgentClientCache) PowerOffMachine(ctx context.Context, id, endpoint, token string) error {
	client, err := c.GetClient(ctx, id, endpoint, token)
	if err != nil {
		return err
	}

	_, err = client.PowerOff(ctx, campagent.OptPowerOffRequestContent{})
	if err != nil {
		return err
	}

	return nil
}

func (c *CachedAgentClientCache) RebootMachine(ctx context.Context, id, endpoint, token string) error {
	client, err := c.GetClient(ctx, id, endpoint, token)
	if err != nil {
		return err
	}

	_, err = client.Reboot(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CachedAgentClientCache) CheckMachineConnectivity(ctx context.Context, id, endpoint, token string) (bool, error) {
	client, err := c.GetClient(ctx, id, endpoint, token)
	if err != nil {
		return false, err
	}

	_, err = client.Health(ctx)
	if err != nil {
		return false, nil
	}

	return true, nil
}
