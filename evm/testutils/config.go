package testutils

import (
	"testing"

	"github.com/smartcontractkit/chainlink-integrations/evm/config"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
)

// Deprecated: use configtest.NewChainScopedConfig
func NewTestChainScopedConfig(t testing.TB, overrideFn func(c *toml.EVMConfig)) config.ChainScopedConfig {
	var chainID = (*big.Big)(FixtureChainID)
	evmCfg := &toml.EVMConfig{
		ChainID: chainID,
		Chain:   toml.Defaults(chainID),
	}

	if overrideFn != nil {
		// We need to get the chainID from the override function first to load the correct chain defaults.
		// Then we apply the override values on top
		overrideFn(evmCfg)
		evmCfg.Chain = toml.Defaults(evmCfg.ChainID)
		overrideFn(evmCfg)
	}

	return config.NewTOMLChainScopedConfig(evmCfg)
}
