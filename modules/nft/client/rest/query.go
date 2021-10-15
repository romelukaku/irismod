package rest

import (
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/nft/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	// Get the total supply of a collection or owner
	r.HandleFunc(fmt.Sprintf("/%s/collections/{%s}/supply", types.ModuleName, RestParamClassID), querySupply(cliCtx, queryRoute)).Methods("GET")
	// Get the collections of NFTs owned by an address
	r.HandleFunc(fmt.Sprintf("/%s/owners/{%s}", types.ModuleName, RestParamOwner), queryOwner(cliCtx, queryRoute)).Methods("GET")
	// Get all the NFTs from a given collection
	r.HandleFunc(fmt.Sprintf("/%s/collections/{%s}", types.ModuleName, RestParamClassID), queryCollection(cliCtx, queryRoute)).Methods("GET")
	// Query all classes
	r.HandleFunc(fmt.Sprintf("/%s/classes", types.ModuleName), queryClasses(cliCtx, queryRoute)).Methods("GET")
	// Query the class
	r.HandleFunc(fmt.Sprintf("/%s/classes/{%s}", types.ModuleName, RestParamClassID), queryClass(cliCtx, queryRoute)).Methods("GET")
	// Query a single NFT
	r.HandleFunc(fmt.Sprintf("/%s/nfts/{%s}/{%s}", types.ModuleName, RestParamClassID, RestParamTokenID), queryNFT(cliCtx, queryRoute)).Methods("GET")
}

func querySupply(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classID := mux.Vars(r)[RestParamClassID]
		err := types.ValidateClassID(classID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		var owner sdk.AccAddress
		ownerStr := r.FormValue(RestParamOwner)
		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		params := types.NewQuerySupplyParams(classID, owner)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySupply), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		out := binary.LittleEndian.Uint64(res)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, out)
	}
}

func queryOwner(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := mux.Vars(r)[RestParamOwner]
		if len(ownerStr) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "param owner should not be empty")
		}

		address, err := sdk.AccAddressFromBech32(ownerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		classID := r.FormValue(RestParamClassID)
		params := types.NewQueryOwnerParams(classID, address)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryOwner), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCollection(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classID := mux.Vars(r)[RestParamClassID]
		if err := types.ValidateClassID(classID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryCollectionParams(classID)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCollection), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryClass(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		classID := mux.Vars(r)[RestParamClassID]
		if err := types.ValidateClassID(classID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryClassParams(classID)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryClass), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryClasses(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryClasses), nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryNFT(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		classID := vars[RestParamClassID]
		if err := types.ValidateClassID(classID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		tokenID := vars[RestParamTokenID]
		if err := types.ValidateTokenID(tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryNFTParams(classID, tokenID)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryNFT), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
