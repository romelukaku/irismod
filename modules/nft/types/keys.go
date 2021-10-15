package types

import (
	"bytes"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "nft"

	// StoreKey is the default store key for NFT
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the NFT store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the NFT module
	RouterKey = ModuleName
)

var (
	PrefixNFT        = []byte{0x01}
	PrefixOwners     = []byte{0x02} // key for a owner
	PrefixCollection = []byte{0x03} // key for balance of NFTs held by the class
	PrefixClass      = []byte{0x04} // key for class of the nft
	PrefixClassName  = []byte{0x05} // key for class name of the nft

	delimiter = []byte("/")
)

// SplitKeyOwner return the address,class,id from the key of stored owner
func SplitKeyOwner(key []byte) (address sdk.AccAddress, classID, tokenID string, err error) {
	key = key[len(PrefixOwners)+len(delimiter):]
	keys := bytes.Split(key, delimiter)
	if len(keys) != 3 {
		return address, classID, tokenID, errors.New("wrong KeyOwner")
	}

	address, _ = sdk.AccAddressFromBech32(string(keys[0]))
	classID = string(keys[1])
	tokenID = string(keys[2])
	return
}

func SplitKeyClass(key []byte) (classID, tokenID string, err error) {
	keys := bytes.Split(key, delimiter)
	if len(keys) != 2 {
		return classID, tokenID, errors.New("wrong KeyOwner")
	}

	classID = string(keys[0])
	tokenID = string(keys[1])
	return
}

// KeyOwner gets the key of a collection owned by an account address
func KeyOwner(address sdk.AccAddress, classID, tokenID string) []byte {
	key := append(PrefixOwners, delimiter...)
	if address != nil {
		key = append(key, []byte(address.String())...)
		key = append(key, delimiter...)
	}

	if address != nil && len(classID) > 0 {
		key = append(key, []byte(classID)...)
		key = append(key, delimiter...)
	}

	if address != nil && len(classID) > 0 && len(tokenID) > 0 {
		key = append(key, []byte(tokenID)...)
	}
	return key
}

// KeyNFT gets the key of nft stored by an class and id
func KeyNFT(classID, tokenID string) []byte {
	key := append(PrefixNFT, delimiter...)
	if len(classID) > 0 {
		key = append(key, []byte(classID)...)
		key = append(key, delimiter...)
	}

	if len(classID) > 0 && len(tokenID) > 0 {
		key = append(key, []byte(tokenID)...)
	}
	return key
}

// KeyCollection gets the storeKey by the collection
func KeyCollection(classID string) []byte {
	key := append(PrefixCollection, delimiter...)
	return append(key, []byte(classID)...)
}

// KeyClassID gets the storeKey by the class id
func KeyClassID(id string) []byte {
	key := append(PrefixClass, delimiter...)
	return append(key, []byte(id)...)
}

// KeyClassName gets the storeKey by the class name
func KeyClassName(name string) []byte {
	key := append(PrefixClassName, delimiter...)
	return append(key, []byte(name)...)
}
