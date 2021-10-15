package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/irisnet/irismod/modules/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(c context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var supply uint64
	switch {
	case len(request.Owner) == 0 && len(request.ClassId) > 0:
		supply = k.GetTotalSupply(ctx, request.ClassId)
	default:
		owner, err := sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address %s", request.Owner)
		}
		supply = k.GetTotalSupplyOfOwner(ctx, request.ClassId, owner)
	}
	return &types.QuerySupplyResponse{Amount: supply}, nil
}

func (k Keeper) Owner(c context.Context, request *types.QueryOwnerRequest) (*types.QueryOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddress, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner address %s", request.Owner)
	}

	owner := types.Owner{
		Address:       ownerAddress.String(),
		IDCollections: types.IDCollections{},
	}
	idsMap := make(map[string][]string)
	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyOwner(ownerAddress, request.ClassId, ""))
	pageRes, err := query.Paginate(nftStore, request.Pagination, func(key []byte, value []byte) error {
		classID := request.ClassId
		tokenID := string(key)
		if len(request.ClassId) == 0 {
			classID, tokenID, _ = types.SplitKeyClass(key)
		}
		if ids, ok := idsMap[classID]; ok {
			idsMap[classID] = append(ids, tokenID)
		} else {
			idsMap[classID] = []string{tokenID}
			owner.IDCollections = append(
				owner.IDCollections,
				types.IDCollection{ClassId: classID},
			)
		}
		return nil
	})
	for i := 0; i < len(owner.IDCollections); i++ {
		owner.IDCollections[i].TokenIds = idsMap[owner.IDCollections[i].ClassId]
	}
	return &types.QueryOwnerResponse{Owner: &owner, Pagination: pageRes}, nil
}

func (k Keeper) Collection(c context.Context, request *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collection, pageRes, err := k.GetPaginateCollection(ctx, request, request.ClassId)
	if err != nil {
		return nil, err
	}
	return &types.QueryCollectionResponse{Collection: &collection, Pagination: pageRes}, nil
}

func (k Keeper) Class(c context.Context, request *types.QueryClassRequest) (*types.QueryClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	classObject, found := k.GetClass(ctx, request.ClassId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", request.ClassId)
	}

	return &types.QueryClassResponse{Class: &classObject}, nil
}

func (k Keeper) Classes(c context.Context, req *types.QueryClassesRequest) (*types.QueryClassesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var classes []types.Class
	store := ctx.KVStore(k.storeKey)
	classStore := prefix.NewStore(store, types.KeyClassID(""))
	pageRes, err := query.Paginate(classStore, req.Pagination, func(key []byte, value []byte) error {
		var class types.Class
		k.cdc.MustUnmarshal(value, &class)
		classes = append(classes, class)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryClassesResponse{
		Classes:    classes,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) NFT(c context.Context, request *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, request.ClassId, request.TokenId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid NFT %s from collection %s", request.TokenId, request.ClassId)
	}

	NFT, ok := nft.(types.NFT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid type NFT %s from collection %s", request.TokenId, request.ClassId)
	}

	return &types.QueryNFTResponse{NFT: &NFT}, nil
}
