package cli_test

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
	val2 := s.network.Validators[1]

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

	//------test GetCmdIssueClass()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagClassName, className),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, symbol),
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
	respType = proto.Message(&nfttypes.Class{})
	bz, err = nfttestutil.QueryClassExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	classItem := respType.(*nfttypes.Class)
	s.Require().Equal(className, classItem.Name)
	s.Require().Equal(schema, classItem.Schema)
	s.Require().Equal(symbol, classItem.Symbol)
	s.Require().Equal(mintRestricted, classItem.MintRestricted)
	s.Require().Equal(updateRestricted, classItem.UpdateRestricted)

	//------test GetCmdQueryClasses()-------------
	respType = proto.Message(&nfttypes.QueryClassesResponse{})
	bz, err = nfttestutil.QueryClassesExec(val.ClientCtx)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
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
	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
	bz, err = nfttestutil.QuerySupplyExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supplyResp := respType.(*nfttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(1), supplyResp.Amount)

	//------test GetCmdQueryNFT()-------------
	respType = proto.Message(&nfttypes.NFT{})
	bz, err = nfttestutil.QueryNFTExec(val.ClientCtx, classID, tokenID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	nftItem := respType.(*nfttypes.NFT)
	s.Require().Equal(tokenID, nftItem.Id)
	s.Require().Equal(tokenName, nftItem.Name)
	s.Require().Equal(tokenURI, nftItem.URI)
	s.Require().Equal(tokenData, nftItem.Data)
	s.Require().Equal(from.String(), nftItem.Owner)

	//------test GetCmdQueryOwner()-------------
	respType = proto.Message(&nfttypes.QueryOwnerResponse{})
	bz, err = nfttestutil.QueryOwnerExec(val.ClientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	ownerResp := respType.(*nfttypes.QueryOwnerResponse)
	s.Require().Equal(from.String(), ownerResp.Owner.Address)
	s.Require().Equal(class, ownerResp.Owner.IDCollections[0].ClassId)
	s.Require().Equal(tokenID, ownerResp.Owner.IDCollections[0].TokenIds[0])

	//------test GetCmdQueryCollection()-------------
	respType = proto.Message(&nfttypes.QueryCollectionResponse{})
	bz, err = nfttestutil.QueryCollectionExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	collectionItem := respType.(*nfttypes.QueryCollectionResponse)
	s.Require().Equal(1, len(collectionItem.Collection.NFTs))

	//------test GetCmdEditNFT()-------------
	newTokenDate := "newdata"
	newTokenURI := "newuri"
	newTokenName := "new Kitty Token"
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, newTokenDate),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, newTokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, newTokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.EditNFTExec(val.ClientCtx, from.String(), classID, tokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.NFT{})
	bz, err = nfttestutil.QueryNFTExec(val.ClientCtx, classID, tokenID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	newNftItem := respType.(*nfttypes.NFT)
	s.Require().Equal(newTokenName, newNftItem.Name)
	s.Require().Equal(newTokenURI, newNftItem.URI)
	s.Require().Equal(newTokenDate, newNftItem.Data)

	//------test GetCmdTransferNFT()-------------
	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.TransferNFTExec(val.ClientCtx, from.String(), recipient.String(), classID, tokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.NFT{})
	bz, err = nfttestutil.QueryNFTExec(val.ClientCtx, classID, tokenID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	nftItem = respType.(*nfttypes.NFT)
	s.Require().Equal(tokenID, nftItem.Id)
	s.Require().Equal(tokenName, nftItem.Name)
	s.Require().Equal(tokenURI, nftItem.URI)
	s.Require().Equal(tokenData, nftItem.Data)
	s.Require().Equal(recipient.String(), nftItem.Owner)

	//------test GetCmdBurnNFT()-------------
	newTokenID := "dgsbl"
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, newTokenDate),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, newTokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, newTokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.MintNFTExec(val.ClientCtx, from.String(), classID, newTokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
	bz, err = nfttestutil.QuerySupplyExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supplyResp = respType.(*nfttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(2), supplyResp.Amount)

	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	bz, err = nfttestutil.BurnNFTExec(val.ClientCtx, from.String(), classID, newTokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val2.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
	bz, err = nfttestutil.QuerySupplyExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supplyResp = respType.(*nfttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(1), supplyResp.Amount)

	//------test GetCmdTransferClass()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.TransferClassExec(val.ClientCtx, from.String(), val2.Address.String(), classID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.Class{})
	bz, err = nfttestutil.QueryClassExec(val.ClientCtx, classID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	classItem2 := respType.(*nfttypes.Class)
	s.Require().Equal(val2.Address.String(), classItem2.Creator)
	s.Require().Equal(className, classItem2.Name)
	s.Require().Equal(schema, classItem2.Schema)
	s.Require().Equal(symbol, classItem2.Symbol)
	s.Require().Equal(mintRestricted, classItem2.MintRestricted)
	s.Require().Equal(updateRestricted, classItem2.UpdateRestricted)
}
