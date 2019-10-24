package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gitlab.com/thorchain/bepswap/thornode/common"
	. "gopkg.in/check.v1"
)

type TypeTxInSuite struct{}

var _ = Suite(&TypeTxInSuite{})

func (s TypeTxInSuite) TestVoter(c *C) {
	txID := GetRandomTxHash()
	txID2 := GetRandomTxHash()
	bnb := GetRandomBNBAddress()
	acc1 := GetRandomBech32Addr()
	acc2 := GetRandomBech32Addr()
	acc3 := GetRandomBech32Addr()
	acc4 := GetRandomBech32Addr()
	accConsPub1 := GetRandomBech32ConsensusPubKey()
	accConsPub2 := GetRandomBech32ConsensusPubKey()
	accConsPub3 := GetRandomBech32ConsensusPubKey()
	accConsPub4 := GetRandomBech32ConsensusPubKey()

	observePoolAddr := GetRandomBNBAddress()
	voter := NewTxInVoter(txID, nil)

	txIn := NewTxIn(nil, "hello", bnb, sdk.ZeroUint(), observePoolAddr)
	txIn2 := NewTxIn(nil, "goodbye", bnb, sdk.ZeroUint(), observePoolAddr)

	voter.Adds([]TxIn{txIn}, acc1)
	c.Assert(voter.Txs, HasLen, 1)

	voter.Adds([]TxIn{txIn}, acc1) // check we don't duplicate the same signer
	c.Assert(voter.Txs, HasLen, 1)
	c.Assert(voter.Txs[0].Signers, HasLen, 1)

	voter.Add(txIn, acc2) // append a signature
	c.Assert(voter.Txs, HasLen, 1)
	c.Assert(voter.Txs[0].Signers, HasLen, 2)

	voter.Add(txIn2, acc1) // same validator seeing a different version of tx
	c.Assert(voter.Txs, HasLen, 1)
	c.Assert(voter.Txs[0].Signers, HasLen, 2)

	voter.Add(txIn2, acc3) // second version
	c.Assert(voter.Txs, HasLen, 2)
	c.Assert(voter.Txs[0].Signers, HasLen, 2)
	c.Assert(voter.Txs[1].Signers, HasLen, 1)

	trusts3 := NodeAccounts{
		NodeAccount{
			NodeAddress: acc1,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc1,
				ValidatorBEPConsPubKey: accConsPub1,
			},
		},
		NodeAccount{
			NodeAddress: acc2,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc2,
				ValidatorBEPConsPubKey: accConsPub2,
			},
		},
		NodeAccount{
			NodeAddress: acc3,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc3,
				ValidatorBEPConsPubKey: accConsPub3,
			},
		},
	}
	trusts4 := NodeAccounts{
		NodeAccount{
			NodeAddress: acc1,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc1,
				ValidatorBEPConsPubKey: accConsPub1,
			},
		},
		NodeAccount{
			NodeAddress: acc2,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc2,
				ValidatorBEPConsPubKey: accConsPub2,
			},
		},
		NodeAccount{
			NodeAddress: acc3,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc3,
				ValidatorBEPConsPubKey: accConsPub3,
			},
		},
		NodeAccount{
			NodeAddress: acc4,
			Status:      Active,
			Accounts: TrustAccount{
				SignerBNBAddress:       bnb,
				ObserverBEPAddress:     acc4,
				ValidatorBEPConsPubKey: accConsPub4,
			},
		},
	}

	tx := voter.GetTx(trusts3)
	c.Check(tx.Memo, Equals, "hello")
	tx = voter.GetTx(trusts4)
	c.Check(tx.Empty(), Equals, true)
	c.Check(voter.HasConensus(trusts3), Equals, true)
	c.Check(voter.HasConensus(trusts4), Equals, false)
	c.Check(voter.Key().Equals(txID), Equals, true)
	c.Check(voter.String() == txID.String(), Equals, true)
	voter.SetDone(txID2)
	for _, transaction := range voter.Txs {
		c.Check(transaction.Done.Equals(txID2), Equals, true)
	}

	txIn.SetReverted(txID2)
	c.Check(txIn.Done.Equals(txID2), Equals, true)
	c.Check(len(txIn.String()) > 0, Equals, true)
	statechainCoins := common.Coins{
		common.NewCoin(common.RuneA1FAsset, sdk.NewUint(100)),
		common.NewCoin(common.BNBAsset, sdk.NewUint(100)),
	}
	inputs := []struct {
		coins           common.Coins
		memo            string
		sender          common.Address
		observePoolAddr common.Address
	}{
		{
			coins:           nil,
			memo:            "test",
			sender:          bnb,
			observePoolAddr: observePoolAddr,
		},
		{
			coins:           common.Coins{},
			memo:            "test",
			sender:          bnb,
			observePoolAddr: observePoolAddr,
		},
		{
			coins:           statechainCoins,
			memo:            "",
			sender:          bnb,
			observePoolAddr: observePoolAddr,
		},
		{
			coins:           statechainCoins,
			memo:            "test",
			sender:          common.NoAddress,
			observePoolAddr: observePoolAddr,
		},
		{
			coins:           statechainCoins,
			memo:            "test",
			sender:          bnb,
			observePoolAddr: common.NoAddress,
		},
	}

	for _, item := range inputs {
		txIn := NewTxIn(item.coins, item.memo, item.sender, sdk.ZeroUint(), item.observePoolAddr)
		c.Assert(txIn.Valid(), NotNil)
	}
}

func (TypeTxInSuite) TestTxInEquals(c *C) {
	coins1 := common.Coins{
		common.NewCoin(common.BNBAsset, sdk.NewUint(100*common.One)),
		common.NewCoin(common.RuneA1FAsset, sdk.NewUint(100*common.One)),
	}
	coins2 := common.Coins{
		common.NewCoin(common.BNBAsset, sdk.NewUint(100*common.One)),
	}
	coins3 := common.Coins{
		common.NewCoin(common.BNBAsset, sdk.NewUint(200*common.One)),
		common.NewCoin(common.RuneA1FAsset, sdk.NewUint(100*common.One)),
	}
	coins4 := common.Coins{
		common.NewCoin(common.RuneB1AAsset, sdk.NewUint(100*common.One)),
		common.NewCoin(common.RuneA1FAsset, sdk.NewUint(100*common.One)),
	}
	bnb, err := common.NewAddress("bnb1xlvns0n2mxh77mzaspn2hgav4rr4m8eerfju38")
	c.Assert(err, IsNil)
	bnb1, err := common.NewAddress("bnb1yk882gllgv3rt2rqrsudf6kn2agr94etnxu9a7")
	c.Assert(err, IsNil)
	observePoolAddr, err := common.NewAddress("bnb1g0xakzh03tpa54khxyvheeu92hwzypkdce77rm")
	c.Assert(err, IsNil)
	observePoolAddr1, err := common.NewAddress("bnb1zxseqkfm3en5cw6dh9xgmr85hw6jtwamnd2y2v")
	c.Assert(err, IsNil)
	inputs := []struct {
		tx    TxIn
		tx1   TxIn
		equal bool
	}{
		{
			tx:    NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo1", bnb, sdk.ZeroUint(), observePoolAddr),
			equal: false,
		},
		{
			tx:    NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo", bnb1, sdk.ZeroUint(), observePoolAddr),
			equal: false,
		},
		{
			tx:    NewTxIn(coins2, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			equal: false,
		},
		{
			tx:    NewTxIn(coins3, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			equal: false,
		},
		{
			tx:    NewTxIn(coins4, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			equal: false,
		},
		{
			tx:    NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr),
			tx1:   NewTxIn(coins1, "memo", bnb, sdk.ZeroUint(), observePoolAddr1),
			equal: false,
		},
	}
	for _, item := range inputs {
		c.Assert(item.tx.Equals(item.tx1), Equals, item.equal)
	}
}
