package keystest

import (
	"context"
	"errors"
	"math/big"
	"math/rand/v2"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-integrations/evm/keys"
	"github.com/smartcontractkit/chainlink-integrations/evm/keys/internal"
)

var _ keys.ChainStore = &FakeChainStore{}

// FakeChainStore is an implementation of keys.ChainStore for testing.
// The zero value is usable.
type FakeChainStore struct {
	Addresses

	TxSigner
	MessageSigner

	internal.Locker[keys.Mutex]

	//TODO round robin state
} //TODO wrapper methods with enabled checks?

func (s *FakeChainStore) ForChain(cid *big.Int) keys.ChainStore {
	return s
}

var _ keys.Addresses = Addresses([]common.Address{})

type Addresses []common.Address

func (a Addresses) CheckEnabled(ctx context.Context, address common.Address) error {
	if slices.Contains(a, address) {
		return nil
	}
	return errors.New("address not found")
}

func (a Addresses) EnabledAddresses(ctx context.Context) (addresses []common.Address, err error) {
	return a, nil
}

// TODO only on keystore? or enhanced at least?
func (a Addresses) GetNextAddress(ctx context.Context, addresses ...common.Address) (address common.Address, err error) {
	addresses = slices.DeleteFunc(addresses, func(address common.Address) bool {
		return !slices.Contains(a, address)
	})
	return addresses[rand.Int()%len(addresses)], nil
}

var _ keys.TxSigner = TxSigner(nil)

type TxSigner func(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error)

func (f TxSigner) SignTx(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
	if f == nil {
		return tx, nil
	}
	return f(ctx, from, tx)
}

var _ keys.MessageSigner = MessageSigner(nil)

type MessageSigner func(ctx context.Context, address common.Address, message []byte) ([]byte, error)

func (f MessageSigner) SignMessage(ctx context.Context, address common.Address, message []byte) ([]byte, error) {
	if f == nil {
		return message, nil
	}
	return f(ctx, address, message)
}
