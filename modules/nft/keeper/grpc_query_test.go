package keeper_test

import (
	gocontext "context"

	"github.com/irisnet/irismod/modules/nft/types"
)

func (suite *KeeperSuite) TestSupply() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.Supply(gocontext.Background(), &types.QuerySupplyRequest{
		ClassId: classID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.Equal(1, int(response.Amount))
}

func (suite *KeeperSuite) TestOwner() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.Owner(gocontext.Background(), &types.QueryOwnerRequest{
		ClassId: classID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.NotNil(response.Owner)
	suite.Contains(response.Owner.IDCollections[0].TokenIds, tokenID)
}

func (suite *KeeperSuite) TestCollection() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.Collection(gocontext.Background(), &types.QueryCollectionRequest{
		ClassId: classID,
	})

	suite.NoError(err)
	suite.NotNil(response.Collection)
	suite.Len(response.Collection.NFTs, 1)
	suite.Equal(response.Collection.NFTs[0].Id, tokenID)
}

func (suite *KeeperSuite) TestClass() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.Class(gocontext.Background(), &types.QueryClassRequest{
		ClassId: classID,
	})

	suite.NoError(err)
	suite.NotNil(response.Class)
	suite.Equal(response.Class.Id, classID)
}

func (suite *KeeperSuite) TestClasses() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.Classes(gocontext.Background(), &types.QueryClassesRequest{})

	suite.NoError(err)
	suite.NotEmpty(response.Classes)
	suite.Equal(response.Classes[0].Id, classID)
}

func (suite *KeeperSuite) TestNFT() {
	err := suite.keeper.MintNFT(suite.ctx, classID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	response, err := suite.queryClient.NFT(gocontext.Background(), &types.QueryNFTRequest{
		ClassId: classID,
		TokenId: tokenID,
	})

	suite.NoError(err)
	suite.NotEmpty(response.NFT)
	suite.Equal(response.NFT.Id, tokenID)
}
