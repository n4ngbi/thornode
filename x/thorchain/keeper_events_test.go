package thorchain

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
)

type KeeperEventsSuite struct{}

var _ = Suite(&KeeperEventsSuite{})

func (s *KeeperEventsSuite) TestEvents(c *C) {
	ctx, k := setupKeeperForTest(c)
	txID, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	inTx := common.NewTx(
		txID,
		GetRandomBNBAddress(),
		GetRandomBNBAddress(),
		common.Coins{
			common.NewCoin(common.BNBAsset, sdk.NewUint(320000000)),
			common.NewCoin(common.RuneAsset(), sdk.NewUint(420000000)),
		},
		BNBGasFeeSingleton,
		"SWAP:BNB.BNB",
	)
	swap := NewEventSwap(
		common.BNBAsset,
		sdk.NewUint(5),
		sdk.NewUint(5),
		sdk.NewUint(5),
		sdk.NewUint(5),
		inTx,
	)
	swapBytes, _ := json.Marshal(swap)
	evt := NewEvent(
		swap.Type(),
		12,
		inTx,
		swapBytes,
		EventSuccess,
	)

	c.Assert(k.UpsertEvent(ctx, evt), IsNil)
	e, err := k.GetEvent(ctx, 1)
	c.Assert(err, IsNil)
	c.Assert(e.Empty(), Equals, false)

	// add another event, and make sure both exists
	c.Assert(k.UpsertEvent(ctx, evt), IsNil)
	e, err = k.GetEvent(ctx, 2)
	c.Assert(err, IsNil)
	c.Assert(e.Empty(), Equals, false)

	// check txIn ID cant be empty
	evt.InTx.ID = ""
	c.Assert(k.UpsertEvent(ctx, evt), NotNil)
	e, err = k.GetEvent(ctx, 3)
	c.Assert(err, IsNil)
	c.Assert(e.Empty(), Equals, true)
}
