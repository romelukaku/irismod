package types

// DONTCOVER

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the NFT Querier
const (
	QuerySupply     = "supply"
	QueryOwner      = "owner"
	QueryCollection = "collection"
	QueryClasses    = "classes"
	QueryClass      = "class"
	QueryNFT        = "nft"
)

// QuerySupplyParams defines the params for queries:
type QuerySupplyParams struct {
	Class string
	Owner sdk.AccAddress
}

// NewQuerySupplyParams creates a new instance of QuerySupplyParams
func NewQuerySupplyParams(class string, owner sdk.AccAddress) QuerySupplyParams {
	return QuerySupplyParams{
		Class: class,
		Owner: owner,
	}
}

// Bytes exports the Class as bytes
func (q QuerySupplyParams) Bytes() []byte {
	return []byte(q.Class)
}

// QueryOwnerParams defines the params for queries:
type QueryOwnerParams struct {
	Class string
	Owner sdk.AccAddress
}

// NewQuerySupplyParams creates a new instance of QuerySupplyParams
func NewQueryOwnerParams(class string, owner sdk.AccAddress) QueryOwnerParams {
	return QueryOwnerParams{
		Class: class,
		Owner: owner,
	}
}

// QuerySupplyParams defines the params for queries:
type QueryCollectionParams struct {
	Class string
}

// NewQueryCollectionParams creates a new instance of QueryCollectionParams
func NewQueryCollectionParams(class string) QueryCollectionParams {
	return QueryCollectionParams{
		Class: class,
	}
}

// QueryClassParams defines the params for queries:
type QueryClassParams struct {
	ID string
}

// NewQueryClassParams creates a new instance of QueryClassParams
func NewQueryClassParams(id string) QueryClassParams {
	return QueryClassParams{
		ID: id,
	}
}

// QueryNFTParams params for query 'custom/nfts/nft'
type QueryNFTParams struct {
	Class   string
	TokenID string
}

// NewQueryNFTParams creates a new instance of QueryNFTParams
func NewQueryNFTParams(class, id string) QueryNFTParams {
	return QueryNFTParams{
		Class:   class,
		TokenID: id,
	}
}
