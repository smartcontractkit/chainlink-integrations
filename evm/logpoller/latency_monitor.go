package logpoller

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
)

const latencyWarningThreshold = 0.7 // Warn if latency exceeds 70% of block production rate

type LatencyMonitorClient interface {
	HeadByNumber(ctx context.Context, n *big.Int) (*types.Head, error)
	HeadByHash(ctx context.Context, n common.Hash) (*types.Head, error)
}

type LatencyMonitorTracker interface {
	LatestAndFinalizedBlock(ctx context.Context) (latest *types.Head, finalized *types.Head, err error)
}

// LatencyMonitor wraps RPC calls with a warning if the latency exceeds the set threshold of block production rate
type LatencyMonitor struct {
	c                   LatencyMonitorClient
	t                   LatencyMonitorTracker
	lggr                logger.Logger
	blockProductionRate time.Duration
}

func NewLatencyMonitor(c LatencyMonitorClient, t LatencyMonitorTracker, lggr logger.Logger, blockProductionRate time.Duration) LatencyMonitor {
	return LatencyMonitor{
		c:                   c,
		t:                   t,
		lggr:                lggr,
		blockProductionRate: blockProductionRate,
	}
}

// latencyMonitoredCall wraps any function and logs a warning if it exceeds the threshold of block production rate
func latencyMonitoredCall[T any](lmc *LatencyMonitor, name string, fn func() (T, error)) (T, error) {
	start := time.Now()
	result, err := fn()
	latency := time.Since(start)

	threshold := time.Duration(float64(lmc.blockProductionRate) * latencyWarningThreshold)
	if latency > threshold {
		lmc.lggr.Warnf(
			"%s latency of %s exceeded threshold of %s (%.0f%% of block production time %s)",
			name, latency, threshold, latencyWarningThreshold*100, lmc.blockProductionRate,
		)
	}

	return result, err
}

func (lmc *LatencyMonitor) HeadByNumber(ctx context.Context, n *big.Int) (*types.Head, error) {
	return latencyMonitoredCall(lmc, "HeadByNumber", func() (*types.Head, error) {
		return lmc.c.HeadByNumber(ctx, n)
	})
}

func (lmc *LatencyMonitor) HeadByHash(ctx context.Context, n common.Hash) (*types.Head, error) {
	return latencyMonitoredCall(lmc, "HeadByHash", func() (*types.Head, error) {
		return lmc.c.HeadByHash(ctx, n)
	})
}

func (lmc *LatencyMonitor) LatestAndFinalizedBlock(ctx context.Context) (latest *types.Head, finalized *types.Head, err error) {
	type HeadPair struct {
		Latest    *types.Head
		Finalized *types.Head
	}
	pair, err := latencyMonitoredCall(lmc, "LatestAndFinalizedBlock", func() (HeadPair, error) {
		latest, finalized, err := lmc.t.LatestAndFinalizedBlock(ctx)
		return HeadPair{latest, finalized}, err
	})
	return pair.Latest, pair.Finalized, err
}
