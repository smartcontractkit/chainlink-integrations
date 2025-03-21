package config

import (
	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
)

type workflowConfig struct {
	c toml.Workflow
}

func (b *workflowConfig) FromAddress() *types.EIP55Address {
	return b.c.FromAddress
}

func (b *workflowConfig) ForwarderAddress() *types.EIP55Address {
	return b.c.ForwarderAddress
}

func (b *workflowConfig) GasLimitDefault() *uint64 {
	return b.c.GasLimitDefault
}
