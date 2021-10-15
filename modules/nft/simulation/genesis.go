package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/irisnet/irismod/modules/nft/types"
)

const (
	kitties = "kitties"
	doggos  = "doggos"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	collections := types.NewCollections(
		types.NewCollection(
			types.Class{
				Id:      doggos,
				Name:    doggos,
				Schema:  "",
				Creator: "",
				Symbol:  "dog",
			},
			types.NFTs{},
		),
		types.NewCollection(
			types.Class{
				Id:      kitties,
				Name:    kitties,
				Schema:  "",
				Creator: "",
				Symbol:  "kit",
			},
			types.NFTs{}),
	)
	for _, acc := range simState.Accounts {
		// 10% of accounts own an NFT
		if simState.Rand.Intn(100) < 10 {
			NFT := types.NewNFT(
				RandnNFTID(simState.Rand, types.MinClassLen, types.MaxClassLen), // id
				simtypes.RandStringOfLength(simState.Rand, 10),
				acc.Address,
				simtypes.RandStringOfLength(simState.Rand, 45), // tokenURI
				simtypes.RandStringOfLength(simState.Rand, 10),
			)

			// 50% doggos and 50% kitties
			if simState.Rand.Intn(100) < 50 {
				collections[0].Class.Creator = NFT.Owner
				collections[0] = collections[0].AddNFT(NFT)
			} else {
				collections[1].Class.Creator = NFT.Owner
				collections[1] = collections[1].AddNFT(NFT)
			}
		}
	}

	nftGenesis := types.NewGenesisState(collections)

	bz, err := json.MarshalIndent(nftGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(nftGenesis)
}

func RandnNFTID(r *rand.Rand, min, max int) string {
	n := simtypes.RandIntBetween(r, min, max)
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(id)
}
