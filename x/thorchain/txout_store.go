package thorchain

import (
	"github.com/blang/semver"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/thorchain/thornode/constants"
)

type VersionedTxOutStore interface {
	GetTxOutStore(keeper Keeper, version semver.Version) (TxOutStore, error)
}

type TxOutStore interface {
	NewBlock(height int64, constAccessor constants.ConstantValues)
	CommitBlock(ctx sdk.Context)
	GetBlockOut() *TxOut
	ClearOutboundItems()
	GetOutboundItems() []*TxOutItem
	TryAddTxOutItem(ctx sdk.Context, toi *TxOutItem) (bool, error)
}

type VersionedTxOutStorage struct {
	txOutStorage *TxOutStorageV1
}

// NewVersionedTxOutStore create a new instance of VersionedTxOutStorage
func NewVersionedTxOutStore() *VersionedTxOutStorage {
	return &VersionedTxOutStorage{}
}

// GetTxOutStore will return an implementation of the txout store that
func (s *VersionedTxOutStorage) GetTxOutStore(keeper Keeper, version semver.Version) (TxOutStore, error) {
	if version.GTE(semver.MustParse("0.1.0")) {
		if s.txOutStorage == nil {
			s.txOutStorage = NewTxOutStorageV1(keeper)
		}
		return s.txOutStorage, nil
	}
	return nil, errInvalidVersion
}
