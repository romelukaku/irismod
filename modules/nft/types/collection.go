package types

import (
	"github.com/irisnet/irismod/modules/nft/exported"
)

// NewCollection creates a new NFT Collection
func NewCollection(class Class, nfts []exported.NFT) (c Collection) {
	c.Class = class
	for _, nft := range nfts {
		c = c.AddNFT(nft.(NFT))
	}
	return c
}

// AddNFT adds an NFT to the collection
func (c Collection) AddNFT(nft NFT) Collection {
	c.NFTs = append(c.NFTs, nft)
	return c
}

func (c Collection) Supply() int {
	return len(c.NFTs)
}

// NewCollection creates a new NFT Collection
func NewCollections(c ...Collection) []Collection {
	return append([]Collection{}, c...)
}
