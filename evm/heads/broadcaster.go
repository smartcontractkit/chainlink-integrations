package heads

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-framework/chains/heads"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
)

type broadcaster = heads.Broadcaster[*evmtypes.Head, common.Hash]

func NewBroadcaster(
	lggr logger.Logger,
) broadcaster {
	return heads.NewBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
