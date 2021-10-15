package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      codec.Codec
}

// NewKeeper creates a new instance of the NFT Keeper
func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// IssueClass issues a class according to the given params
func (k Keeper) IssueClass(ctx sdk.Context,
	id, name, schema, symbol string,
	creator sdk.AccAddress,
	mintRestricted, updateRestricted bool,
) error {
	return k.SetClass(ctx, types.NewClass(id, name, schema, symbol, creator, mintRestricted, updateRestricted))
}

// MintNFT mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFT(
	ctx sdk.Context, classID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", classID)
	}

	if class.MintRestricted && class.Creator != owner.String() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to mint NFT of class %s", class.Creator, classID)
	}

	if k.HasNFT(ctx, classID, tokenID) {
		return sdkerrors.Wrapf(types.ErrNFTAlreadyExists, "NFT %s already exists in collection %s", tokenID, classID)
	}

	k.setNFT(
		ctx, classID,
		types.NewNFT(
			tokenID,
			tokenNm,
			owner,
			tokenURI,
			tokenData,
		),
	)
	k.setOwner(ctx, classID, tokenID, owner)
	k.increaseSupply(ctx, classID)

	return nil
}

// EditNFT updates an already existing NFT
func (k Keeper) EditNFT(
	ctx sdk.Context, classID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", classID)
	}

	if class.UpdateRestricted {
		// if true , nobody can update the NFT under this class
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "nobody can update the NFT under this class %s", class.Id)
	}

	// just the owner of NFT can edit
	nft, err := k.Authorize(ctx, classID, tokenID, owner)
	if err != nil {
		return err
	}

	if types.Modified(tokenNm) {
		nft.Name = tokenNm
	}

	if types.Modified(tokenURI) {
		nft.URI = tokenURI
	}

	if types.Modified(tokenData) {
		nft.Data = tokenData
	}

	k.setNFT(ctx, classID, nft)

	return nil
}

// TransferOwner transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferOwner(
	ctx sdk.Context, classID, tokenID, tokenNm, tokenURI,
	tokenData string, srcOwner, dstOwner sdk.AccAddress,
) error {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", classID)
	}

	nft, err := k.Authorize(ctx, classID, tokenID, srcOwner)
	if err != nil {
		return err
	}

	nft.Owner = dstOwner.String()

	if class.UpdateRestricted && (types.Modified(tokenNm) || types.Modified(tokenURI) || types.Modified(tokenData)) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "It is restricted to update NFT under this class %s", class.Id)
	}

	if types.Modified(tokenNm) {
		nft.Name = tokenNm
	}
	if types.Modified(tokenURI) {
		nft.URI = tokenURI
	}
	if types.Modified(tokenData) {
		nft.Data = tokenData
	}

	k.setNFT(ctx, classID, nft)
	k.swapOwner(ctx, classID, tokenID, srcOwner, dstOwner)
	return nil
}

// BurnNFT deletes a specified NFT
func (k Keeper) BurnNFT(ctx sdk.Context, classID, tokenID string, owner sdk.AccAddress) error {
	if !k.HasClassID(ctx, classID) {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", classID)
	}

	nft, err := k.Authorize(ctx, classID, tokenID, owner)
	if err != nil {
		return err
	}

	k.deleteNFT(ctx, classID, nft)
	k.deleteOwner(ctx, classID, tokenID, owner)
	k.decreaseSupply(ctx, classID)

	return nil
}

// TransferClassOwner transfers the ownership of the given class to the new owner
func (k Keeper) TransferClassOwner(
	ctx sdk.Context, classID string, srcOwner, dstOwner sdk.AccAddress,
) error {
	class, found := k.GetClass(ctx, classID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "class ID %s not exists", classID)
	}

	// authorize
	if srcOwner.String() != class.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to transfer class %s", srcOwner.String(), classID)
	}

	class.Creator = dstOwner.String()

	err := k.UpdateClass(ctx, class)
	if err != nil {
		return err
	}

	return nil
}
