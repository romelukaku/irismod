package keeper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/irisnet/irismod/modules/nft/exported"
	"github.com/irisnet/irismod/modules/nft/types"
)

// SetCollection saves all NFTs and returns an error if there already exists
func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	for _, nft := range collection.NFTs {
		if err := k.MintNFT(
			ctx,
			collection.Class.Id,
			nft.GetID(),
			nft.GetName(),
			nft.GetURI(),
			nft.GetData(),
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

// GetCollection returns the collection by the specified class ID
func (k Keeper) GetCollection(ctx sdk.Context, classID string) (types.Collection, error) {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidClass, "classID %s not existed ", classID)
	}

	nfts := k.GetNFTs(ctx, classID)
	return types.NewCollection(class, nfts), nil
}

// GetPaginateCollection returns the collection by the specified class ID
func (k Keeper) GetPaginateCollection(ctx sdk.Context, request *types.QueryCollectionRequest, classID string) (types.Collection, *query.PageResponse, error) {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return types.Collection{}, nil, sdkerrors.Wrapf(types.ErrInvalidClass, "classID %s not existed ", classID)
	}
	var nfts []exported.NFT
	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyNFT(classID, ""))
	pageRes, err := query.Paginate(nftStore, request.Pagination, func(key []byte, value []byte) error {
		var NFT types.NFT
		k.cdc.MustUnmarshal(value, &NFT)
		nfts = append(nfts, NFT)
		return nil
	})
	if err != nil {
		return types.Collection{}, nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return types.NewCollection(class, nfts), pageRes, nil
}

// GetCollections returns all the collections
func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection) {
	for _, class := range k.GetClasses(ctx) {
		nfts := k.GetNFTs(ctx, class.Id)
		cs = append(cs, types.NewCollection(class, nfts))
	}
	return cs
}

// GetTotalSupply returns the number of NFTs by the specified class ID
func (k Keeper) GetTotalSupply(ctx sdk.Context, classID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollection(classID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

// GetTotalSupplyOfOwner returns the amount of NFTs by the specified conditions
func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(owner, id, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		supply++
	}
	return supply
}

func (k Keeper) increaseSupply(ctx sdk.Context, classID string) {
	supply := k.GetTotalSupply(ctx, classID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(classID), bz)
}

func (k Keeper) decreaseSupply(ctx sdk.Context, classID string) {
	supply := k.GetTotalSupply(ctx, classID)
	supply--

	store := ctx.KVStore(k.storeKey)
	if supply == 0 {
		store.Delete(types.KeyCollection(classID))
		return
	}

	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(classID), bz)
}
