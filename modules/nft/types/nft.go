package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/exported"
)

var _ exported.NFT = NFT{}

// NewNFT creates a new NFT instance
func NewNFT(id, name string, owner sdk.AccAddress, uri, data string) NFT {
	return NFT{
		Id:    id,
		Name:  name,
		Owner: owner.String(),
		URI:   uri,
		Data:  data,
	}
}

// GetID return the id of NFT
func (bnft NFT) GetID() string {
	return bnft.Id
}

// GetName return the name of NFT
func (bnft NFT) GetName() string {
	return bnft.Name
}

// GetOwner return the owner of NFT
func (bnft NFT) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(bnft.Owner)
	return owner
}

// GetURI return the URI of NFT
func (bnft NFT) GetURI() string {
	return bnft.URI
}

// GetData return the Data of NFT
func (bnft NFT) GetData() string {
	return bnft.Data
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts)
}
