package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewClass return a new class
func NewClass(id, name, schema, symbol string, creator sdk.AccAddress, mintRestricted, updateRestricted bool) Class {
	return Class{
		Id:               id,
		Name:             name,
		Schema:           schema,
		Creator:          creator.String(),
		Symbol:           symbol,
		MintRestricted:   mintRestricted,
		UpdateRestricted: updateRestricted,
	}
}
