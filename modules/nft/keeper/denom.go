package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/types"
)

// HasClassID returns whether the specified class ID exists
func (k Keeper) HasClassID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyClassID(id))
}

// SetClass is responsible for saving the definition of class
func (k Keeper) SetClass(ctx sdk.Context, class types.Class) error {
	if k.HasClassID(ctx, class.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "classID %s has already exists", class.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&class)
	store.Set(types.KeyClassID(class.Id), bz)
	store.Set(types.KeyClassName(class.Name), []byte(class.Id))
	return nil
}

// GetClass returns the class by id
func (k Keeper) GetClass(ctx sdk.Context, id string) (class types.Class, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyClassID(id))
	if len(bz) == 0 {
		return class, false
	}

	k.cdc.MustUnmarshal(bz, &class)
	return class, true
}

// GetClasses returns all the classes
func (k Keeper) GetClasses(ctx sdk.Context) (classes []types.Class) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyClassID(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		k.cdc.MustUnmarshal(iterator.Value(), &class)
		classes = append(classes, class)
	}
	return classes
}

// UpdateClass is responsible for updating the definition of class
func (k Keeper) UpdateClass(ctx sdk.Context, class types.Class) error {
	if !k.HasClassID(ctx, class.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidClass, "classID %s not exists", class.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&class)
	store.Set(types.KeyClassID(class.Id), bz)
	return nil
}
