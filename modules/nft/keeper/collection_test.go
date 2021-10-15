package keeper_test

import (
	"github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"
)

func (suite *KeeperSuite) TestSetCollection() {
	nft := types.NewNFT(tokenID, tokenNm, address, tokenURI, tokenData)
	// create a new NFT and add it to the collection created with the NFT mint
	nft2 := types.NewNFT(tokenID2, tokenNm, address, tokenURI, tokenData)

	classE := types.Class{
		Id:               classID,
		Name:             classNm,
		Schema:           schema,
		Creator:          address.String(),
		Symbol:           classSymbol,
		MintRestricted:   true,
		UpdateRestricted: true,
	}

	collection2 := types.Collection{
		Class: classE,
		NFTs:  []types.NFT{nft2, nft},
	}

	err := suite.keeper.SetCollection(suite.ctx, collection2)
	suite.Nil(err)

	collection2, err = suite.keeper.GetCollection(suite.ctx, classID)
	suite.NoError(err)
	suite.Len(collection2.NFTs, 2)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetCollection() {
	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// collection should exist
	collection, err := suite.keeper.GetCollection(suite.ctx, classID)
	suite.NoError(err)
	suite.NotEmpty(collection)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetCollections() {

	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetSupply() {
	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.MintNFT(suite.ctx, classID, tokenID2, tokenNm2, tokenURI, tokenData, address2)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.MintNFT(suite.ctx, classID2, tokenID, tokenNm2, tokenURI, tokenData, address2)
	suite.NoError(err)

	supply := suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupplyOfOwner(suite.ctx, classID, address)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupplyOfOwner(suite.ctx, classID, address2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID2)
	suite.Equal(uint64(1), supply)

	//burn nft
	err = suite.keeper.BurnNFT(suite.ctx, classID, tokenID, address)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(1), supply)

	//burn nft
	err = suite.keeper.BurnNFT(suite.ctx, classID, tokenID2, address2)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(0), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, classID)
	suite.Equal(uint64(0), supply)
}
