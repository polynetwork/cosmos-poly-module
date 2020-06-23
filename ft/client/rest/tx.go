/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package rest

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"encoding/hex"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/ft/create_coins/{%s}", Coins), CreateCoinsRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/ft/create_denom/{%s}", Denom), CreateDenomRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/ft/bind_asset_hash", BindAssetHashRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/ft/lock", LockRequestHandlerFn(cliCtx)).Methods("POST")

}

type CreateReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

type BindAssetHashReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom       string       `json:"denom" yaml:"denom"`
	ToChainId   uint64       `json:"to_chain_id" yaml:"to_chain_id"`
	ToAssetHash string       `json:"to_asset_hash" yaml:"to_asset_hash"`
}
type LockReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom     string       `json:"denom" yaml:"denom"`
	ToChainId uint64       `json:"to_chain_id" yaml:"to_chain_id"`
	ToAddress []byte       `json:"to_address" yaml:"to_address"`
	Amount    *big.Int     `json:"amount" yaml:"amount"`
}

func CreateCoinsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		_, err := sdk.ParseCoin(vars[Coins])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req CreateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		msg := types.NewMsgCreateCoins(cliCtx.GetFromAddress(), vars[Coins])
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func CreateDenomRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[Denom]

		var req CreateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgCreateDenom(cliCtx.GetFromAddress(), denom)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func BindAssetHashRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BindAssetHashReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		toAssetHash, err := hex.DecodeString(req.ToAssetHash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), req.Denom, req.ToChainId, toAssetHash)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func LockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		msg := types.NewMsgLock(cliCtx.GetFromAddress(), req.Denom, req.ToChainId, req.ToAddress, sdk.NewIntFromBigInt(req.Amount))
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
