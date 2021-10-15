package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/nft/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixNFT):
			var nftA, nftB types.NFT
			cdc.MustUnmarshal(kvA.Value, &nftA)
			cdc.MustUnmarshal(kvB.Value, &nftB)
			return fmt.Sprintf("%v\n%v", nftA, nftB)
		case bytes.Equal(kvA.Key[:1], types.PrefixOwners):
			idA := types.MustUnMarshalTokenID(cdc, kvA.Value)
			idB := types.MustUnMarshalTokenID(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)
		case bytes.Equal(kvA.Key[:1], types.PrefixCollection):
			supplyA := types.MustUnMarshalSupply(cdc, kvA.Value)
			supplyB := types.MustUnMarshalSupply(cdc, kvB.Value)
			return fmt.Sprintf("%d\n%d", supplyA, supplyB)
		case bytes.Equal(kvA.Key[:1], types.PrefixClass):
			var classA, classB types.Class
			cdc.MustUnmarshal(kvA.Value, &classA)
			cdc.MustUnmarshal(kvB.Value, &classB)
			return fmt.Sprintf("%v\n%v", classA, classB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
