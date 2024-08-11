package common

import (
	. "gopkg.in/check.v1"
)

type AssetSuite struct{}

var _ = Suite(&AssetSuite{})

func (s AssetSuite) TestAsset(c *C) {
	asset, err := NewAsset("bnb.rune-a1f")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(RuneA1FAsset), Equals, true)
	c.Check(asset.IsRune(), Equals, true)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "BNB.RUNE-A1F")

	c.Check(asset.Chain.Equals(BNBChain), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("RUNE-A1F")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("RUNE")), Equals, true)

	// parse without chain
	asset, err = NewAsset("rune")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(RuneNative), Equals, true)

	// ETH test
	asset, err = NewAsset("eth.knc")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(ETHChain), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("KNC")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("KNC")), Equals, true)
}
