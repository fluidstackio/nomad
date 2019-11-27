package consul

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/nomad/structs"
)

// Implementation of ConsulTokenAPI used to interact with Nomad Server from
// Nomad Client for acquiring Consul Service Identity tokens.
//
// This client is split from the other consul client(s) to avoid a circular
// dependency between themselves and client.Client
type tokenClient struct {
	tokenDeriver TokenDeriverFunc
	logger       hclog.Logger
}

func NewTokensClient(logger hclog.Logger, tokenDeriver TokenDeriverFunc) *tokenClient {
	return &tokenClient{
		tokenDeriver: tokenDeriver,
		logger:       logger,
	}
}

func (c *tokenClient) DeriveSITokens(alloc *structs.Allocation, tasks []string) (map[string]string, error) {
	tokens, err := c.tokenDeriver(alloc, tasks)
	if err != nil {
		c.logger.Error("error deriving SI token", "error", err, "alloc_id", alloc.ID, "task_names", tasks)
		return nil, err
	}
	return tokens, nil
}
