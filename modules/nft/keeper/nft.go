package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/exported"
	"github.com/irisnet/irismod/modules/nft/types"
)

// GetNFT gets the the specified NFT
func (k Keeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nft exported.NFT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNFT(classID, tokenID))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownCollection, "not found NFT: %s", classID)
	}

	var NFT types.NFT
	k.cdc.MustUnmarshal(bz, &NFT)

	return NFT, nil
}

// GetNFTs returns all NFTs by the specified class ID
func (k Keeper) GetNFTs(ctx sdk.Context, class string) (nfts []exported.NFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyNFT(class, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var NFT types.NFT
		k.cdc.MustUnmarshal(iterator.Value(), &NFT)
		nfts = append(nfts, NFT)
	}

	return nfts
}

// Authorize checks if the sender is the owner of the given NFT
// Return the NFT if true, an error otherwise
func (k Keeper) Authorize(ctx sdk.Context, classID, tokenID string, owner sdk.AccAddress) (types.NFT, error) {
	nft, err := k.GetNFT(ctx, classID, tokenID)
	if err != nil {
		return types.NFT{}, err
	}

	if !owner.Equals(nft.GetOwner()) {
		return types.NFT{}, sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}

	return nft.(types.NFT), nil
}

// HasNFT checks if the specified NFT exists
func (k Keeper) HasNFT(ctx sdk.Context, classID, tokenID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyNFT(classID, tokenID))
}

func (k Keeper) setNFT(ctx sdk.Context, classID string, nft types.NFT) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&nft)
	store.Set(types.KeyNFT(classID, nft.GetID()), bz)
}

// deleteNFT deletes an existing NFT from store
func (k Keeper) deleteNFT(ctx sdk.Context, classID string, nft exported.NFT) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNFT(classID, nft.GetID()))
}
