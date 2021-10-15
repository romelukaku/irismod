package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/nft/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	// Issue a class
	r.HandleFunc("/nft/nfts/classes/issue", issueClassHandlerFn(cliCtx)).Methods("POST")
	// Mint an NFT
	r.HandleFunc("/nft/nfts/mint", mintNFTHandlerFn(cliCtx)).Methods("POST")
	// Update an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}", RestParamClassID, RestParamTokenID), editNFTHandlerFn(cliCtx)).Methods("PUT")
	// Transfer an NFT to an address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/transfer", RestParamClassID, RestParamTokenID), transferNFTHandlerFn(cliCtx)).Methods("POST")
	// Burn an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/burn", RestParamClassID, RestParamTokenID), burnNFTHandlerFn(cliCtx)).Methods("POST")
}

func issueClassHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueClassReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgIssueClass(req.ID, req.Name, req.Schema, req.Owner, req.Symbol, req.MintRestricted, req.UpdateRestricted)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func mintNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if req.Recipient == "" {
			req.Recipient = req.Owner
		}
		// create the message
		msg := types.NewMsgMintNFT(
			req.ID,
			req.ClassID,
			req.Name,
			req.URI,
			req.Data,
			req.Owner,
			req.Recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgEditNFT(
			vars[RestParamTokenID],
			vars[RestParamClassID],
			req.Name,
			req.URI,
			req.Data, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		if _, err := sdk.AccAddressFromBech32(req.Recipient); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgTransferNFT(
			vars[RestParamTokenID],
			vars[RestParamClassID],
			req.Name,
			req.URI,
			req.Data,
			req.Owner,
			req.Recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func burnNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)

		// create the message
		msg := types.NewMsgBurnNFT(
			req.Owner,
			vars[RestParamTokenID],
			vars[RestParamClassID],
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
