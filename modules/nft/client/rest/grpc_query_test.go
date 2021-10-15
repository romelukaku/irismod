package rest_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
	nfttestutil "github.com/irisnet/irismod/modules/nft/client/testutil"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 2

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestNft() {
	val := s.network.Validators[0]
	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	// ---------------------------------------------------------------------------

	from := val.Address
	tokenName := "Kitty Token"
	tokenURI := "uri"
	tokenData := "data"
	tokenID := "kitty"
	//owner     := "owner"
	className := "name"
	class := "class"
	schema := "schema"
	symbol := "symbol"
	mintRestricted := true
	updateRestricted := false
	baseURL := val.APIAddress

	//------test GetCmdIssueClass()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagClassName, className),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%t", nftcli.FlagMintRestricted, mintRestricted),
		fmt.Sprintf("--%s=%t", nftcli.FlagUpdateRestricted, updateRestricted),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := nfttestutil.IssueClassExec(val.ClientCtx, from.String(), class, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	classID := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	//------test GetCmdQueryClass()-------------
	url := fmt.Sprintf("%s/irismod/nft/classes/%s", baseURL, classID)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryClassResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	classItem := respType.(*nfttypes.QueryClassResponse)
	s.Require().Equal(className, classItem.Class.Name)
	s.Require().Equal(schema, classItem.Class.Schema)
	s.Require().Equal(symbol, classItem.Class.Symbol)
	s.Require().Equal(mintRestricted, classItem.Class.MintRestricted)
	s.Require().Equal(updateRestricted, classItem.Class.UpdateRestricted)

	//------test GetCmdQueryClasses()-------------
	url = fmt.Sprintf("%s/irismod/nft/classes", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryClassesResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	classesResp := respType.(*nfttypes.QueryClassesResponse)
	s.Require().Equal(1, len(classesResp.Classes))
	s.Require().Equal(classID, classesResp.Classes[0].Id)

	//------test GetCmdMintNFT()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.MintNFTExec(val.ClientCtx, from.String(), classID, tokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQuerySupply()-------------
	url = fmt.Sprintf("%s/irismod/nft/collections/%s/supply", baseURL, classID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	supplyResp := respType.(*nfttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(1), supplyResp.Amount)

	//------test GetCmdQueryNFT()-------------
	url = fmt.Sprintf("%s/irismod/nft/nfts/%s/%s", baseURL, classID, tokenID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryNFTResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	nftItem := respType.(*nfttypes.QueryNFTResponse)
	s.Require().Equal(tokenID, nftItem.NFT.Id)
	s.Require().Equal(tokenName, nftItem.NFT.Name)
	s.Require().Equal(tokenURI, nftItem.NFT.URI)
	s.Require().Equal(tokenData, nftItem.NFT.Data)
	s.Require().Equal(from.String(), nftItem.NFT.Owner)

	//------test GetCmdQueryOwner()-------------
	url = fmt.Sprintf("%s/irismod/nft/nfts?owner=%s", baseURL, from.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryOwnerResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	ownerResp := respType.(*nfttypes.QueryOwnerResponse)
	s.Require().Equal(from.String(), ownerResp.Owner.Address)
	s.Require().Equal(class, ownerResp.Owner.IDCollections[0].ClassId)
	s.Require().Equal(tokenID, ownerResp.Owner.IDCollections[0].TokenIds[0])

	//------test GetCmdQueryCollection()-------------
	url = fmt.Sprintf("%s/irismod/nft/collections/%s", baseURL, classID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryCollectionResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	collectionResp := respType.(*nfttypes.QueryCollectionResponse)
	s.Require().Equal(1, len(collectionResp.Collection.NFTs))

	//------test GetCmdTransferClass()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.TransferClassExec(val.ClientCtx, from.String(), recipient.String(), classID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.Class{})
	bz, err = nfttestutil.QueryClassExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	classItem2 := respType.(*nfttypes.Class)
	s.Require().Equal(recipient.String(), classItem2.Creator)
	s.Require().Equal(className, classItem2.Name)
	s.Require().Equal(schema, classItem2.Schema)
	s.Require().Equal(symbol, classItem2.Symbol)
	s.Require().Equal(mintRestricted, classItem2.MintRestricted)
	s.Require().Equal(updateRestricted, classItem2.UpdateRestricted)
}
