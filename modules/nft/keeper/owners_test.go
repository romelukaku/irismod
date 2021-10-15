package keeper_test

import (
	"github.com/irisnet/irismod/modules/nft/keeper"
)

func (suite *KeeperSuite) TestGetOwners() {

	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, classID, tokenID2, tokenNm2, tokenURI, tokenData, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, classID, tokenID3, tokenNm3, tokenURI, tokenData, address3)
	suite.NoError(err)

	owners := suite.keeper.GetOwners(suite.ctx)
	suite.Equal(3, len(owners))

	err = suite.keeper.MintNFT(suite.ctx, classID2, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, classID2, tokenID2, tokenNm2, tokenURI, tokenData, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, classID2, tokenID3, tokenNm3, tokenURI, tokenData, address3)
	suite.NoError(err)

	owners = suite.keeper.GetOwners(suite.ctx)
	suite.Equal(3, len(owners))

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}
